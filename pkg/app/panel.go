package app

import (
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
		Name:      detailViewName,
		Wrap:      true,
		Title:     "",
		Clickable: true,
		OnRender:  detailRender,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 2, gui.MaxWidth() - 1, gui.MaxHeight() - 2
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			{
				Name: detailArrowUp,
				Keys: keyMap[detailArrowUp],
				Handler: func(gui *guilib.Gui, _ *guilib.View) error {
					err := gui.FocusView(navigationViewName, false)
					if err != nil {
						return err
					}
					return nil
				},
				Mod: gocui.ModNone,
			},
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
		OnRender:             podRender,
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
			if resizeableViews[len(resizeableViews)-1] == view.Name {
				return gui.MaxHeight() - 2
			}
			_, _, _, y1 := view.DimensionFunc(gui, view)
			return y1
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			previousLine,
			nextLine,
			filterResource,
			editResourceAction,
			newMoreActions(moreActionsMap[namespaceViewName]),
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
)

func getViewResourceName(viewName string) string {
	return viewNameResourceMap[viewName]
}

func newCustomResourcePanel(resource string) *guilib.View {
	viewName := resourceViewName(resource)
	return &guilib.View{
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
			if resizeableViews[len(resizeableViews)-1] == view.Name {
				return gui.MaxHeight() - 2
			}
			_, _, _, y1 := view.DimensionFunc(gui, view)
			return y1
		},
	}
}

func addCustomResourcePanel(gui *guilib.Gui, resource string) error {
	var customResourcePanel *guilib.View
	customResourcePanel, _ = gui.GetView(resourceViewName(resource))
	if customResourcePanel != nil {
		return nil
	}

	customResourcePanel = newCustomResourcePanel(resource)
	viewNameResourceMap[customResourcePanel.Name] = resource
	functionViews = append(functionViews, customResourcePanel.Name)
	resizeableViews = append(resizeableViews, customResourcePanel.Name)
	viewNavigationMap[customResourcePanel.Name] = []string{"Config", "Describe"}
	detailRenderMap[navigationPath(customResourcePanel.Name, "Config")] = clearBeforeRender(configRender)
	detailRenderMap[navigationPath(customResourcePanel.Name, "Describe")] = clearBeforeRender(describeRender)
	if err := gui.AddView(customResourcePanel); err != nil {
		return err
	}
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
