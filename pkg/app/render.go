package app

import (
	"bytes"
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
)

const (
	OptSeparator       = "   "
	navigationPathJoin = " + "
	logsTail           = "500"
)

var (
	// Todo: use state to control.
	activeView *guilib.View

	navigationIndex     int
	activeNavigationOpt string

	functionViews     = []string{clusterInfoViewName, namespaceViewName, serviceViewName, deploymentViewName, podViewName}
	viewNavigationMap = map[string][]string{
		clusterInfoViewName: {"Nodes", "Top Nodes"},
		namespaceViewName:   {"Config", "Deployments", "Pods"},
		serviceViewName:     {"Config", "Pods", "Pods Log", "Top Pods"},
		deploymentViewName:  {"Config", "Pods", "Pods Log", "Describe", "Top Pods"},
		podViewName:         {"Log", "Config", "Top", "Describe"},
	}

	detailRenderMap = map[string]func(gui *guilib.Gui, view *guilib.View) error{
		navigationPath(clusterInfoViewName, "Nodes"):     renderAfterClear(clusterNodesRender),
		navigationPath(clusterInfoViewName, "Top Nodes"): renderAfterClear(topNodesRender),
		navigationPath(namespaceViewName, "Deployments"): renderAfterClear(deploymentRender),
		navigationPath(namespaceViewName, "Pods"):        renderAfterClear(podRender),
		navigationPath(namespaceViewName, "Config"):      renderAfterClear(configRender),
		navigationPath(serviceViewName, "Config"):        renderAfterClear(configRender),
		navigationPath(serviceViewName, "Pods"):          renderAfterClear(labelsPodsRender),
		navigationPath(serviceViewName, "Pods Log"):      podsLogsRender,
		navigationPath(serviceViewName, "Top Pods"):      renderAfterClear(topPodsRender),
		navigationPath(deploymentViewName, "Config"):     renderAfterClear(configRender),
		navigationPath(deploymentViewName, "Pods"):       renderAfterClear(labelsPodsRender),
		navigationPath(deploymentViewName, "Describe"):   renderAfterClear(describeRender),
		navigationPath(deploymentViewName, "Pods Log"):   podsLogsRender,
		navigationPath(deploymentViewName, "Top Pods"):   renderAfterClear(topPodsRender),
		navigationPath(podViewName, "Config"):            renderAfterClear(configRender),
		navigationPath(podViewName, "Log"):               podLogsRender,
		navigationPath(podViewName, "Describe"):          renderAfterClear(describeRender),
		navigationPath(podViewName, "Top"):               podMetricsPlotRender,
	}
)

func notResourceSelected(selectedName string) bool {
	if selectedName == "" || selectedName == "NAME" || selectedName == "NAMESPACE" || selectedName == "No" {
		return true
	}
	return false
}

func renderAfterClear(render func(gui *guilib.Gui, view *guilib.View) error) func(gui *guilib.Gui, view *guilib.View) error {
	return func(gui *guilib.Gui, view *guilib.View) error {
		view.Clear()
		return render(gui, view)
	}
}

func navigationPath(args ...string) string {
	return strings.Join(args, navigationPathJoin)
}

func switchNavigation(index int) string {
	Detail.SetOrigin(0, 0)
	Detail.Clear()
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
			if err := gui.FocusView(clusterInfoViewName, false); err != nil {
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

func navigationOnClick(gui *guilib.Gui, view *guilib.View) error {
	cx, cy := view.Cursor()
	log.Logger.Debugf("navigationOnClick - cx %d cy %d", cx, cy)

	options := viewNavigationMap[activeView.Name]
	sep := len(OptSeparator)
	sections := make([]int, 0)
	preFix := 0

	var selected string
	for i, opt := range options {
		left := preFix + i*sep

		words := len([]rune(opt))

		right := left + words - 1
		preFix += words

		sections = append(sections, left, right)
	}

	log.Logger.Debugf("navigationOnClick - sections %+v", sections)

	for i := 0; i < len(sections); i += 2 {
		left := sections[i]
		right := sections[i+1]
		if cx >= left && cx <= right {
			optionIndex := i / 2
			log.Logger.Debugf("navigationOnClick - cx %d in selection(%d)[%d, %d]", cx, optionIndex, left, right)
			selected = switchNavigation(optionIndex)
			view.ReRender()
			Detail.ReRender()
			break
		}
	}

	log.Logger.Debugf("navigationOnClick - cx %d selected '%s'", cx, selected)

	return nil
}

func renderClusterInfo(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	currentContext, err := kubecli.Cli.CurrentContext()
	if err != nil {
		return nil
	}

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

func clusterNodesRender(gui *guilib.Gui, view *guilib.View) error {
	kubecli.Cli.Get(viewStreams(view), "nodes").Run()
	return nil
}

func topNodesRender(gui *guilib.Gui, view *guilib.View) error {
	kubecli.Cli.TopNode(viewStreams(view), nil, "").Run()
	view.ReRender()
	return nil
}

func namespaceRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	kubecli.Cli.Get(viewStreams(view), "namespaces").Run()
	return nil
}

func serviceRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if kubecli.Cli.Namespace() == "" {
		kubecli.Cli.Get(viewStreams(view), "services").SetFlag("all-namespaces", "true").Run()
		return nil
	}
	kubecli.Cli.Get(viewStreams(view), "services").Run()
	return nil
}

func deploymentRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if kubecli.Cli.Namespace() == "" {
		kubecli.Cli.Get(viewStreams(view), "deployments").SetFlag("all-namespaces", "true").Run()
		return nil
	}
	kubecli.Cli.Get(viewStreams(view), "deployments").Run()
	return nil
}

func podRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if kubecli.Cli.Namespace() == "" {
		kubecli.Cli.Get(viewStreams(view), "pods").SetFlag("all-namespaces", "true").SetFlag("output", "wide").Run()
		return nil
	}
	kubecli.Cli.Get(viewStreams(view), "pods").SetFlag("output", "wide").Run()
	return nil
}

func newStream() genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     &bytes.Buffer{},
		Out:    &bytes.Buffer{},
		ErrOut: &bytes.Buffer{},
	}
}

func streamCopyTo(streams genericclioptions.IOStreams, writer io.Writer) {
	if _, err := io.Copy(writer, (streams.Out).(io.Reader)); err != nil {
		log.Logger.Warningf("streamCopyTo - streams.Out copy error %s", err)
	}
	if _, err := io.Copy(writer, (streams.ErrOut).(io.Reader)); err != nil {
		log.Logger.Warningf("streamCopyTo - streams.ErrOut copy error %s", err)
	}
}

func streamToString(streams genericclioptions.IOStreams) string {
	buf := new(strings.Builder)
	streamCopyTo(streams, buf)
	// check errors
	return buf.String()
}

func showPleaseSelected(view *guilib.View, name string) {
	fmt.Fprintf(view, "Please select a %s. ", name)
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

	namespace := formatSelectedNamespace(namespaceView.SelectedLine)
	selected := activeView.SelectedLine
	resource := ""
	switch activeView.Name {
	case serviceViewName:
		resource = "service"
		break
	case deploymentViewName:
		resource = "deployment"
		break
	case podViewName:
		resource = "pod"
		break
	}

	if resource == "" {
		return nil
	}

	if selected == "" {
		showPleaseSelected(view, resource)
		return nil
	}

	if !notResourceSelected(namespace) {
		resourceName := formatResourceName(selected, 0)
		if notResourceSelected(resourceName) {
			showPleaseSelected(view, resource)
			return nil
		}
		kubecli.Cli.Get(viewStreams(view), resource, resourceName).SetFlag("output", "yaml").Run()
		return nil
	}

	namespace = formatResourceName(selected, 0)
	resourceName := formatResourceName(selected, 1)
	if notResourceSelected(resourceName) {
		showPleaseSelected(view, resource)
		return nil
	}

	kubecli.Cli.WithNamespace(namespace).Get(viewStreams(view), resource, resourceName).SetFlag("output", "yaml").Run()
	return nil
}

