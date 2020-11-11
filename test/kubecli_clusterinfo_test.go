package test

import (
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"testing"
)

func TestClusterInfo(t *testing.T) {
	cli := kubecli.NewKubeCLI()
	t.Logf("cli %+v \n", cli)

	info, err := cli.ClusterInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", info, "\n")
}
