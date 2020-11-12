package app

import (
	"fmt"
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/gookit/color"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"strings"
)

const (
	OptSeparator       = "   "
	navigationPathJoin = " + "
)

var (
	// Todo: use state to control.
	activeView *gui.View

	navigationIndex     int
	activeNavigationOpt string

	functionViews     = []string{clusterInfoViewName, namespaceViewName, serviceViewName, deploymentViewName, podViewName}
	viewNavigationMap = map[string][]string{
		clusterInfoViewName: []string{"Nodes", "Top Nodes"},
		namespaceViewName:   []string{"Deployments", "Pods", "Config"},
		serviceViewName:     []string{"Pods Log", "Config"},
		deploymentViewName:  []string{"Pods Log", "Config", "Describe", "Top Pods"},
		podViewName:         []string{"Log", "Top", "Config", "Describe"},
	}

	detailRenderMap = map[string]func(gui *gui.Gui, view *gui.View) error{
		navigationPath(clusterInfoViewName, "Nodes"):     clusterNodesRender,
		navigationPath(clusterInfoViewName, "Top Nodes"): topNodesRender,
	}
)

func navigationPath(args ...string) string {
	return strings.Join(args, navigationPathJoin)
}

func switchNavigation(index int) string {
	if index < 0 {
		return ""
	}

	if activeView != nil {
		if index >= len(viewNavigationMap[activeView.Name]) {
			return ""
		}
		navigationIndex = index
		activeNavigationOpt = viewNavigationMap[activeView.Name][index]
		return activeNavigationOpt
	}
	return ""
}

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
			if err := gui.FocusView(namespaceViewName, false); err != nil {
				log.Logger.Println(err)
			}
		}
		activeView = gui.CurrentView()
	}

	options := viewNavigationMap[activeView.Name]
	if changeNavigation {
		switchNavigation(0)
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
			selected = switchNavigation(i)
			break
		}
	}

	log.Logger.Debugf("navigationOnClick - selected '%s'", selected)

	return nil
}

func renderClusterInfo(gui *gui.Gui, view *gui.View) error {
	view.Clear()
	currentContext, err := kubecli.Cli.CurrentContext()
	if err != nil {
		return nil
	}

	if _, err := fmt.Fprintf(view, "Current Context: %s", color.Green.Sprint(currentContext)); err != nil {
		return err
	}

	//clusterInfo, err := kubecli.Cli.ClusterInfo()
	//if err != nil {
	//	return nil
	//}
	//if _, err := fmt.Fprintln(view, clusterInfo); err != nil {
	//	return err
	//}

	return nil
}

func detailRender(gui *gui.Gui, view *gui.View) error {
	view.Clear()
	if activeView == nil {
		return nil
	}
	renderFunc := detailRenderMap[navigationPath(activeView.Name, activeNavigationOpt)]
	if renderFunc != nil {
		return renderFunc(gui, view)
	}
	return nil
}

func viewStreams(view *gui.View) genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    view,
		ErrOut: view,
	}
}

func clusterNodesRender(gui *gui.Gui, view *gui.View) error {
	kubecli.Cli.Get(viewStreams(view), "nodes")
	return nil
}

func topNodesRender(gui *gui.Gui, view *gui.View) error {
	kubecli.Cli.TopNode(viewStreams(view), nil, "")
	return nil
}

func namespaceRender(gui *gui.Gui, view *gui.View) error {
	view.Clear()
	kubecli.Cli.Get(viewStreams(view), "namespaces")
	return nil
}
