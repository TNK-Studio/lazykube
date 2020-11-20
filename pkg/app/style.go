package app

import (
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
)

const (
	tallPanels = 4
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

func usableSpace(_ *guilib.Gui, maxHeight int) int {
	if maxHeight < 28 {
		return maxHeight - 2
	}
	return maxHeight - 8
}

func resizePanelHeight(gui *guilib.Gui) error {
	_, maxHeight := gui.Size()

	space := usableSpace(gui, maxHeight)

	viewHeights[clusterInfoViewName] = 2
	viewHeights[namespaceViewName] = space / tallPanels
	viewHeights[serviceViewName] = space / tallPanels
	viewHeights[deploymentViewName] = space / tallPanels
	viewHeights[podViewName] = space / tallPanels
	viewHeights[optionViewName] = 1

	return nil
}

func reactiveHeight(gui *guilib.Gui, view *guilib.View) int {
	var resizeView string
	currentView := gui.CurrentView()
	if currentView == nil {
		resizeView = namespaceViewName
	} else {
		resizeView = currentView.Name
	}

	// When cluster info 、navigation or detail panel selected.
	if !resizeAbleView(resizeView) {
		resizeView = gui.PeekPreviousView()

		// If previous view is cluster info 、navigation or detail pane.
		if !resizeAbleView(resizeView) {
			resizeView = namespaceViewName
		}
	}

	maxHeight := gui.MaxHeight()
	height := viewHeights[view.Name]
	if resizeView == view.Name {
		height += usableSpace(gui, maxHeight) % tallPanels
	}

	if maxHeight < 28 {
		if view.Name == podViewName {
			// First time
			if currentView == nil {
				return height
			}

			// When pod panel selected.
			if resizeView == podViewName {
				return height - migrateTopFunc(gui, view)*2
			}
		}
	}

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
