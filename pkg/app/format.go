package app

import (
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"strings"
)

func formatSelectedNamespace(selected string) string {
	return formatResourceName(selected, 0)
}

func formatResourceName(selected string, index int) string {
	if selected == "" {
		return ""
	}
	selected = utils.DeleteExtraSpace(selected)
	formatted := strings.Split(selected, " ")
	length := len(formatted)
	if index < 0 && length-index >= 0 {
		resourceName := formatted[length-index]
		if validateResourceName(resourceName) {
			return resourceName
		}
	}

	if index < length {
		resourceName := formatted[index]
		if validateResourceName(resourceName) {
			return resourceName
		}
	}
	return ""
}

func validateResourceName(resourceName string) bool {
	return !utils.IsUpper(resourceName)
}
