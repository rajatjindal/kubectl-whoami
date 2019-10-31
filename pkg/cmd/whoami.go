package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rajatjindal/kubectl-whoami/pkg/k8s"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

//Version is set during build time
var Version = "unknown"

//WhoAmIOptions is struct for modify secret
type WhoAmIOptions struct {
	configFlags *genericclioptions.ConfigFlags
	iostreams   genericclioptions.IOStreams

	args         []string
	kubeclient   kubernetes.Interface
	printVersion bool
}

// NewWhoAmIOptions provides an instance of WhoAmIOptions with default values
func NewWhoAmIOptions(streams genericclioptions.IOStreams) *WhoAmIOptions {
	return &WhoAmIOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		iostreams:   streams,
	}
}

// NewCmdModifySecret provides a cobra command wrapping WhoAmIOptions
func NewCmdModifySecret(streams genericclioptions.IOStreams) *cobra.Command {
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

	cmd.Flags().BoolVar(&o.printVersion, "version", false, "prints version of plugin")
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

// Run fetches the given secret manifest from the cluster, decodes the payload, opens an editor to make changes, and applies the modified manifest when done
func (o *WhoAmIOptions) Run() error {
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	c, err := config.TransportConfig()
	if err != nil {
		return err
	}

	// from vendor/k8s.io/client-go/transport/round_trippers.go:HTTPWrappersForConfig function, tokenauth has preference over basicauth
	if c.HasTokenAuth() {
		var token string
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

		username, err := k8s.WhoAmI(o.kubeclient, token)
		if err != nil {
			return err
		}

		if username == "" {
			return fmt.Errorf("failed to find subject of token. please report a ticket at https://github.com/rajatjindal/kubectl-whoami")
		}

		fmt.Println(username)
		return nil
	}

	if c.HasBasicAuth() {
		fmt.Printf("kubecfg:basicauth:%s", config.Username)
		return nil
	}

	if c.HasCertAuth() {
		fmt.Println("kubecfg:certauth:admin")
		return nil

	}

	return fmt.Errorf("could not identify auth type")
}
