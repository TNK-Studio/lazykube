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

func usableSpace(maxHeight int) int {
	return maxHeight - 4
}

func reactiveHeight(gui *gui.Gui, view *gui.View) int {
	_, maxHeight := gui.Size()

	space := usableSpace(maxHeight)

	tallPanels := 4
	viewHeights = map[string]int{
		clusterInfoViewName: 3,
		namespaceViewName:   space / tallPanels,
		serviceViewName:     space / tallPanels,
		deploymentViewName:  space / tallPanels,
		podViewName:         space / tallPanels,
		optionViewName:      1,
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
			clusterInfoViewName: defaultHeight,
			namespaceViewName:   defaultHeight,
			serviceViewName:     defaultHeight,
			deploymentViewName:  defaultHeight,
			podViewName:         defaultHeight,
			optionViewName:      defaultHeight,
		}

		viewHeights[resizeView] = maxHeight - defaultHeight*tallPanels - 1
	}

	viewHeights[clusterInfoViewName] -= 1
	if viewHeights[clusterInfoViewName] == 0 {
		viewHeights[clusterInfoViewName] = 1
	}
	height := viewHeights[view.Name]
	return height
}

func migrateTopFunc(gui *gui.Gui, view *gui.View) int {
	return 1
}
