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
	}

	Deployment = &guilib.View{
		Name:        deploymentViewName,
		Title:       "Deployments",
		FgColor:     gocui.ColorDefault,
		Clickable:   true,
		OnRender:    deploymentRender,
		OnLineClick: viewLineClickHandler,
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
	}

	Namespace = &guilib.View{
		Name:        namespaceViewName,
		Title:       "Namespaces",
		Clickable:   true,
		OnRender:    namespaceRender,
		OnLineClick: viewLineClickHandler,
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
	}

	Option = &guilib.View{
		Name: optionViewName,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			return 0, maxHeight - 2, maxWidth, maxHeight
		},
		NoFrame: true,
		FgColor: gocui.ColorBlue,
	}

	Pod = &guilib.View{
		Name:        podViewName,
		Title:       "Pods",
		Clickable:   true,
		OnRender:    podRender,
		OnLineClick: viewLineClickHandler,
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
	}

	Service = &guilib.View{
		Name:        serviceViewName,
		Title:       "Services",
		Clickable:   true,
		OnRender:    serviceRender,
		OnLineClick: viewLineClickHandler,
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
	}
)

func setViewSelectedLine(gui *guilib.Gui, view *guilib.View, selectedLine string) error {
	formatted := formatSelectedName(selectedLine, 0)
	if formatted == "NAME" || formatted == "NAMESPACE" {
		formatted = ""
	}

	if formatted == "" {
		if err := view.State.Set(selectedViewLine, nil); err != nil {
			return err
		}
		return nil
	}

	if view.Name == namespaceViewName {
		switchNamespace(gui, formatSelectedNamespace(selectedLine))
	}

	if err := view.State.Set(selectedViewLine, selectedLine); err != nil {
		return err
	}
	return nil
}
