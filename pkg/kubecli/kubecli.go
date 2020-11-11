package kubecli

import (
	"github.com/TNK-Studio/lazykube/pkg/kubecli/clusterinfo"
	"github.com/TNK-Studio/lazykube/pkg/kubecli/config"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
)

var Cli *KubeCLI

func init() {
	Cli = NewKubeCLI()
}

type KubeCLI struct {
	factory util.Factory
}

func NewKubeCLI() *KubeCLI {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionKubeConfigFlags := util.NewMatchVersionFlags(kubeConfigFlags)

	k := &KubeCLI{
		factory: util.NewFactory(matchVersionKubeConfigFlags),
	}
	return k
}

func (cli *KubeCLI) CurrentContext() (string, error) {
	return config.CurrentContext()
}

func (cli *KubeCLI) ClusterInfo() (string, error) {
	return clusterinfo.ClusterInfo(cli.factory)
}
