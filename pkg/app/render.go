package app

import (
	"errors"
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/gookit/color"
	"io"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"strings"
	"time"
)

const (
	optSeparator       = "   "
	navigationPathJoin = " + "
	logsTail           = "500"

	namespaceResource  = "namespace"
	serviceResource    = "service"
	deploymentResource = "deployment"
	podResource        = "pod"

	reRenderIntervalDuration = 3 * time.Second
)

var (
	// Todo: use state to control.
	activeView *guilib.View

	navigationIndex     int
	activeNavigationOpt string

	navigationOptNodes       = "Nodes"
	navigationOptTopNodes    = "Top Nodes"
	navigationOptDeployments = "Deployments"
	navigationOptPods        = "Pods"
	navigationOptPodsLog     = "Pods Log"
	navigationOptTopPods     = "Top Pods"
	navigationOptServices    = "Services"
	navigationOptConfig      = "Config"
	navigationOptDescribe    = "Describe"
	navigationOptTop         = "Top"
	navigationOptLog         = "Log"

	viewNavigationMap = map[string][]string{
		clusterInfoViewName: {navigationOptNodes, navigationOptTopNodes},
		namespaceViewName:   {navigationOptConfig, navigationOptServices, navigationOptDeployments, navigationOptPods},
		serviceViewName:     {navigationOptConfig, navigationOptPods, navigationOptPodsLog, navigationOptTopPods},
		deploymentViewName:  {navigationOptConfig, navigationOptDescribe, navigationOptPods, navigationOptPodsLog, navigationOptTopPods},
		podViewName:         {navigationOptLog, navigationOptConfig, navigationOptDescribe, navigationOptTop},
	}

	detailRenderMap = map[string]guilib.ViewHandler{
		navigationPath(clusterInfoViewName, navigationOptNodes):     reRenderInterval(clearBeforeRender(clusterNodesRender), reRenderIntervalDuration),
		navigationPath(clusterInfoViewName, navigationOptTopNodes):  reRenderInterval(clearBeforeRender(topNodesRender), reRenderIntervalDuration),
		navigationPath(namespaceViewName, navigationOptDeployments): reRenderInterval(clearBeforeRender(namespaceResourceListRender("deployments")), reRenderIntervalDuration),
		navigationPath(namespaceViewName, navigationOptPods):        reRenderInterval(clearBeforeRender(namespaceResourceListRender("pods")), reRenderIntervalDuration),
		navigationPath(namespaceViewName, navigationOptServices):    reRenderInterval(clearBeforeRender(namespaceResourceListRender("services")), reRenderIntervalDuration),
		navigationPath(namespaceViewName, navigationOptConfig):      reRenderInterval(clearBeforeRender(configRender), reRenderIntervalDuration),
		navigationPath(serviceViewName, navigationOptConfig):        reRenderInterval(clearBeforeRender(configRender), reRenderIntervalDuration),
		navigationPath(serviceViewName, navigationOptPods):          reRenderInterval(clearBeforeRender(labelsPodsRender), reRenderIntervalDuration),
		navigationPath(serviceViewName, navigationOptPodsLog):       reRenderInterval(podsLogsRender, reRenderIntervalDuration),
		navigationPath(serviceViewName, navigationOptTopPods):       reRenderInterval(clearBeforeRender(topPodsRender), reRenderIntervalDuration),
		navigationPath(deploymentViewName, navigationOptConfig):     reRenderInterval(clearBeforeRender(configRender), reRenderIntervalDuration),
		navigationPath(deploymentViewName, navigationOptPods):       reRenderInterval(clearBeforeRender(labelsPodsRender), reRenderIntervalDuration),
		navigationPath(deploymentViewName, navigationOptDescribe):   reRenderInterval(clearBeforeRender(describeRender), reRenderIntervalDuration),
		navigationPath(deploymentViewName, navigationOptPodsLog):    reRenderInterval(podsLogsRender, reRenderIntervalDuration),
		navigationPath(deploymentViewName, navigationOptTopPods):    reRenderInterval(clearBeforeRender(topPodsRender), reRenderIntervalDuration),
		navigationPath(podViewName, navigationOptConfig):            reRenderInterval(clearBeforeRender(configRender), reRenderIntervalDuration),
		navigationPath(podViewName, navigationOptLog):               reRenderInterval(podLogsRender, reRenderIntervalDuration),
		navigationPath(podViewName, navigationOptDescribe):          reRenderInterval(clearBeforeRender(describeRender), reRenderIntervalDuration),
		navigationPath(podViewName, navigationOptTop):               reRenderInterval(podMetricsPlotRender, reRenderIntervalDuration),
	}
)

func notResourceSelected(selectedName string) bool {
	if selectedName == "" || selectedName == "NAME" || selectedName == "NAMESPACE" || selectedName == "No" {
		return true
	}
	return false
}

