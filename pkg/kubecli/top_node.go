package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/top"
)

func (cli *KubeCLI) TopNode(streams genericclioptions.IOStreams, o *top.TopNodeOptions, args ...string) *Cmd {
	cmd := top.NewCmdTopNode(cli.factory, o, streams)
	return NewCmd(cmd, args, streams)
}
