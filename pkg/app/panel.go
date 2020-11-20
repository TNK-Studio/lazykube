package app

import (
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
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
		LowerRightPointXFunc: func(gui *guilib.Gui, view *guilib.View) int {
			return leftSideWidth(gui.MaxWidth())
		},
		LowerRightPointYFunc: reactiveHeight,
		OnRender:             renderClusterInfo,
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
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
		ZIndex:               3,
		Clickable:            true,
		Highlight:            true,
		SelFgColor:           gocui.ColorGreen,
		OnRender:             deploymentRender,
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
			filterAction,
			editResourceAction,
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
				Name: "navigationArrowLeft",
				Keys: []interface{}{
					gocui.KeyArrowLeft,
					'k',
				},
				Handler: navigationArrowLeftHandler,
				Mod:     gocui.ModNone,
			},
			{
				Name: "navigationArrowRight",
				Keys: []interface{}{
					gocui.KeyArrowRight,
					'l',
				},
				Handler: navigationArrowRightHandler,
				Mod:     gocui.ModNone,
			},
			{
				Name: "navigationDown",
				Key:  gocui.KeyArrowDown,
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
				Name: "detailArrowUp",
				Key:  gocui.KeyArrowUp,
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
		ZIndex:    1,
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
			filterAction,
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
		ZIndex:               4,
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
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toNavigation,
			nextCyclicView,
			previousLine,
			nextLine,
			filterAction,
			editResourceAction,
			newMoreActions(moreActionsMap[namespaceViewName]),
		}),
	}

	Service = &guilib.View{
		Name:                 serviceViewName,
		Title:                "Services",
		ZIndex:               2,
		Clickable:            true,
		OnRender:             serviceRender,
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
			filterAction,
			editResourceAction,
			newMoreActions(moreActionsMap[namespaceViewName]),
		}),
	}
)

func getViewResourceName(viewName string) string {
	var resource string
	switch viewName {
	case namespaceViewName:
		resource = namespaceResource
	case serviceViewName:
		resource = serviceResource
	case deploymentViewName:
		resource = deploymentResource
	case podViewName:
		resource = podResource
	}

	// Todo
	//if resource == "" {
	//
	//}

	return resource
}
