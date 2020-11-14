package app

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
	"math"
)

const (
	selectedViewLine = "selectedViewLine"
)

var (
	cyclicViews = []string{clusterInfoViewName, namespaceViewName, serviceViewName, deploymentViewName, podViewName}

	nextCyclicView = &gui.Action{
		Name:    "nextCyclicView",
		Keys:    []interface{}{gocui.KeyTab},
		Handler: nextCyclicViewHandler,
		Mod:     gocui.ModNone,
	}

	//previousCyclicView = &gui.Action{
	//	Name:    "previousCyclicView",
	//	Key:     gocui.KeyArrowUp,
	//	Handler: previousCyclicViewHandler,
	//	Mod:     gocui.ModNone,
	//}

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

	previousLine = &gui.Action{
		Name:    "previousLine",
		Key:     gocui.KeyArrowUp,
		Handler: previousLineHandler,
		Mod:     gocui.ModNone,
	}

	nextLine = &gui.Action{
		Name:    "nextLine",
		Key:     gocui.KeyArrowDown,
		Handler: nextLineHandler,
		Mod:     gocui.ModNone,
	}

	actions = []*gui.Action{
		backToPreviousView,
		&gui.Action{
			Name: "scrollUp",
			Keys: []interface{}{
				gocui.KeyPgup,
				gocui.MouseWheelUp,
			},
			Handler: scrollUpHandler,
			Mod:     gocui.ModNone,
		},
		&gui.Action{
			Name: "scrollDown",
			Keys: []interface{}{
				gocui.KeyPgdn,
				gocui.MouseWheelDown,
			},
			Handler: scrollDownHandler,
			Mod:     gocui.ModNone,
		},
		&gui.Action{
			Name:    "scrollTop",
			Key:     gocui.KeyHome,
			Handler: scrollTopHandler,
			Mod:     gocui.ModNone,
		},
		&gui.Action{
			Name:    "scrollBottom",
			Key:     gocui.KeyEnd,
			Handler: scrollBottomHandler,
			Mod:     gocui.ModNone,
		},
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
			//previousCyclicView,
		},
		namespaceViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			//previousCyclicView,
			previousLine,
			nextLine,
		},
		serviceViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			//previousCyclicView,
			previousLine,
			nextLine,
		},
		deploymentViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			//previousCyclicView,
			previousLine,
			nextLine,
		},
		podViewName: []*gui.Action{
			toNavigation,
			nextCyclicView,
			//previousCyclicView,
			previousLine,
			nextLine,
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

		return gui.FocusView(clusterInfoViewName, false)
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
		if navigationIndex+1 >= len(options) {
			return nil
		}
		switchNavigation(navigationIndex + 1)
		return nil
	}
}

func navigationArrowLeftHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		if navigationIndex-1 < 0 {
			return gui.ReturnPreviousView()
		}
		switchNavigation(navigationIndex - 1)
		return nil
	}
}

func scrollUpHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, oy := v.Origin()
		newOy := int(math.Max(0, float64(oy-2)))
		return v.SetOrigin(ox, newOy)
	}
}

func scrollDownHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, oy := v.Origin()

		reservedLines := 0
		_, sizeY := v.Size()
		reservedLines = sizeY

		totalLines := len(v.ViewBufferLines())
		if oy+reservedLines >= totalLines {
			v.Autoscroll = true
			return nil
		}

		return v.SetOrigin(ox, oy+2)
	}
}

func scrollTopHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, _ := v.Origin()
		return v.SetOrigin(ox, 0)
	}
}

func scrollBottomHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		totalLines := len(v.ViewBufferLines())
		if totalLines == 0 {
			return nil
		}
		_, vy := v.Size()
		if totalLines <= vy {
			return nil
		}

		ox, _ := v.Origin()
		v.Autoscroll = true
		return v.SetOrigin(ox, totalLines-1)
	}
}

func viewLineClickHandler(gui *gui.Gui, view *gui.View, cy int, lineString string) error {
	if cy == 0 {
		selected := formatSelectedName(lineString, 0)
		if selected == "NAME" || selected == "NAMESPACE" {
			log.Logger.Debugf("viewLineClickHandler - view: '%s' cy == 0, view.State.Set(selectedViewLine, nil)", view.Name)
			if view.Name == namespaceViewName {
				kubecli.Cli.SetNamespace("")
			}
			return view.State.Set(selectedViewLine, nil)
		}
	}

	log.Logger.Debugf("viewLineClickHandler - view: '%s' view.State.Set(selectedViewLine, \"%s\")", view.Name, lineString)
	if view.Name == namespaceViewName {
		namespace := formatSelectedNamespace(lineString)
		log.Logger.Debugf("viewLineClickHandler - switch namespace to %s", namespace)
		kubecli.Cli.SetNamespace(namespace)
	}
	return view.State.Set(selectedViewLine, lineString)
}

func previousLineHandler(gui *gui.Gui) func(gui *gocui.Gui, view *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		currentView := gui.CurrentView()
		if currentView == nil {
			return nil
		}

		_, cy := v.Cursor()
		if cy-1 < 0 {
			return nil
		}
		lineStr, err := v.Line(cy - 1)
		if err != nil {
			log.Logger.Warningf("previousLineHandler - v.Line(cy - 1)", cy)
		}
		v.MoveCursor(0, -1, false)

		if cy-1 != 0 {
			if currentView.Name == namespaceViewName {
				namespace := formatSelectedNamespace(lineStr)
				log.Logger.Debugf("previousLineHandler - switch namespace to %s", namespace)
				kubecli.Cli.SetNamespace(namespace)
			}
			return currentView.State.Set(selectedViewLine, lineStr)
		}
		if currentView.Name == namespaceViewName {
			namespace := ""
			log.Logger.Debugf("previousLineHandler - switch namespace to %s", namespace)
			kubecli.Cli.SetNamespace(namespace)
		}
		return currentView.State.Set(selectedViewLine, lineStr)
	}
}

func nextLineHandler(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		currentView := gui.CurrentView()
		if currentView == nil {
			return nil
		}

		_, cy := v.Cursor()
		if cy+1 >= len(v.ViewBufferLines()) {
			return nil
		}
		lineStr, err := v.Line(cy + 1)
		if err != nil {
			log.Logger.Warningf("nextLineHandler - v.Line(%d + 1)", cy)
		}
		v.MoveCursor(0, 1, false)

		if currentView.Name == namespaceViewName {
			namespace := formatSelectedNamespace(lineStr)
			log.Logger.Debugf("nextLineHandler - switch namespace to %s", namespace)
			kubecli.Cli.SetNamespace(namespace)
		}
		return currentView.State.Set(selectedViewLine, lineStr)
	}
}
