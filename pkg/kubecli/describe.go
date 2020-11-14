package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/describe"
)

func (cli *KubeCLI) Describe(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := describe.NewCmdDescribe("kubectl", cli.factory, streams)
	return &Cmd{
		cmd:  cmd,
		args: args,
	}
}
