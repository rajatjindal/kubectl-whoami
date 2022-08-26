package cmd

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rajatjindal/kubectl-whoami/pkg/k8s"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/transport"
)

//Version is set during build time
var Version = "unknown"

const JSONFORMAT = "json"

//WhoAmIOptions is struct for whoami command
type WhoAmIOptions struct {
	configFlags *genericclioptions.ConfigFlags
	iostreams   genericclioptions.IOStreams

	args         []string
	kubeclient   kubernetes.Interface
	printVersion bool
	all          bool
	outputFormat string

	tokenRetriever *tokenRetriever
}

// tokenRetriever helps to retrieve token
type tokenRetriever struct {
	rountTripper http.RoundTripper
	token        string
}

//RoundTrip gets token
func (t *tokenRetriever) RoundTrip(req *http.Request) (*http.Response, error) {
	header := req.Header.Get("authorization")
	switch {
	case strings.HasPrefix(header, "Bearer "):
		t.token = strings.ReplaceAll(header, "Bearer ", "")
	}

	return t.rountTripper.RoundTrip(req)
}

// NewWhoAmIOptions provides an instance of WhoAmIOptions with default values
func NewWhoAmIOptions(streams genericclioptions.IOStreams) *WhoAmIOptions {
	return &WhoAmIOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		iostreams:   streams,
	}
}

// NewCmdWhoAmI provides a cobra command wrapping WhoAmIOptions
func NewCmdWhoAmI(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewWhoAmIOptions(streams)

	cmd := &cobra.Command{
		Use:          "whoami [flags]",
		Short:        "find out the subject for the current context of kubeconfig",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if o.printVersion {
				fmt.Println(Version)
				os.Exit(0)
			}

			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().BoolVar(&o.all, "all", false, "Prints information about user, groups and ARN")
	cmd.Flags().BoolVar(&o.printVersion, "version", false, "prints version of plugin")
	cmd.Flags().StringVarP(&o.outputFormat, "output", "o", "", "Output format. One of: json")
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets all information required for updating the current context
func (o *WhoAmIOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	o.tokenRetriever = &tokenRetriever{}
	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		o.tokenRetriever.rountTripper = rt
		return o.tokenRetriever
	})

	o.kubeclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *WhoAmIOptions) Validate() error {
	if len(o.args) > 0 {
		return fmt.Errorf("no arguments expected. got %d arguments", len(o.args))
	}

	return nil
}

// Run retrieves and print the subject that's currently authenticated
func (o *WhoAmIOptions) Run() error {
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	c, err := config.TransportConfig()
	if err != nil {
		return err
	}

	var token string
	// from vendor/k8s.io/client-go/transport/round_trippers.go:HTTPWrappersForConfig function, tokenauth has preference over basicauth
	if c.HasTokenAuth() {
		if config.BearerTokenFile != "" {
			d, err := ioutil.ReadFile(config.BearerTokenFile)
			if err != nil {
				return err
			}

			token = string(d)
		}

		if config.BearerToken != "" {
			token = config.BearerToken
		}
	}

	if token == "" && (config.AuthProvider != nil || config.ExecProvider != nil) {
		token, err = o.getToken()
		if err != nil {
			return err
		}
	}

	if token != "" {
		username, tokenreview, err := k8s.WhoAmI(o.kubeclient, token)
		if err != nil {
			return err
		}

		if o.outputFormat == JSONFORMAT {
			if tokenreview == nil {
				return fmt.Errorf("tokenreview did not return a status")
			}
			out, err := json.MarshalIndent(tokenreview, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal json from tokenreview")
			}
			fmt.Println(string(out))
			return nil
		}
		if o.all {
			username = fmt.Sprintf("User:\t%s\nGroups:\n\t%s", tokenreview.User.Username,
				strings.Join(tokenreview.User.Groups, "\n\t"))
			if len(tokenreview.User.Extra["arn"]) > 0 {
				username = username + "\nARN:\n\t" + strings.Join(tokenreview.User.Extra["arn"], "\n\t")
			}
		}

		if o.outputFormat != JSONFORMAT {
			if username == "" {
				return fmt.Errorf("failed to find subject of token. please report a ticket at https://github.com/rajatjindal/kubectl-whoami")
			}
			fmt.Println(username)
		}
		return nil
	}

	if c.HasBasicAuth() {
		fmt.Printf("kubecfg:basicauth:%s", config.Username)
		return nil
	}

	if c.HasCertAuth() {
		cert, err := getClientCertificate(c)
		if err != nil {
			return err
		}
		username := cert.Subject.CommonName
		fmt.Println(username)
		return nil
	}

	return fmt.Errorf("unsupported auth mechanism. kindly report a ticket at https://github.com/rajatjindal/kubectl-whoami")
}

func (o *WhoAmIOptions) getToken() (string, error) {
	err := k8s.WhatCanI(o.kubeclient)
	if err != nil {
		return "", err
	}
	return o.tokenRetriever.token, nil
}

func getClientCertificate(c *transport.Config) (*x509.Certificate, error) {
	tlsConfig, err := transport.TLSConfigFor(c)
	if err != nil {
		return nil, err
	}
	// GetClientCertificate has been set in transport.TLSConfigFor,
	// so it is not nil
	cert, err := tlsConfig.GetClientCertificate(nil)
	if err != nil {
		return nil, err
	}
	if cert.Leaf != nil {
		return cert.Leaf, nil
	}
	return x509.ParseCertificate(cert.Certificate[0])
}