func clearBeforeRender(render guilib.ViewHandler) guilib.ViewHandler {
	return func(gui *guilib.Gui, view *guilib.View) error {
		view.Clear()
		return render(gui, view)
	}
}

func reRenderInterval(handler guilib.ViewHandler, interval time.Duration) guilib.ViewHandler {
	return func(gui *guilib.Gui, view *guilib.View) error {
		now := time.Now()
		view.ReRender()
		val, _ := view.GetState(viewLastRenderTimeStateKey)
		if val == nil {
			if err := view.SetState(viewLastRenderTimeStateKey, now, true); err != nil {
				return err
			}
			log.Logger.Debugf("reRenderInterval - interval: %+v handler %+v view %s", interval, handler, view.Name)
			if err := handler(gui, view); err != nil {
				return nil
			}
			return nil
		}

		viewLastRenderTime := val.(time.Time)
		if viewLastRenderTime.Add(interval).After(now) {
			return nil
		}
		if err := view.SetState(viewLastRenderTimeStateKey, now, true); err != nil {
			return err
		}
		log.Logger.Debugf("reRenderInterval - interval: %+v handler %+v view %s", interval, handler, view.Name)
		if err := handler(gui, view); err != nil {
			return nil
		}
		return nil
	}
}

func clearLastRenderTime(gui *guilib.Gui, viewName string) error {
	view, err := gui.GetView(viewName)
	if err != nil {
		return err
	}
	if err := view.SetState(viewLastRenderTimeStateKey, nil, true); err != nil {
		return err
	}
	return nil
}

func clearDetailViewState(gui *guilib.Gui) {
	detailView, err := gui.GetView(detailViewName)
	if err != nil {
		log.Logger.Warningf("clearDetailViewState - get view error %s", err)
		return
	}

	if err := clearLastRenderTime(gui, detailViewName); err != nil {
		log.Logger.Warningf("clearDetailViewState - clearLastRenderTime err %s", err)
		return
	}

	if err := detailView.SetState(logSinceTimeStateKey, nil, true); err != nil {
		log.Logger.Warningf("clearDetailViewState - clear logSinceTimeStateKey err %s", err)
		return
	}

	if err := detailView.SetState(logContainerStateKey, nil, true); err != nil {
		log.Logger.Warningf("clearDetailViewState - clear logContainerStateKey err %s", err)
		return
	}
	_ = detailView.SetOrigin(0, 0)
	_ = detailView.SetCursor(0, 0)
	detailView.Clear()
}

func navigationPath(args ...string) string {
	return strings.Join(args, navigationPathJoin)
}

func switchNavigation(gui *guilib.Gui, index int) string {
	err := Detail.SetOrigin(0, 0)
	if err != nil {
		log.Logger.Warningf("switchNavigation - Detail.SetOrigin(0, 0) error %s", err)
	}

	Detail.Clear()
	clearDetailViewState(gui)
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

func navigationRender(gui *guilib.Gui, view *guilib.View) error {
	currentView := gui.CurrentView()
	// Change navigation render
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
			if err := gui.FocusView(functionViews[0], false); err != nil {
				log.Logger.Println(err)
			}
		}
		activeView = gui.CurrentView()
	}

	options := viewNavigationMap[activeView.Name]
	if activeNavigationOpt == "" {
		activeNavigationOpt = options[navigationIndex]
	}
	if changeNavigation {
		switchNavigation(gui, 0)
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
	str := strings.Join(colorfulOptions, optSeparator)

	_, err := fmt.Fprint(view, str)
	if err != nil {
		return err
	}

	return nil
}

func navigationOnClick(gui *guilib.Gui, view *guilib.View) error {
	cx, cy := view.Cursor()
	log.Logger.Debugf("navigationOnClick - cx %d cy %d", cx, cy)

	options := viewNavigationMap[activeView.Name]
	optionIndex, selected := utils.ClickOption(options, optSeparator, cx, 0)
	if optionIndex < 0 {
		return nil
	}
	log.Logger.Debugf("navigationOnClick - cx %d selected '%s'", cx, selected)
	_ = switchNavigation(gui, optionIndex)
	view.ReRender()
	Detail.ReRender()
	return nil
}

func renderClusterInfo(_ *guilib.Gui, view *guilib.View) error {
	view.Clear()
	currentContext := kubecli.Cli.CurrentContext()

	if _, err := fmt.Fprintf(view, "Current Context: %s", color.Green.Sprint(currentContext)); err != nil {
		return err
	}
	return nil
}

func detailRender(gui *guilib.Gui, view *guilib.View) error {
	if activeView == nil {
		return nil
	}
	renderFunc := detailRenderMap[navigationPath(activeView.Name, activeNavigationOpt)]
	if renderFunc != nil {
		return renderFunc(gui, view)
	}
	return nil
}

