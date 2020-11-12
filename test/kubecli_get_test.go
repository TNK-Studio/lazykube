package test

import (
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"testing"
)

func TestKubeCLI_Get(t *testing.T) {
	cli := kubecli.NewKubeCLI()
	t.Logf("cli %+v \n", cli)

	cli.Get(
		genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stdout},
		"nodes",
	)
}
