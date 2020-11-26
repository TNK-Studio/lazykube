package app

import (
	"github.com/TNK-Studio/lazykube/pkg/config"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"github.com/fatih/camelcase"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	clusterInfoViewName = "clusterInfo"
	deploymentViewName  = "deployment"
	navigationViewName  = "navigation"
	detailViewName      = "detail"
	namespaceViewName   = "namespace"
	optionViewName      = "option"
	podViewName         = "pod"
	serviceViewName     = "service"
)

var (
	ClusterInfo = &guilib.View{
		Name:      clusterInfoViewName,
		Title:     "Cluster Info",
		Clickable: true,
		ZIndex:    zIndexOfFunctionView(clusterInfoViewName),
		LowerRightPointXFunc: func(gui *guilib.Gui, view *guilib.View) int {
			return leftSideWidth(gui.MaxWidth())
		},
		LowerRightPointYFunc: reactiveHeight,
		OnRender:             renderClusterInfo,
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			changeContext,
			newMoreActions(moreActionsMap[clusterInfoViewName]),
		}),
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			gui.ReRenderViews(navigationViewName, detailViewName)
			return nil
		},
	}

	Deployment = &guilib.View{
		Name:                 deploymentViewName,
		Title:                "Deployments",
		FgColor:              gocui.ColorDefault,
		ZIndex:               zIndexOfFunctionView(deploymentViewName),
		Clickable:            true,
		Highlight:            true,
		SelFgColor:           gocui.ColorGreen,
		OnRender:             resourceListRender,
		OnSelectedLineChange: viewSelectedLineChangeHandler,
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		DimensionFunc: guilib.BeneathView(
			serviceViewName,
			reactiveHeight,
			migrateTopFunc,
		),
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			previousLine,
			nextLine,
			copySelectedLine,
			filterResource,
			editResourceAction,
			newConfirmDialogAction(deploymentViewName, rolloutRestartAction),
			newMoreActions(moreActionsMap[deploymentViewName]),
		}),
	}

	Navigation = &guilib.View{
		Name:         navigationViewName,
		Title:        "Navigation",
		Clickable:    true,
		CanNotReturn: true,
		OnClick:      navigationOnClick,
		FgColor:      gocui.ColorGreen,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 0, gui.MaxWidth() - 1, 2
		},
		OnRender: navigationRender,
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			{
				Name:    navigationArrowLeft,
				Keys:    keyMap[navigationArrowLeft],
				Handler: navigationArrowLeftHandler,
				Mod:     gocui.ModNone,
			},
			{
				Name:    navigationArrowRight,
				Keys:    keyMap[navigationArrowRight],
				Handler: navigationArrowRightHandler,
				Mod:     gocui.ModNone,
			},
			{
				Name: navigationDown,
				Keys: keyMap[navigationDown],
				Handler: func(gui *guilib.Gui, _ *guilib.View) error {
					if err := gui.FocusView(detailViewName, false); err != nil {
						return err
					}
					return nil
				},
				Mod: gocui.ModNone,
			},
		}),
	}

	Detail = &guilib.View{
		Name:       detailViewName,
		Wrap:       true,
		Title:      "",
		Clickable:  true,
		OnRender:   detailRender,
		Highlight:  true,
		SelFgColor: gocui.ColorGreen,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 2, gui.MaxWidth() - 1, gui.MaxHeight() - 2
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			editResourceAction,
			copySelectedLine,
			{
				Keys: keyMap[detailToNavigation],
				Name: detailToNavigation,
				Handler: func(gui *guilib.Gui, view *guilib.View) error {
					return gui.FocusView(navigationViewName, false)
				},
				Mod: gocui.ModNone,
			},
			{
				Name: detailArrowUp,
				Keys: keyMap[detailArrowUp],
				Handler: func(gui *guilib.Gui, view *guilib.View) error {
					_, oy := view.Origin()
					if oy == 0 {
						err := gui.FocusView(navigationViewName, false)
						if err != nil {
							return err
						}
					}
					return scrollUpHandler(gui, view)
				},
				Mod: gocui.ModNone,
			},
			{
				Keys:    keyMap[detailArrowDown],
				Name:    detailArrowDown,
				Handler: scrollDownHandler,
				Mod:     gocui.ModNone,
			},
			changePodLogsContainerAction,
			newMoreActions(moreActionsMap[detailViewName]),
		}),
	}

	Namespace = &guilib.View{
		Name:      namespaceViewName,
		Title:     "Namespaces",
		ZIndex:    zIndexOfFunctionView(deploymentViewName),
		Clickable: true,
		OnRender:  namespaceRender,
		OnSelectedLineChange: func(gui *guilib.Gui, view *guilib.View, selectedLine string) error {
			formatted := formatResourceName(selectedLine, 0)
			if notResourceSelected(formatted) {
				formatted = ""
			}

			if formatted == "" {
				switchNamespace(gui, "")
				return nil
			}
			switchNamespace(gui, formatSelectedNamespace(selectedLine))
			return nil
		},
		Highlight:  true,
		SelFgColor: gocui.ColorGreen,
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		FgColor: gocui.ColorDefault,
		DimensionFunc: guilib.BeneathView(
			clusterInfoViewName,
			reactiveHeight,
			migrateTopFunc,
		),
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			previousLine,
			nextLine,
			copySelectedLine,
			filterResource,
			editResourceAction,
			newMoreActions(moreActionsMap[namespaceViewName]),
		}),
	}

	Option = &guilib.View{
		Name: optionViewName,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			return 0, maxHeight - 2, maxWidth, maxHeight
		},
		AlwaysOnTop: true,
		NoFrame:     true,
		FgColor:     gocui.ColorBlue,
	}

	Pod = &guilib.View{
		Name:                 podViewName,
		Title:                "Pods",
		ZIndex:               zIndexOfFunctionView(deploymentViewName),
		Clickable:            true,
		OnRender:             namespaceResourceListRender("pods"),
		OnSelectedLineChange: viewSelectedLineChangeHandler,
		Highlight:            true,
		SelFgColor:           gocui.ColorGreen,
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		FgColor: gocui.ColorDefault,
		DimensionFunc: guilib.BeneathView(
			deploymentViewName,
			reactiveHeight,
			migrateTopFunc,
		),
		LowerRightPointXFunc: func(gui *guilib.Gui, view *guilib.View) int {
			if resizeableViews[len(resizeableViews)-1] == view.Name {
				return leftSideWidth(gui.MaxWidth())
			}

			_, _, x1, _ := view.DimensionFunc(gui, view)
			return x1
		},
		LowerRightPointYFunc: func(gui *guilib.Gui, view *guilib.View) int {
			_, y0, _, y1 := view.DimensionFunc(gui, view)

			if resizeableViews[len(resizeableViews)-1] == view.Name {
				height := gui.MaxHeight() - 2
				if height < y0+1 {
					return y0 + 1
				}

				return height
			}
			return y1
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			previousLine,
			nextLine,
			copySelectedLine,
			filterResource,
			editResourceAction,
			containerExecCommandAction,
			runPodAction,
			newMoreActions(moreActionsMap[podViewName]),
		}),
	}

	Service = &guilib.View{
		Name:                 serviceViewName,
		Title:                "Services",
		ZIndex:               zIndexOfFunctionView(deploymentViewName),
		Clickable:            true,
		OnRender:             resourceListRender,
		OnSelectedLineChange: viewSelectedLineChangeHandler,
		Highlight:            true,
		SelFgColor:           gocui.ColorGreen,
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		DimensionFunc: guilib.BeneathView(
			namespaceViewName,
			reactiveHeight,
			migrateTopFunc,
		),
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			previousLine,
			nextLine,
			copySelectedLine,
			filterResource,
			editResourceAction,
			newMoreActions(moreActionsMap[namespaceViewName]),
		}),
	}

	viewNameResourceMap = map[string]string{
		namespaceViewName:  namespaceResource,
		serviceViewName:    serviceResource,
		deploymentViewName: deploymentResource,
		podViewName:        podResource,
	}

	restartableResource = []string{"deployments", "statefulsets", "daemonsets"}
)

