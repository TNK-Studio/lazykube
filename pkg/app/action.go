package app

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
)

var (
	cyclicViews = []string{clusterInfoViewName, namespaceViewName, serviceViewName, deploymentViewName, podViewName}

	nextCyclicView = &gui.Action{
		Name:    "nextCyclicView",
		Keys:    []interface{}{gocui.KeyTab, gocui.KeyArrowDown},
		Handler: nextCyclicViewHandler,
		Mod:     gocui.ModNone,
	}

	previousCyclicView = &gui.Action{
		Name:    "previousCyclicView",
		Key:     gocui.KeyArrowUp,
		Handler: previousCyclicViewHandler,
		Mod:     gocui.ModNone,
	}

	backToPreviousView = &gui.Action{
		Name:    "backToPreviousView",
		Key:     gocui.KeyEsc,
		Handler: backToPreviousViewHandler,
		Mod:     gocui.ModNone,
	}

	toNavigation = &gui.Action{
		Name: "toNavigation",
		Keys: []interface{}{
			gocui.KeyEnter,
			gocui.KeyArrowRight,
		},
		Handler: toNavigationHandler,
		Mod:     gocui.ModNone,
	}

	actions = []*gui.Action{
		backToPreviousView,
	}

	viewActionsMap = map[string][]*gui.Action{
		navigationViewName: []*gui.Action{
			&gui.Action{
				Name:    "navigationArrowLeft",
				Key:     gocui.KeyArrowLeft,
				Handler: navigationArrowLeftHandler,
				Mod:     gocui.ModNone,
			},
			&gui.Action{
				Name:    "navigationArrowRight",
				Key:     gocui.KeyArrowRight,
				Handler: navigationArrowRightHandler,
				Mod:     gocui.ModNone,
			},
			&gui.Action{
				Name: "navigationDown",
				Key:  gocui.KeyArrowDown,
				Handler: func(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
					return func(*gocui.Gui, *gocui.View) error {
						gui.FocusView(detailViewName, false)
						return nil
					}
				},
				Mod: gocui.ModNone,
			},
		},
		detailViewName: []*gui.Action{
			&gui.Action{
				Name: "detailArrowUp",
				Key:  gocui.KeyArrowUp,
				Handler: func(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
					return func(*gocui.Gui, *gocui.View) error {
						gui.FocusView(navigationViewName, false)
						return nil
					}
				},
				Mod: gocui.ModNone,
			},
		},
		clusterInfoViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			previousCyclicView,
		},
		namespaceViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			previousCyclicView,
		},
		serviceViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			previousCyclicView,
		},
		deploymentViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			previousCyclicView,
		},
		podViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			previousCyclicView,
		},
	}
)

func nextCyclicViewHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {

		currentView := gui.CurrentView()
		if currentView == nil {
			return nil
		}

		for index, viewName := range cyclicViews {
			if currentView.Name == viewName {
				nextIndex := index + 1
				if nextIndex >= len(cyclicViews) {
					nextIndex = 0
				}
				nextViewName := cyclicViews[nextIndex]
				log.Logger.Debugf("nextCyclicViewHandler - nextViewName: %s", nextViewName)
				return gui.FocusView(nextViewName, true)
			}
		}
		return nil
	}
}

func previousCyclicViewHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {

		currentView := gui.CurrentView()
		if currentView == nil {
			return nil
		}

		for index, viewName := range cyclicViews {
			if currentView.Name == viewName {
				nextIndex := index - 1
				if nextIndex < 0 {
					nextIndex = len(cyclicViews) - 1
				}
				previousViewName := cyclicViews[nextIndex]
				log.Logger.Debugf("previousCyclicViewHandler - previousViewName: %s", previousViewName)
				return gui.FocusView(cyclicViews[nextIndex], true)
			}
		}
		return nil
	}
}

func backToPreviousViewHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		if gui.HasPreviousView() {
			return gui.ReturnPreviousView()
		}

		return gui.FocusView(namespaceViewName, false)
	}
}

func toNavigationHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		return gui.FocusView(navigationViewName, true)
	}
}

func navigationArrowRightHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		options := viewNavigationMap[activeView.Name]
		navigationIndex += 1
		if navigationIndex >= len(options) {
			navigationIndex = len(options) - 1
		}
		return nil
	}
}

func navigationArrowLeftHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		navigationIndex -= 1
		if navigationIndex < 0 {
			navigationIndex = 0
			return gui.ReturnPreviousView()
		}
		return nil
	}
}
