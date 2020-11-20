package kubecli

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/rollout"
)

func (cli *KubeCLI) RolloutRestart(streams genericclioptions.IOStreams, args ...string) *Cmd {
	cmd := rollout.NewCmdRolloutRestart(cli.factory, streams)
	return NewCmd(cmd, args, streams)
}