func getViewResourceName(viewName string) string {
	return viewNameResourceMap[viewName]
}

func newCustomResourcePanel(resource string) *guilib.View {
	viewName := resourceViewName(resource)
	customResourcePanel := &guilib.View{
		Name:                 resourceViewName(resource),
		Title:                resourceViewTitle(resource),
		ZIndex:               zIndexOfFunctionView(viewName),
		Clickable:            true,
		OnRender:             resourceListRender,
		OnSelectedLineChange: viewSelectedLineChangeHandler,
		Highlight:            true,
		SelFgColor:           gocui.ColorGreen,
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		DimensionFunc: guilib.BeneathView(
			functionViews[len(functionViews)-1],
			reactiveHeight,
			migrateTopFunc,
		),
		LowerRightPointXFunc: func(gui *guilib.Gui, view *guilib.View) int {
			if resizeableViews[len(resizeableViews)-1] == view.Name {
				return leftSideWidth(gui.MaxWidth())
			}

			_, _, x1, _ := view.DimensionFunc(gui, view)
			return x1
		},
		LowerRightPointYFunc: func(gui *guilib.Gui, view *guilib.View) int {
			_, y0, _, y1 := view.DimensionFunc(gui, view)

			if resizeableViews[len(resizeableViews)-1] == view.Name {
				height := gui.MaxHeight() - 2
				if height < y0+1 {
					return y0 + 1
				}

				return height
			}
			return y1
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			previousLine,
			nextLine,
			filterResource,
			editResourceAction,
		}),
	}

	customPanelMoreActions := []*moreAction{
		// initialization loop
		//addCustomResourcePanelMoreAction,
		editResourceMoreAction,
		deleteCustomResourcePanelMoreAction,
	}
	if resourceRestartable(resource) {
		customPanelMoreActions = append(
			customPanelMoreActions,
			&moreAction{
				NeedSelectResource: true,
				Action:             *newConfirmDialogAction(customResourcePanel.Name, rolloutRestartAction),
			},
		)
	}

	customResourcePanel.Actions = append(customResourcePanel.Actions, newMoreActions(customPanelMoreActions))
	return customResourcePanel
}

