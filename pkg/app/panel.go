package app

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
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
	ClusterInfo = &gui.View{
		Name:      clusterInfoViewName,
		Title:     "Cluster Info",
		Clickable: true,
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			return leftSideWidth(gui.MaxWidth())
		},
		LowerRightPointYFunc: reactiveHeight,
		OnRender:             renderClusterInfo,
	}

	Deployment = &gui.View{
		Name:      deploymentViewName,
		Title:     "Deployments",
		FgColor:   gocui.ColorDefault,
		Clickable: true,
		//Highlight:   true,
		//SelFgColor:  gocui.ColorGreen,
		//SelBgColor:  gocui.ColorDefault,
		OnRender:    deploymentRender,
		OnLineClick: viewLineClickHandler,
		OnFocus: func(gui *gui.Gui, view *gui.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		DimensionFunc: gui.BeneathView(
			serviceViewName,
			reactiveHeight,
			migrateTopFunc,
		),
	}

	Navigation = &gui.View{
		Name:         navigationViewName,
		Title:        "Navigation",
		Clickable:    true,
		CanNotReturn: true,
		OnClick:      navigationOnClick,
		FgColor:      gocui.ColorGreen,
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 0, gui.MaxWidth() - 1, 2
		},
		OnRender: navigationRender,
	}

	Detail = &gui.View{
		Name:      detailViewName,
		Wrap:      true,
		Title:     "",
		Clickable: true,
		OnRender:  detailRender,
		OnFocusLost: func(gui *gui.Gui, view *gui.View) error {
			if err := view.SetCursor(0, 0); err != nil {
				return err
			}

			return nil
		},
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 2, gui.MaxWidth() - 1, gui.MaxHeight() - 2
		},
	}

	Namespace = &gui.View{
		Name:      namespaceViewName,
		Title:     "Namespaces",
		Clickable: true,
		//Highlight:   true,
		//SelFgColor:  gocui.ColorGreen,
		//SelBgColor:  gocui.ColorDefault,
		OnRender:    namespaceRender,
		OnLineClick: viewLineClickHandler,
		OnFocus: func(gui *gui.Gui, view *gui.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		FgColor: gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			clusterInfoViewName,
			reactiveHeight,
			migrateTopFunc,
		),
	}

	Option = &gui.View{
		Name: optionViewName,
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			return 0, maxHeight - 2, maxWidth, maxHeight
		},
		NoFrame: true,
		FgColor: gocui.ColorBlue,
	}

	Pod = &gui.View{
		Name:      podViewName,
		Title:     "Pods",
		Clickable: true,
		//Highlight:   true,
		//SelFgColor:  gocui.ColorGreen,
		//SelBgColor:  gocui.ColorDefault,
		OnRender:    podRender,
		OnLineClick: viewLineClickHandler,
		OnFocus: func(gui *gui.Gui, view *gui.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		FgColor: gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			deploymentViewName,
			reactiveHeight,
			migrateTopFunc,
		),
	}

	Service = &gui.View{
		Name:      serviceViewName,
		Title:     "Services",
		Clickable: true,
		//Highlight:   true,
		//SelFgColor:  gocui.ColorGreen,
		//SelBgColor:  gocui.ColorDefault,
		OnRender:    serviceRender,
		OnLineClick: viewLineClickHandler,
		OnFocus: func(gui *gui.Gui, view *gui.View) error {
			if err := onFocusClearSelected(gui, view); err != nil {
				return err
			}
			return nil
		},
		DimensionFunc: gui.BeneathView(
			namespaceViewName,
			reactiveHeight,
			migrateTopFunc,
		),
	}
)
