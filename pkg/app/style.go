package app

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
)

var (
	viewHeights     = map[string]int{}
	resizeView      = ""
	resizeableViews = []string{namespaceViewName, serviceViewName, deploymentViewName, podViewName}
)

func leftSideWidth(maxWidth int) int {
	return maxWidth / 3
}

func usableSpace(gui *gui.Gui, maxHeight int) int {
	if maxHeight < 28 {
		if currentView := gui.CurrentView(); currentView != nil && currentView.Name == podViewName {
			return maxHeight
		}

		return maxHeight - 2
	}
	return maxHeight - 8
}

func reactiveHeight(gui *gui.Gui, view *gui.View) int {
	_, maxHeight := gui.Size()

	space := usableSpace(gui, maxHeight)

	tallPanels := 4
	viewHeights[clusterInfoViewName] = 2
	viewHeights[namespaceViewName] = space / tallPanels
	viewHeights[serviceViewName] = space / tallPanels
	viewHeights[deploymentViewName] = space / tallPanels
	viewHeights[podViewName] = space / tallPanels
	viewHeights[optionViewName] = 1

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

	height := viewHeights[view.Name]
	return height
}

func migrateTopFunc(gui *gui.Gui, view *gui.View) int {
	if gui.MaxHeight() < 28 {
		currentView := gui.CurrentView()
		if currentView != nil {
			for i, viewName := range functionViews {
				if currentView.Name == viewName {
					index := i + 1
					if index < len(functionViews) && functionViews[index] == view.Name {
						return 1
					}
				}
			}
		}
	}
	if gui.MaxHeight() < 28 {
		return -1
	}
	return 1
}