func addCustomResourcePanel(gui *guilib.Gui, resource string) error {
	var customResourcePanel *guilib.View
	customResourcePanel, _ = gui.GetView(resourceViewName(resource))
	if customResourcePanel != nil {
		return nil
	}

	customResourcePanel = newCustomResourcePanel(resource)
	viewNameResourceMap[customResourcePanel.Name] = resource

	// Add to function views and resizeable views.
	functionViews = append(functionViews, customResourcePanel.Name)
	resizeableViews = append(resizeableViews, customResourcePanel.Name)

	// Add custom panel navigation.
	viewNavigationMap[customResourcePanel.Name] = []string{navigationOptConfig, navigationOptDescribe}
	detailRenderMap[navigationPath(customResourcePanel.Name, navigationOptConfig)] = clearBeforeRender(configRender)
	detailRenderMap[navigationPath(customResourcePanel.Name, navigationOptDescribe)] = reRenderInterval(clearBeforeRender(describeRender), reRenderIntervalDuration)

	// Add pods and pods log navigation
	if resourceRestartable(resource) {
		detailRenderMap[navigationPath(customResourcePanel.Name, navigationOptPods)] = reRenderInterval(clearBeforeRender(labelsPodsRender), reRenderIntervalDuration)
		detailRenderMap[navigationPath(customResourcePanel.Name, navigationOptPodsLog)] = reRenderInterval(podsLogsRender, reRenderIntervalDuration)
		detailRenderMap[navigationPath(customResourcePanel.Name, navigationOptTopPods)] = reRenderInterval(clearBeforeRender(topPodsRender), reRenderIntervalDuration)
		viewNavigationMap[customResourcePanel.Name] = append(viewNavigationMap[customResourcePanel.Name], navigationOptPods, navigationOptPodsLog, navigationOptTopPods)
	}

	// Add namespace navigation options.
	viewNavigationMap[namespaceViewName] = append(viewNavigationMap[namespaceViewName], customResourcePanel.Title)
	detailRenderMap[navigationPath(namespaceViewName, customResourcePanel.Title)] = reRenderInterval(
		clearBeforeRender(namespaceResourceListRender(resource)),
		reRenderIntervalDuration,
	)

	if err := resizePanelHeight(gui); err != nil {
		return err
	}
	if err := gui.AddView(customResourcePanel); err != nil {
		return err
	}

	if err := gui.FocusView(customResourcePanel.Name, false); err != nil {
		return err
	}
	config.Conf.UserConfig.AddCustomResourcePanels(resource)
	config.Save()
	return nil
}

func deleteCustomResourcePanel(gui *guilib.Gui, viewName string) error {
	var customResourcePanel *guilib.View
	customResourcePanel, _ = gui.GetView(viewName)
	if customResourcePanel == nil {
		return nil
	}

	for index, eachViewName := range functionViews {
		if eachViewName == viewName {
			functionViews = append(functionViews[:index], functionViews[index+1:]...)
		}
	}

	for index, eachViewName := range resizeableViews {
		if eachViewName == viewName {
			resizeableViews = append(resizeableViews[:index], resizeableViews[index+1:]...)
		}
	}

	for index, option := range viewNavigationMap[namespaceViewName] {
		if option == customResourcePanel.Title {
			viewNavigationMap[namespaceViewName] = append(
				viewNavigationMap[namespaceViewName][:index],
				viewNavigationMap[namespaceViewName][index+1:]...,
			)
		}
	}

	if err := resizePanelHeight(gui); err != nil {
		return err
	}
	if err := gui.DeleteView(customResourcePanel.Name); err != nil {
		return err
	}
	if err := gui.FocusView(functionViews[0], false); err != nil {
		return err
	}
	config.Conf.UserConfig.DeleteCustomResourcePanels(getViewResourceName(customResourcePanel.Name))
	config.Save()
	return nil
}

func resourceViewName(resource string) string {
	gvk := kubecli.Cli.GetResourceGroupVersionKind(resource)
	return strings.ToLower(gvk.Kind)
}

func resourceViewTitle(resource string) string {
	gvk := kubecli.Cli.GetResourceGroupVersionKind(resource)
	return strings.Join(camelcase.Split(gvk.Kind), " ")
}

func zIndexOfFunctionView(viewName string) int {
	i := 0
	for i < len(functionViews) {
		if functionViews[i] == viewName {
			return i
		}
		i++
	}
	return i
}

func resourceRestartable(resource string) bool {
	for _, restartable := range restartableResource {
		if resource == restartable {
			return true
		}
	}
	return false
}
