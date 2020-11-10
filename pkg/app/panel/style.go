package panel

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
)

func leftSideWidth(maxWidth int) int {
	return maxWidth / 3
}

func usableSpace(maxHeight int) int {
	return maxHeight - 4
}

func reactiveHeight(gui *gui.Gui, view *gui.View) int {
	_, maxHeight := gui.Size()
	currView := gui.CurrentView()
	currentCyclebleView := gui.PeekPreviousView()
	if currView != nil {
		viewName := currView.Name
		usePreviouseView := true
		for _, view := range []string{"namespace", "service", "deployment", "pod"} {
			if view == viewName {
				currentCyclebleView = viewName
				usePreviouseView = false
				break
			}
		}
		if usePreviouseView {
			currentCyclebleView = gui.PeekPreviousView()
		}
	}

	space := usableSpace(maxHeight)

	tallPanels := 4
	vHeights := map[string]int{
		"clusterInfo": 3,
		"namespace":   space / tallPanels,
		"service":     space / tallPanels,
		"deployment":  space / tallPanels,
		"pod":         space / tallPanels,
		"option":      1,
	}

	currentView := gui.CurrentView()
	if currentView != nil {
		vHeights[currentView.Name] += space % tallPanels
	}
	if maxHeight < 28 {
		defaultHeight := 3
		// Todo: Folding panel
		if maxHeight < 21 {
			defaultHeight = 1
		}
		vHeights = map[string]int{
			"clusterInfo": defaultHeight,
			"namespace":   defaultHeight,
			"service":     defaultHeight,
			"deployment":  defaultHeight,
			"pod":         defaultHeight,
			"option":      defaultHeight,
		}

		vHeights[currentCyclebleView] = maxHeight - defaultHeight*tallPanels - 1
	}

	vHeights["clusterInfo"] -= 1
	if vHeights["clusterInfo"] == 0 {
		vHeights["clusterInfo"] = 1
	}
	height := vHeights[view.Name]
	return height
}
