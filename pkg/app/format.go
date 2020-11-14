package app

import (
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"strings"
)

func formatSelectedNamespace(selected string) string {
	return formatSelectedName(selected, 0)
}

func formatSelectedName(selected string, index int) string {
	if selected == "" {
		return ""
	}
	selected = utils.DeleteExtraSpace(selected)
	formatted := strings.Split(selected, " ")
	length := len(formatted)
	if index < 0 && length-index >= 0 {
		return formatted[length-index]
	}

	if index < length {
		return formatted[index]
	}
	return ""
}
