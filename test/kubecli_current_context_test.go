package test

import (
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"testing"
)

func TestCurrentContext(t *testing.T) {
	cli := kubecli.NewKubeCLI()
	t.Logf("cli %+v \n", cli)

	currentContext, err := cli.CurrentContext()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(currentContext, "\n")
}
