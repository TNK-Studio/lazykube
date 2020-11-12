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
		Render:               renderClusterInfo,
	}

	Deployment = &gui.View{
		Name:      deploymentViewName,
		Title:     "Deployments",
		FgColor:   gocui.ColorDefault,
		Clickable: true,
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
		SelFgColor:   gocui.ColorBlack | gocui.ColorRed | gocui.AttrBold,
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 0, gui.MaxWidth() - 1, 2
		},
		Render: navigationRender,
	}

	Detail = &gui.View{
		Name:      detailViewName,
		Title:     "",
		Clickable: true,
		Render:    detailRender,
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 2, gui.MaxWidth() - 1, gui.MaxHeight() - 2
		},
	}

	Namespace = &gui.View{
		Name:      namespaceViewName,
		Title:     "Namespaces",
		Highlight: true,
		Clickable: true,
		Render:    namespaceRender,
		FgColor:   gocui.ColorDefault,
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
		Highlight: true,
		Clickable: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			deploymentViewName,
			reactiveHeight,
			migrateTopFunc,
		),
	}

	Service = &gui.View{
		Name:      serviceViewName,
		Title:     "Services",
		Highlight: true,
		Clickable: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			namespaceViewName,
			reactiveHeight,
			migrateTopFunc,
		),
	}
)
