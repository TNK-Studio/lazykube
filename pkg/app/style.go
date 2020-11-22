package app

import (
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"strings"
)

const (
	resizeableViewMinHeight = 5
	clusterInfoViewHeight   = 2
	optionViewHeight        = 1
)

var (
	viewHeights     = map[string]int{}
	functionViews   = []string{clusterInfoViewName, namespaceViewName, serviceViewName, deploymentViewName, podViewName}
	resizeableViews = []string{namespaceViewName, serviceViewName, deploymentViewName, podViewName}

	// Function cache
	reactiveHeightCache = map[string]int{}
	migrateTopCache     = map[string]int{}
)

func resizeableView(viewName string) bool {
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

func usableSpace(maxHeight int) int {
	usedSpace := clusterInfoViewHeight + optionViewHeight + (len(functionViews) - 1)
	if maxHeight < resizeableViewMinHeight*len(resizeableViews)-usedSpace {
		return maxHeight - 2
	}
	return maxHeight - usedSpace
}

func resizePanelHeight(gui *guilib.Gui) error {
	_, maxHeight := gui.Size()

	space := usableSpace(maxHeight)

	n := len(resizeableViews)
	viewHeights[clusterInfoViewName] = clusterInfoViewHeight
	for _, resizeableView := range resizeableViews {
		viewHeights[resizeableView] = space / n
	}
	viewHeights[optionViewName] = optionViewHeight

	return nil
}

// Todo: cache result
func reactiveHeight(gui *guilib.Gui, view *guilib.View) int {
	var currentViewName string
	currentView := gui.CurrentView()
	if currentView != nil {
		currentViewName = currentView.Name
	}
	previousViewName := gui.PeekPreviousView()

	height := cacheAbleReactiveHeight(gui.MaxHeight(), resizeableViews, currentViewName, previousViewName, view.Name)
	return height
}

func cacheAbleReactiveHeight(maxHeight int, resizeableViews []string, currentViewName, previousViewName, viewName string) int {
	key := fmt.Sprintf("%d,%s,%s,%s,%s", maxHeight, strings.Join(resizeableViews, ","), currentViewName, previousViewName, viewName)
	cacheVal, ok := reactiveHeightCache[key]
	if ok {
		return cacheVal
	}

	var resizeView string
	if currentViewName == "" {
		resizeView = resizeableViews[0]
	} else {
		resizeView = currentViewName
	}

	// When cluster info 、navigation or detail panel selected.
	if !resizeableView(resizeView) {
		resizeView = previousViewName

		// If previous view is cluster info 、navigation or detail pane.
		if !resizeableView(resizeView) {
			resizeView = resizeableViews[0]
		}
	}

	n := len(resizeableViews)
	height := viewHeights[viewName]

	if maxHeight < heightBoundary() && resizeableView(viewName) {
		if resizeView == viewName {
			cacheVal = height + len(resizeableViews)
		} else {
			cacheVal = 2
		}
	} else {
		if resizeView == viewName {
			height += usableSpace(maxHeight) % n
		}
	}

	cacheVal = height
	reactiveHeightCache[key] = cacheVal
	return cacheVal
}

func migrateTopFunc(gui *guilib.Gui, view *guilib.View) int {
	var currentViewName string
	currentView := gui.CurrentView()
	if currentView != nil {
		currentViewName = currentView.Name
	}
	return cacheAbleMigrateTopFunc(gui.MaxHeight(), resizeableViews, currentViewName, view.Name)
}

func cacheAbleMigrateTopFunc(maxHeight int, resizeableViews []string, currentViewName, viewName string) int {
	key := fmt.Sprintf("%d,%s,%s,%s", maxHeight, strings.Join(resizeableViews, ","), currentViewName, viewName)
	cacheVal, ok := migrateTopCache[key]
	if ok {
		return cacheVal
	}

	if maxHeight < heightBoundary() {
		if currentViewName != "" {
			if !resizeableView(viewName) {
				if viewName == resizeableViews[0] {
					cacheVal = 1
				}
			} else {
				for i, viewName := range functionViews {
					if currentViewName == viewName {
						index := i + 1
						if index < len(functionViews) && functionViews[index] == viewName {
							cacheVal = 1
							break
						}
					}
				}
			}
		} else {
			cacheVal = -1
		}
	}
	cacheVal = 1
	migrateTopCache[key] = cacheVal
	return cacheVal
}

func heightBoundary() int {
	usedSpace := clusterInfoViewHeight - optionViewHeight - (len(functionViews) - 1)
	boundary := resizeableViewMinHeight*len(resizeableViews) - usedSpace
	return boundary
}
