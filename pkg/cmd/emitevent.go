package cmd

import (
	"fmt"
	"os"

	"github.com/rajatjindal/kubectl-emitevent/pkg/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

//Version is set during build time
var Version = "unknown"

//EmitEventOptions is struct for emitevent command
type EmitEventOptions struct {
	configFlags *genericclioptions.ConfigFlags
	iostreams   genericclioptions.IOStreams

	args         []string
	kubeclient   kubernetes.Interface
	printVersion bool
}

// NewEmitEventOptions provides an instance of EmitEventOptions with default values
func NewEmitEventOptions(streams genericclioptions.IOStreams) *EmitEventOptions {
	return &EmitEventOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		iostreams:   streams,
	}
}

// NewCmdEmitEvent provides a cobra command wrapping EmitEventOptions
func NewCmdEmitEvent(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewEmitEventOptions(streams)

	cmd := &cobra.Command{
		Use:          "emitevent [flags]",
		Short:        "emits an event to given object",
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
func (o *EmitEventOptions) Complete(cmd *cobra.Command, args []string) error {
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
func (o *EmitEventOptions) Validate() error {
	if len(o.args) > 0 {
		return fmt.Errorf("no arguments expected. got %d arguments", len(o.args))
	}

	return nil
}

// Run emits the event to given k8s object
func (o *EmitEventOptions) Run() error {
	pod, err := o.kubeclient.CoreV1().Pods("kube-system").Get("coredns-6955765f44-7lh5l", metav1.GetOptions{})
	if err != nil {
		return err
	}

	k8s.EmitEvent(o.kubeclient, pod, "my-own-reason", "my-own-message")
	return nil
}
