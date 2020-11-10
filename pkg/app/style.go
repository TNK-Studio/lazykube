package app

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
)

var (
	viewHeights     = map[string]int{}
	resizeView      = ""
	resizeableViews = []string{"namespace", "service", "deployment", "pod"}
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

	tallPanels := 4
	viewHeights = map[string]int{
		"clusterInfo": 3,
		"namespace":   space / tallPanels,
		"service":     space / tallPanels,
		"deployment":  space / tallPanels,
		"pod":         space / tallPanels,
		"option":      1,
	}

	currentView := gui.CurrentView()
	if currentView != nil {
		for _, viewName := range resizeableViews {
			if currentView.Name == viewName {
				resizeView = viewName
				break
			}
		}

		viewHeights[resizeView] += space % tallPanels
	}

	if maxHeight < 28 {
		defaultHeight := 3
		// Todo: Folding panel
		if maxHeight < 21 {
			defaultHeight = 1
		}
		viewHeights = map[string]int{
			"clusterInfo": defaultHeight,
			"namespace":   defaultHeight,
			"service":     defaultHeight,
			"deployment":  defaultHeight,
			"pod":         defaultHeight,
			"option":      defaultHeight,
		}

		viewHeights[resizeView] = maxHeight - defaultHeight*tallPanels - 1
	}

	viewHeights["clusterInfo"] -= 1
	if viewHeights["clusterInfo"] == 0 {
		viewHeights["clusterInfo"] = 1
	}
	height := viewHeights[view.Name]
	return height
}

func migrateTopFunc(gui *gui.Gui, view *gui.View) int {
	return 1
}
