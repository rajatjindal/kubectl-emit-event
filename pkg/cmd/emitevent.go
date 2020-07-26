package cmd

import (
	"fmt"

	"github.com/rajatjindal/kubectl-emit-event/pkg/k8s"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	//needed for cloud provider auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

const usage = `kubectl emit-event [flags]

Example

## emitevent
➜  kubectl emit-event daemonset/kube-proxy -n kube-system --reason "foo-reason" --message "bar-message"

## verify event was generated and attached to the object
➜  kubectl describe daemonset/kube-proxy -n kube-system

Name:           kube-proxy
Selector:       k8s-app=kube-proxy
Node-Selector:  beta.kubernetes.io/os=linux
Labels:         k8s-app=kube-proxy
Annotations:    deprecated.daemonset.template.generation: 1
.
.
.
.
.
Events:
  Type    Reason      Age    From                Message
  ----    ------      ----   ----                -------
  Normal  foo-reason  13s    kubectl-emit-event  bar-message
`

//Version is set during build time
var Version = "unknown"

//EmitEventOptions is struct for emitevent command
type EmitEventOptions struct {
	configFlags *genericclioptions.ConfigFlags
	iostreams   genericclioptions.IOStreams
	factory     cmdutil.Factory

	Namespace        string
	EnforceNamespace bool
	BuilderArgs      []string

	Builder       func() *resource.Builder
	DynamicClient dynamic.Interface
	kubeclient    kubernetes.Interface

	FilenameOptions *resource.FilenameOptions

	eventType string
	reason    string
	message   string
}

// NewEmitEventOptions provides an instance of EmitEventOptions with default values
func NewEmitEventOptions(streams genericclioptions.IOStreams) *EmitEventOptions {
	return &EmitEventOptions{
		configFlags:     genericclioptions.NewConfigFlags(true),
		FilenameOptions: &resource.FilenameOptions{},
		iostreams:       streams,
	}
}

// NewCmdEmitEvent provides a cobra command wrapping EmitEventOptions
func NewCmdEmitEvent(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewEmitEventOptions(streams)

	cmd := &cobra.Command{
		Use:          usage,
		Short:        "Emits a Kubernetes Event for the requested object",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
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

	cmd.Flags().StringVar(&o.eventType, "type", corev1.EventTypeNormal, "Type of event")

	cmd.Flags().StringVar(&o.reason, "reason", "", "reason of event")
	cmd.MarkFlagRequired("reason")

	cmd.Flags().StringVar(&o.message, "message", "", "message to be passed")
	cmd.MarkFlagRequired("message")

	o.configFlags.AddFlags(cmd.Flags())

	usage := "identifying the resource to get from a server."
	cmdutil.AddFilenameOptionFlags(cmd, o.FilenameOptions, usage)

	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(o.configFlags)
	o.factory = cmdutil.NewFactory(matchVersionKubeConfigFlags)

	return cmd
}

// Complete sets all information required for updating the current context
func (o *EmitEventOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Builder = o.factory.NewBuilder
	o.BuilderArgs = args

	var err error
	o.Namespace, o.EnforceNamespace, err = o.factory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	clientConfig, err := o.factory.ToRESTConfig()
	if err != nil {
		return err
	}

	o.kubeclient, err = kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	o.DynamicClient, err = dynamic.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *EmitEventOptions) Validate() error {
	if len(o.BuilderArgs) == 0 && cmdutil.IsFilenameSliceEmpty(o.FilenameOptions.Filenames, o.FilenameOptions.Kustomize) {
		return fmt.Errorf("required resource not specified")
	}

	return nil
}

// Run emits the event to given k8s object
func (o *EmitEventOptions) Run() error {
	r := o.Builder().
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		NamespaceParam(o.Namespace).DefaultNamespace().
		FilenameParam(o.EnforceNamespace, o.FilenameOptions).
		ResourceTypeOrNameArgs(true, o.BuilderArgs...).
		SingleResourceType().
		Latest().
		Do()

	err := r.Err()
	if err != nil {
		return err
	}

	object, err := r.Object()
	if err != nil {
		return err
	}

	k8s.EmitEvent(o.kubeclient, object, o.eventType, o.reason, o.message)
	return nil
}
