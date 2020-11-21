package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apiresources"
)

func (cli *KubeCLI) APIResources(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := apiresources.NewCmdAPIResources(cli.factory, streams)
	return NewCmd(cmd, args, streams)
}
