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
	optSeparator       = "   "
	navigationPathJoin = " + "
	logsTail           = "500"

	namespaceResource  = "namespace"
	serviceResource    = "service"
	deploymentResource = "deployment"
	podResource        = "pod"
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

	detailRenderMap = map[string]guilib.ViewHandler{
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

func renderAfterClear(render guilib.ViewHandler) guilib.ViewHandler {
	return func(gui *guilib.Gui, view *guilib.View) error {
		view.Clear()
		return render(gui, view)
	}
}

func navigationPath(args ...string) string {
	return strings.Join(args, navigationPathJoin)
}

func switchNavigation(index int) string {
	err := Detail.SetOrigin(0, 0)
	if err != nil {
		log.Logger.Warningf("switchNavigation - Detail.SetOrigin(0, 0) error %s", err)
	}

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
	str := strings.Join(colorfulOptions, optSeparator)

	_, err := fmt.Fprint(view, str)
	if err != nil {
		return err
	}

	return nil
}

func navigationOnClick(_ *guilib.Gui, view *guilib.View) error {
	cx, cy := view.Cursor()
	log.Logger.Debugf("navigationOnClick - cx %d cy %d", cx, cy)

	options := viewNavigationMap[activeView.Name]
	optionIndex, selected := utils.ClickOption(options, optSeparator, cx, 0)
	if optionIndex < 0 {
		return nil
	}
	log.Logger.Debugf("navigationOnClick - cx %d selected '%s'", cx, selected)
	selected = switchNavigation(optionIndex)
	view.ReRender()
	Detail.ReRender()
	return nil
}

func renderClusterInfo(_ *guilib.Gui, view *guilib.View) error {
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

func clusterNodesRender(_ *guilib.Gui, view *guilib.View) error {
	kubecli.Cli.Get(viewStreams(view), "nodes").Run()
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

func serviceRender(_ *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if kubecli.Cli.Namespace() == "" {
		kubecli.Cli.Get(viewStreams(view), "services").SetFlag("all-namespaces", "true").Run()
		return nil
	}
	kubecli.Cli.Get(viewStreams(view), "services").Run()
	return nil
}

func deploymentRender(_ *guilib.Gui, view *guilib.View) error {
	view.Clear()
	if kubecli.Cli.Namespace() == "" {
		kubecli.Cli.Get(viewStreams(view), "deployments").SetFlag("all-namespaces", "true").Run()
		return nil
	}
	kubecli.Cli.Get(viewStreams(view), "deployments").Run()
	return nil
}

func podRender(_ *guilib.Gui, view *guilib.View) error {
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

func newStdStream() genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
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

func showPleaseSelected(view io.Writer, name string) {
	_, err := fmt.Fprintf(view, "Please select a %s. ", name)
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
	if activeView == Namespace {
		return namespaceConfigRender(gui, view)
	}

	resource := getViewResourceName(activeView.Name)

	if resource == "" {
		return nil
	}

	namespace, resourceName, err := getResourceNamespaceAndName(gui, activeView)
	if err != nil {
		return err
	}

	if notResourceSelected(resourceName) {
		showPleaseSelected(view, resource)
		return nil
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
	podView, err := gui.GetView(podViewName)
	if err != nil {
		return err
	}

	resource := "pod"

	namespace, resourceName, err := getResourceNamespaceAndName(gui, podView)
	if err != nil {
		return err
	}

	if notResourceSelected(resourceName) {
		showPleaseSelected(view, resource)
		return nil
	}

	cli(namespace).
		Logs(viewStreams(view), resourceName).
		SetFlag("all-containers", "true").
		SetFlag("tail", logsTail).
		SetFlag("prefix", "true").
		Run()

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
		if activeView == Namespace {
			return namespaceConfigRender(gui, view)
		}
		selected := activeView.SelectedLine
		var resource, jsonPath string
		switch activeView.Name {
		case serviceViewName:
			resource = "service"
			jsonPath = "jsonpath='{.spec.selector}'"
		case deploymentViewName:
			resource = "deployment"
			jsonPath = "jsonpath='{.spec.selector.matchLabels}'"
		}

		if resource == "" {
			return nil
		}

		if notResourceSelected(selected) {
			showPleaseSelected(view, resource)
			return nil
		}

		namespace, resourceName, err := getResourceNamespaceAndName(gui, activeView)
		if err != nil {
			return err
		}

		if notResourceSelected(resourceName) {
			showPleaseSelected(view, resource)
			return nil
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
