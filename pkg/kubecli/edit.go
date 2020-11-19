package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/edit"
)

func (cli *KubeCLI) Edit(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := edit.NewCmdEdit(cli.factory, streams)
	return NewCmd(cmd, args, streams)
}