func describeRender(gui *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if activeView == nil {
		return nil
	}
	if activeView == Namespace {
		return namespaceConfigRender(gui, view)
	}
	namespaceView, err := gui.GetView(namespaceViewName)
	if err != nil {
		return nil
	}

	namespace := formatSelectedNamespace(namespaceView.SelectedLine)
	selected := activeView.SelectedLine
	resource := ""
	switch activeView.Name {
	case deploymentViewName:
		resource = "deployment"
		break
	case podViewName:
		resource = "pod"
		break
	}

	if resource == "" {
		return nil
	}

	if notResourceSelected(selected) {
		showPleaseSelected(view, resource)
		return nil
	}

	if !notResourceSelected(namespace) {
		resourceName := formatResourceName(selected, 0)
		if notResourceSelected(resourceName) {
			showPleaseSelected(view, resource)
			return nil
		}
		kubecli.Cli.Describe(viewStreams(view), resource, resourceName).Run()
		return nil
	}

	namespace = formatResourceName(selected, 0)
	resourceName := formatResourceName(selected, 1)
	if notResourceSelected(resourceName) {
		showPleaseSelected(view, resource)
		return nil
	}

	kubecli.Cli.WithNamespace(namespace).Describe(viewStreams(view), resource, resourceName).Run()
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
	namespaceView, err := gui.GetView(namespaceViewName)
	if err != nil {
		return err
	}

	namespace := formatSelectedNamespace(namespaceView.SelectedLine)

	podView, err := gui.GetView(podViewName)
	if err != nil {
		return err
	}

	selected := podView.SelectedLine
	resource := "pod"
	if notResourceSelected(selected) {
		showPleaseSelected(view, resource)
		return nil
	}

	if !notResourceSelected(namespace) {
		resourceName := formatResourceName(selected, 0)
		if notResourceSelected(resourceName) {
			showPleaseSelected(view, resource)
			return nil
		}
		streams := newStream()
		kubecli.Cli.Logs(streams, resourceName).SetFlag("all-containers", "true").SetFlag("tail", logsTail).SetFlag("prefix", "true").Run()
		view.Clear()
		streamCopyTo(streams, view)
		view.ReRender()
		return nil
	}

	namespace = formatResourceName(selected, 0)
	resourceName := formatResourceName(selected, 1)
	if notResourceSelected(resourceName) {
		showPleaseSelected(view, resource)
		return nil
	}

	streams := newStream()
	kubecli.Cli.WithNamespace(namespace).Logs(streams, resourceName).SetFlag("all-containers", "true").SetFlag("tail", logsTail).SetFlag("prefix", "true").Run()
	streamCopyTo(streams, view)
	view.ReRender()
	return nil
}

func podsLogsRender(gui *guilib.Gui, view *guilib.View) error {
	if err := podsSelectorRenderHelper(func(namespace string, labelsArr []string) error {
		streams := newStream()
		cmd := kubecli.Cli.WithNamespace(namespace).Logs(streams)
		cmd.SetFlag("selector", strings.Join(labelsArr, ","))
		cmd.SetFlag("all-containers", "true").SetFlag("tail", logsTail).SetFlag("prefix", "true").Run()
		view.Clear()
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

func podsSelectorRenderHelper(cmdFunc func(namespace string, labelsArr []string) error) func(gui *guilib.Gui, view *guilib.View) error {
	return func(gui *guilib.Gui, view *guilib.View) error {
		if activeView == nil {
			return nil
		}
		if activeView == Namespace {
			return namespaceConfigRender(gui, view)
		}
		namespaceView, err := gui.GetView(namespaceViewName)
		if err != nil {
			return nil
		}

		namespace := formatSelectedNamespace(namespaceView.SelectedLine)
		selected := activeView.SelectedLine
		var resource string
		var jsonPath string
		switch activeView.Name {
		case serviceViewName:
			resource = "service"
			jsonPath = "jsonpath='{.spec.selector}'"
			break
		case deploymentViewName:
			resource = "deployment"
			jsonPath = "jsonpath='{.spec.selector.matchLabels}'"
			break
		}

		if resource == "" {
			return nil
		}

		if notResourceSelected(selected) {
			showPleaseSelected(view, resource)
			return nil
		}

		output := newStream()
		if !notResourceSelected(namespace) {
			resourceName := formatResourceName(selected, 0)
			if notResourceSelected(resourceName) {
				showPleaseSelected(view, resource)
				return nil
			}
			kubecli.Cli.Get(output, resource, resourceName).SetFlag("output", jsonPath).Run()
		} else {
			namespace = formatResourceName(selected, 0)
			resourceName := formatResourceName(selected, 1)
			if notResourceSelected(resourceName) {
				showPleaseSelected(view, resource)
				return nil
			}
			kubecli.Cli.WithNamespace(namespace).Get(output, resource, resourceName).SetFlag("output", jsonPath).Run()
		}

		labelJson := streamToString(output)
		if labelJson == "" {
			fmt.Fprint(view, "Pods not found.")
			return nil
		}
		labelsArr := utils.LabelsToStringArr(labelJson[1 : len(labelJson)-1])
		if len(labelsArr) == 0 {
			showPleaseSelected(view, resource)
			return nil
		}

		if namespace == "" {
			namespace = kubecli.Cli.Namespace()
		}

		if err := cmdFunc(namespace, labelsArr); err != nil {
			return err
		}
		return nil
	}
}
