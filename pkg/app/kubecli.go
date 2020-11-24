package app

import (
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"strings"
)

const matchLabels = "jsonpath='{.spec.selector.matchLabels}'"

func cli(namespace string) *kubecli.KubeCLI {
	if namespace == kubecli.Cli.Namespace() {
		return kubecli.Cli
	}
	return kubecli.Cli.WithNamespace(namespace)
}

func resourceLabelSelectorJSONPath(resource string) string {
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

func getPodContainers(namespace, podName string) []string {
	// Todo: support others resource
	stream := newStream()
	cli(namespace).
		Get(stream, "pods", podName).
		SetFlag("output", "jsonpath='{.spec.containers[*].name}'").
		Run()

	result := strings.ReplaceAll(streamToString(stream), "'", "")
	containers := strings.Split(result, " ")
	if len(containers) == 0 {
		return []string{}
	}
	return containers
}
