package app

import (
	"fmt"
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/gookit/color"
	"strings"
)

var (
	// Todo: use state to control.
	activeView *gui.View

	navigationIndex int

	functionViews     = []string{"clusterInfo", "namespace", "service", "deployment", "pod"}
	viewNavigationMap = map[string][]string{
		"clusterInfo": []string{"Nodes", "Top Nodes"},
		"namespace":   []string{"Config", "Deployments", "Pods"},
		"service":     []string{"Config", "Pods Log"},
		"deployment":  []string{"Config", "Describe", "Pods Log", "Top Pods"},
		"pod":         []string{"Config", "Describe", "Log", "Top"},
	}

	OptSeparator = "   "
)

func navigationRender(gui *gui.Gui, view *gui.View) error {
	currentView := gui.CurrentView()
	var changeNavigation bool
	if currentView != nil {
		for _, viewName := range functionViews {
			if currentView.Name == viewName {
				if activeView != currentView {
					changeNavigation = true
				}
				activeView = currentView
				break
			}
		}
	}

	if activeView == nil {
		if gui.CurrentView() == nil {
			if err := gui.FocusView("namespace", false); err != nil {
				log.Logger.Println(err)
			}
		}
		activeView = gui.CurrentView()
	}

	options := viewNavigationMap[activeView.Name]
	if changeNavigation {
		navigationIndex = 0
	}

	colorfulOptions := make([]string, 0)
	for index, opt := range options {
		colorfulOpt := color.White.Sprint(opt)
		if navigationIndex == index {
			colorfulOpt = color.Green.Sprint(opt)
		}
		colorfulOptions = append(colorfulOptions, colorfulOpt)
	}

	view.Clear()
	str := strings.Join(colorfulOptions, OptSeparator)
	fmt.Fprint(view, str)

	return nil
}

func navigationOnClick(gui *gui.Gui, view *gui.View) error {
	cx, cy := view.Cursor()
	log.Logger.Debugf("navigationOnClick - cx %d cy %d", cx, cy)

	options := viewNavigationMap[activeView.Name]
	sep := len(OptSeparator)
	halfSep := sep / 2
	preFix := 0

	var selected string
	for i, opt := range options {
		left := preFix + i*sep

		words := len([]rune(opt))

		right := left + words - 1
		preFix += words - 1

		if cx >= left-halfSep && cx <= right+halfSep {
			log.Logger.Debugf("navigationOnClick - cx %d in selection[%d, %d]", cx, left, right)
			navigationIndex = i
			selected = options[i]
			break
		}
	}

	log.Logger.Debugf("navigationOnClick - selected '%s'", selected)

	return nil
}
