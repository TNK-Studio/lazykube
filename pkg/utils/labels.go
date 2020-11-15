package utils

import (
	"encoding/json"
	"fmt"
)

// Convert labels which like '{"k8s-app":"kube-dns"}' to '[]string{"k8s-app=kube-dns"}'.
func LabelsToStringArr(labelsMapString string) []string {
	labelsArr := make([]string, 0)
	if labelsMapString == "" {
		return labelsArr
	}

	//labelsMapString = strings.ReplaceAll(labelsMapString, `"`, `\"`)
	labelsMap := make(map[string]string, 0)
	b := []byte(labelsMapString)
	err := json.Unmarshal(b, &labelsMap)
	if err != nil {
		return labelsArr
	}

	for key, val := range labelsMap {
		labelsArr = append(labelsArr, fmt.Sprintf("%s=%s", key, val))
	}
	return labelsArr
}
