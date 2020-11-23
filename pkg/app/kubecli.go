package app

import "github.com/TNK-Studio/lazykube/pkg/kubecli"

const matchLabels = "jsonpath='{.spec.selector.matchLabels}'"

func cli(namespace string) *kubecli.KubeCLI {
	if namespace == kubecli.Cli.Namespace() {
		return kubecli.Cli
	}
	return kubecli.Cli.WithNamespace(namespace)
}

func resourceLabelSelectorJsonPath(resource string) string {
	var jsonPath string
	switch resource {
	case "services", "service", "svc":
		jsonPath = "jsonpath='{.spec.selector}'"
	case "deployments", "deployment", "deploy":
		jsonPath = matchLabels
	case "statefulsets", "statefulset", "sts":
		jsonPath = matchLabels
	case "daemonsets", "daemonset", "ds":
		jsonPath = matchLabels
	}
	return jsonPath
}
