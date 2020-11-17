package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
)

func (cli *KubeCLI) Get(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := get.NewCmdGet("kubectl", cli.factory, streams)
	return NewCmd(cmd, args, streams)
}
