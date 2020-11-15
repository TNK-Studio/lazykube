package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/logs"
)

func (cli *KubeCLI) Logs(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := logs.NewCmdLogs(cli.factory, streams)
	return &Cmd{
		cmd:  cmd,
		args: args,
	}
}
