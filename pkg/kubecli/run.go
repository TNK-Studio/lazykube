package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/run"
)

func (cli *KubeCLI) Run(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := run.NewCmdRun(cli.factory, streams)
	return NewCmd(cmd, args, streams)
}
