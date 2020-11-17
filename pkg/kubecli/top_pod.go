package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/top"
)

func (cli *KubeCLI) TopPod(streams genericclioptions.IOStreams, o *top.TopPodOptions, args ...string) *Cmd {
	cmd := top.NewCmdTopPod(cli.factory, o, streams)
	return NewCmd(cmd, args, streams)
}
