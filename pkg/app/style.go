package app

import (
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
)

const (
	resizeAbleViewMinHeight = 5
	clusterInfoViewHeight   = 2
	optionViewHeight        = 1
)

var (
	viewHeights     = map[string]int{}
	functionViews   = []string{clusterInfoViewName, namespaceViewName, serviceViewName, deploymentViewName, podViewName}
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
	usedSpace := clusterInfoViewHeight + optionViewHeight + (len(functionViews) - 1)
	if maxHeight < resizeAbleViewMinHeight*len(resizeableViews)-usedSpace {
		return maxHeight - 2
	}
	return maxHeight - usedSpace
}

func resizePanelHeight(gui *guilib.Gui) error {
	_, maxHeight := gui.Size()

	space := usableSpace(gui, maxHeight)

	n := len(resizeableViews)
	viewHeights[clusterInfoViewName] = clusterInfoViewHeight
	viewHeights[namespaceViewName] = space / n
	viewHeights[serviceViewName] = space / n
	viewHeights[deploymentViewName] = space / n
	viewHeights[podViewName] = space / n
	viewHeights[optionViewName] = optionViewHeight

	return nil
}

func reactiveHeight(gui *guilib.Gui, view *guilib.View) int {
	var resizeView string
	currentView := gui.CurrentView()
	if currentView == nil {
		resizeView = resizeableViews[0]
	} else {
		resizeView = currentView.Name
	}

	// When cluster info 、navigation or detail panel selected.
	if !resizeAbleView(resizeView) {
		resizeView = gui.PeekPreviousView()

		// If previous view is cluster info 、navigation or detail pane.
		if !resizeAbleView(resizeView) {
			resizeView = resizeableViews[0]
		}
	}

	n := len(resizeableViews)
	maxHeight := gui.MaxHeight()
	height := viewHeights[view.Name]
	if resizeView == view.Name {
		height += usableSpace(gui, maxHeight) % n
	}

	return height
}

func migrateTopFunc(gui *guilib.Gui, view *guilib.View) int {
	maxHeight := gui.MaxHeight()
	usedSpace := clusterInfoViewHeight - optionViewHeight - (len(functionViews) - 1)
	HeightLine := resizeAbleViewMinHeight*len(resizeableViews) - usedSpace
	if maxHeight < HeightLine {
		currentView := gui.CurrentView()
		if currentView != nil {
			if !resizeAbleView(view.Name) {
				if view.Name == resizeableViews[0] {
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
		return -1
	}
	return 1
}
