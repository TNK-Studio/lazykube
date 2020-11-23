package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/exec"
)

func (cli *KubeCLI) Exec(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := exec.NewCmdExec(cli.factory, streams)
	return NewCmd(cmd, args, streams)
}
