package app

import "github.com/TNK-Studio/lazykube/pkg/kubecli"

func cli(namespace string) *kubecli.KubeCLI {
	if namespace == kubecli.Cli.Namespace() {
		return kubecli.Cli
	}
	return kubecli.Cli.WithNamespace(namespace)
}
