package app

import (
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
)

var (
	viewHeights     = map[string]int{}
	resizeableViews = []string{namespaceViewName, serviceViewName, deploymentViewName, podViewName}
)

func resizeAbleView(viewName string) bool {
	for _, resizeableView := range resizeableViews {
		if resizeableView == viewName {
			return true
		}
	}
	return false
}

func leftSideWidth(maxWidth int) int {
	return maxWidth / 3
}

func usableSpace(gui *guilib.Gui, maxHeight int) int {
	if maxHeight < 28 {
		if currentView := gui.CurrentView(); currentView != nil && currentView.Name == podViewName {
			return maxHeight
		}

		return maxHeight - 2
	}
	return maxHeight - 8
}

func resizePanelHeight(gui *guilib.Gui) error {
	_, maxHeight := gui.Size()

	space := usableSpace(gui, maxHeight)

	tallPanels := 4
	viewHeights[clusterInfoViewName] = 2
	viewHeights[namespaceViewName] = space / tallPanels
	viewHeights[serviceViewName] = space / tallPanels
	viewHeights[deploymentViewName] = space / tallPanels
	viewHeights[podViewName] = space / tallPanels
	viewHeights[optionViewName] = 1

	resizeView := namespaceViewName
	currentView := gui.CurrentView()
	if currentView != nil && resizeAbleView(currentView.Name) {
		resizeView = currentView.Name
	} else if gui.PeekPreviousView() != "" && gui.PeekPreviousView() != clusterInfoViewName {
		resizeView = gui.PeekPreviousView()
	}

	viewHeights[resizeView] += space % tallPanels
	return nil
}

func reactiveHeight(_ *guilib.Gui, view *guilib.View) int {
	height := viewHeights[view.Name]
	return height
}

func migrateTopFunc(gui *guilib.Gui, view *guilib.View) int {
	maxHeight := gui.MaxHeight()

	if maxHeight < 28 {
		currentView := gui.CurrentView()
		if currentView != nil {
			if currentView.Name == navigationViewName || currentView.Name == detailViewName {
				if view.Name == namespaceViewName {
					return 1
				}
			}

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
	if maxHeight < 28 {
		return -1
	}
	return 1
}