func viewStreams(view *guilib.View) genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    view,
		ErrOut: view,
	}
}

func clusterNodesRender(_ *guilib.Gui, view *guilib.View) error {
	kubecli.Cli.Get(viewStreams(view), navigationOptNodes).Run()
	return nil
}

func topNodesRender(_ *guilib.Gui, view *guilib.View) error {
	kubecli.Cli.TopNode(viewStreams(view), nil, "").Run()
	view.ReRender()
	return nil
}

func namespaceRender(_ *guilib.Gui, view *guilib.View) error {
	view.Clear()
	kubecli.Cli.Get(viewStreams(view), "namespaces").Run()
	return nil
}

func namespaceResourceListRender(resource string) guilib.ViewHandler {
	return func(gui *guilib.Gui, view *guilib.View) error {
		view.Clear()
		if kubecli.Cli.Namespace() == "" {
			kubecli.Cli.Get(viewStreams(view), resource).SetFlag("all-namespaces", "true").SetFlag("output", "wide").Run()
			return nil
		}
		kubecli.Cli.Get(viewStreams(view), resource).SetFlag("output", "wide").Run()
		return nil
	}
}

func resourceListRender(_ *guilib.Gui, view *guilib.View) error {
	view.Clear()
	resource := getViewResourceName(view.Name)
	if kubecli.Cli.Namespace() == "" {
		kubecli.Cli.Get(viewStreams(view), resource).SetFlag("all-namespaces", "true").Run()
		return nil
	}
	kubecli.Cli.Get(viewStreams(view), resource).Run()
	return nil
}

func showPleaseSelected(view io.Writer, name string) {
	_, err := fmt.Fprintf(view, "Please select a %s.\n ", name)
	if err != nil {
		log.Logger.Warningf("showPleaseSelected - error %s", err)
	}
}

func namespaceConfigRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	namespaceView, err := gui.GetView(namespaceViewName)
	if err != nil {
		return nil
	}
	namespace := formatSelectedNamespace(namespaceView.SelectedLine)
	if notResourceSelected(namespace) {
		showPleaseSelected(view, namespaceViewName)
		return nil
	}

	kubecli.Cli.Get(viewStreams(view), "namespaces", namespace).SetFlag("output", "yaml").Run()
	return nil
}

func configRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if activeView == nil {
		return nil
	}

	namespaceView, err := gui.GetView(namespaceViewName)
	if err != nil {
		return nil
	}

	if activeView == namespaceView {
		return namespaceConfigRender(gui, view)
	}

	resource := getViewResourceName(activeView.Name)

	if resource == "" {
		return nil
	}

	namespace, resourceName, err := getResourceNamespaceAndName(gui, activeView)
	if err != nil {
		if errors.Is(err, noResourceSelectedErr) {
			showPleaseSelected(view, resource)
			return nil
		}
		return err
	}

	cli(namespace).Get(viewStreams(view), resource, resourceName).SetFlag("output", "yaml").Run()
	return nil
}

func describeRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if activeView == nil {
		return nil
	}
	if activeView.Name == namespaceViewName {
		return namespaceConfigRender(gui, view)
	}

	resource := getViewResourceName(activeView.Name)

	if resource == "" {
		return nil
	}

	namespace, resourceName, err := getResourceNamespaceAndName(gui, activeView)
	if err != nil {
		if errors.Is(err, noResourceSelectedErr) {
			showPleaseSelected(view, resource)
			return nil
		}

		return err
	}

	cli(namespace).Describe(viewStreams(view), resource, resourceName).Run()

	view.ReRender()
	return nil
}

func onFocusClearSelected(gui *guilib.Gui, view *guilib.View) error {
	for _, functionViewName := range functionViews {
		if functionViewName == view.Name || functionViewName == namespaceViewName {
			continue
		}
		functionView, err := gui.GetView(functionViewName)
		if err != nil {
			log.Logger.Warningf("onFocusClearSelected - view name %s gui.GetView(\"%s\") error %s", view.Name, functionView, err)
			continue
		}
		if err := functionView.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := functionView.SetCursor(0, 0); err != nil {
			return err
		}
	}
	return nil
}

