package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rajatjindal/kubectl-whoami/pkg/k8s"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

//Version is set during build time
var Version = "unknown"

//WhoAmIOptions is struct for modify secret
type WhoAmIOptions struct {
	configFlags *genericclioptions.ConfigFlags
	iostreams   genericclioptions.IOStreams

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
	return nil
}

// Run fetches the given secret manifest from the cluster, decodes the payload, opens an editor to make changes, and applies the modified manifest when done
func (o *WhoAmIOptions) Run() error {
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	if config.Username != "" {
		fmt.Printf("kubecfg:basicauth:%s", config.Username)
		return nil
	}

	if (config.CAData != nil || config.CAFile != "") && (config.CertData != nil || config.CertFile != "") {
		fmt.Println("kubecfg:certauth:admin")
		return nil
	}

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

	fmt.Println(username)
	return nil
}
