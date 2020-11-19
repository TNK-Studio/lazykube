package kubecli

import (
	"flag"
	"fmt"
	"github.com/TNK-Studio/lazykube/pkg/kubecli/clusterinfo"
	"github.com/TNK-Studio/lazykube/pkg/kubecli/config"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/cmd/util"
)

var Cli *KubeCLI

func init() {
	Cli = NewKubeCLI()
	// To disable aws warning
	disableKlog()
}

type KubeCLI struct {
	factory   util.Factory
	namespace *string
}

type Cmd struct {
	cmd     *cobra.Command
	args    []string
	streams genericclioptions.IOStreams
}

func NewCmd(cmd *cobra.Command, args []string, streams genericclioptions.IOStreams) *Cmd {
	return &Cmd{
		cmd:     cmd,
		args:    args,
		streams: streams,
	}
}

func (c *Cmd) Run() {
	util.BehaviorOnFatal(func(s string, i int) {
		_, _ = fmt.Fprint(c.streams.ErrOut, s)
	})
	c.cmd.Run(c.cmd, c.args)
}

func (c *Cmd) SetFlag(name, value string) *Cmd {
	if err := c.cmd.Flags().Set(name, value); err != nil {
		log.Logger.Panicln(err)
	}
	return c
}

func NewKubeCLI() *KubeCLI {
	namespace := ""
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.Namespace = &namespace

	matchVersionKubeConfigFlags := util.NewMatchVersionFlags(kubeConfigFlags)

	k := &KubeCLI{
		factory:   util.NewFactory(matchVersionKubeConfigFlags),
		namespace: &namespace,
	}
	return k
}

func (cli *KubeCLI) SetNamespace(namespace string) {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.Namespace = &namespace

	matchVersionKubeConfigFlags := util.NewMatchVersionFlags(kubeConfigFlags)
	cli.factory = util.NewFactory(matchVersionKubeConfigFlags)
	cli.namespace = &namespace
}

func (cli *KubeCLI) WithNamespace(namespace string) *KubeCLI {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.Namespace = &namespace

	matchVersionKubeConfigFlags := util.NewMatchVersionFlags(kubeConfigFlags)

	k := &KubeCLI{
		factory:   util.NewFactory(matchVersionKubeConfigFlags),
		namespace: &namespace,
	}
	return k
}

func (cli *KubeCLI) Namespace() string {
	return *cli.namespace
}

func (cli *KubeCLI) CurrentContext() (string, error) {
	return config.CurrentContext()
}

func (cli *KubeCLI) ClusterInfo() (string, error) {
	return clusterinfo.ClusterInfo(cli.factory)
}

func disableKlog() {
	flagSet := &flag.FlagSet{}
	klog.InitFlags(flagSet)
	_ = flagSet.Set("logtostderr", "false")
	klog.SetOutput(ioutil.Discard)
}