func podLogsRender(gui *guilib.Gui, view *guilib.View) error {
	// Todo: Fix chinese character of logs.
	scrollLogs := true
	if val, _ := view.GetState(ScrollingLogsStateKey); val != nil {
		var ok bool
		scrollLogs, ok = val.(bool)
		if !ok {
			scrollLogs = true
		}
	}

	if !scrollLogs {
		return nil
	}

	podView, err := gui.GetView(podViewName)
	if err != nil {
		return err
	}

	resource := "pod"
	namespace, resourceName, err := getResourceNamespaceAndName(gui, podView)
	if err != nil {
		if errors.Is(err, noResourceSelectedErr) {
			showPleaseSelected(view, resource)
			return nil
		}
		return err
	}

	containers := getPodContainers(namespace, resourceName)

	if err := view.SetState(podContainersStateKey, containers, true); err != nil {
		return err
	}

	var since time.Time
	var hasSince bool
	if val, _ := view.GetState(logSinceTimeStateKey); val != nil {
		hasSince = true
		since = val.(time.Time)
	}

	var logContainer string
	if val, _ := view.GetState(logContainerStateKey); val != nil {
		logContainer = val.(string)
	}

	cmd := cli(namespace).
		Logs(viewStreams(view), resourceName).
		SetFlag("tail", logsTail).
		SetFlag("prefix", "true")

	if logContainer == "" {
		cmd.SetFlag("all-containers", "true")
	} else {
		cmd.SetFlag("container", logContainer)
	}

	if hasSince {
		cmd.SetFlag("since-time", since.Format(time.RFC3339))
	}

	cmd.Run()

	if err := view.SetState(logSinceTimeStateKey, time.Now(), true); err != nil {
		return err
	}

	view.ReRender()
	return nil
}
func podsLogsRender(gui *guilib.Gui, view *guilib.View) error {
	// Todo: Fix chinese character of logs.
	if err := podsSelectorRenderHelper(func(namespace string, labelsArr []string) error {
		var since time.Time
		var hasSince bool
		val, _ := view.GetState(logSinceTimeStateKey)
		if val != nil {
			hasSince = true
			since = val.(time.Time)
		}

		streams := newStream()
		cmd := kubecli.Cli.WithNamespace(namespace).Logs(streams)
		cmd.SetFlag("selector", strings.Join(labelsArr, ",")).
			SetFlag("all-containers", "true").
			SetFlag("tail", logsTail).
			SetFlag("prefix", "true")

		if hasSince {
			cmd.SetFlag("since-time", since.Format(time.RFC3339))
		}

		cmd.Run()

		if err := view.SetState(logSinceTimeStateKey, time.Now(), true); err != nil {
			return err
		}
		streamCopyTo(streams, view)
		view.ReRender()
		return nil
	})(gui, view); err != nil {
		return err
	}
	return nil
}

func labelsPodsRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if err := podsSelectorRenderHelper(func(namespace string, labelsArr []string) error {
		cmd := kubecli.Cli.WithNamespace(namespace).Get(viewStreams(view), "pods")
		cmd.SetFlag("selector", strings.Join(labelsArr, ","))
		cmd.SetFlag("output", "wide")
		cmd.Run()
		view.ReRender()
		return nil
	})(gui, view); err != nil {
		return err
	}
	return nil
}

func topPodsRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if err := podsSelectorRenderHelper(func(namespace string, labelsArr []string) error {
		cmd := kubecli.Cli.WithNamespace(namespace).TopPod(viewStreams(view), nil)
		cmd.SetFlag("selector", strings.Join(labelsArr, ","))
		cmd.Run()
		view.ReRender()
		return nil
	})(gui, view); err != nil {
		return err
	}
	return nil
}

//nolint:funlen
//nolint:funlen
//nolint:funlen
//nolint:funlen
//nolint:funlen
//nolint:funlen
func podsSelectorRenderHelper(cmdFunc func(namespace string, labelsArr []string) error) func(gui *guilib.Gui, view *guilib.View) error {
	return func(gui *guilib.Gui, view *guilib.View) error {
		if activeView == nil {
			return nil
		}
		if activeView.Name == namespaceViewName {
			return namespaceConfigRender(gui, view)
		}
		selected := activeView.SelectedLine
		resource := getViewResourceName(activeView.Name)
		if resource == "" {
			return nil
		}

		jsonPath := resourceLabelSelectorJSONPath(resource)
		if jsonPath == "" {
			return nil
		}

		if notResourceSelected(selected) {
			showPleaseSelected(view, resource)
			return nil
		}

		namespace, resourceName, err := getResourceNamespaceAndName(gui, activeView)
		if err != nil {
			if errors.Is(err, noResourceSelectedErr) {
				showPleaseSelected(view, resource)
				return nil
			}
			return err
		}

		output := newStream()
		cli(namespace).Get(output, resource, resourceName).SetFlag("output", jsonPath).Run()

		var labelJSON = streamToString(output)
		if labelJSON == "" {
			_, err := fmt.Fprint(view, "Pods not found.")
			if err != nil {
				return err
			}

			return nil
		}
		labelsArr := utils.LabelsToStringArr(labelJSON[1 : len(labelJSON)-1])
		if len(labelsArr) == 0 {
			showPleaseSelected(view, resource)
			return nil
		}

		if err := cmdFunc(namespace, labelsArr); err != nil {
			return err
		}
		return nil
	}
}
