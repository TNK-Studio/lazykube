package panel

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/sirupsen/logrus"
)

func leftSideWidth(maxWidth int) int {
	return maxWidth / 3
}

func usableSpace(maxHeight int) int {
	return maxHeight - 4
}

func reactiveHeight(gui *gui.Gui, view *gui.View) int {
	_, maxHeight := gui.Size()
	space := usableSpace(maxHeight)

	// Todo: dynamic calculate
	tallPanels := 4
	vHeights := map[string]int{
		"clusterInfo": 3,
		"namespace":   space / tallPanels,
		"service":     space / tallPanels,
		"deployment":  space / tallPanels,
		"pod":         space / tallPanels,
		"options":     1,
	}
	//if maxHeight < 28 {
	//	defaultHeight := 3
	//	if maxHeight < 21 {
	//		defaultHeight = 1
	//	}
	//	vHeights = map[string]int{
	//		"clusterInfo": defaultHeight,
	//		"namespace":   defaultHeight,
	//		"images":      defaultHeight,
	//		"volumes":     defaultHeight,
	//		"options":     defaultHeight,
	//	}
	//	//if gui.DockerCommand.InDockerComposeProject {
	//	//	vHeights["services"] = defaultHeight
	//	//}
	//	//vHeights[currentCyclebleView] = height - defaultHeight*tallPanels - 1
	//}

	height := vHeights[view.Name] - 1
	logrus.Debugf("View '%s' height %d", view.Name, height)
	return height
}
