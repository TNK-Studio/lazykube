package app

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	ClusterInfo = &gui.View{
		Name:      "clusterInfo",
		Title:     "Cluster Info",
		Highlight: true,
		Clickable: true,
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			return leftSideWidth(gui.MaxWidth())
		},
		LowerRightPointYFunc: reactiveHeight,
	}

	Deployment = &gui.View{
		Name:      "deployment",
		Title:     "Deployments",
		FgColor:   gocui.ColorDefault,
		Clickable: true,
		DimensionFunc: gui.BeneathView(
			"service",
			reactiveHeight,
			migrateTopFunc,
		),
	}

	Navigation = &gui.View{
		Name:         "navigation",
		Title:        "Navigation",
		Clickable:    true,
		CanNotReturn: true,
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 0, gui.MaxWidth() - 1, 2
		},
	}

	Detail = &gui.View{
		Name:       "detail",
		Title:      "",
		Clickable:  true,
		FgColor:    gocui.ColorGreen,
		SelFgColor: gocui.ColorBlack | gocui.ColorRed | gocui.AttrBold,
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			return leftSideWidth(gui.MaxWidth()) + 1, 2, gui.MaxWidth() - 1, gui.MaxHeight() - 2
		},
	}

	Namespace = &gui.View{
		Name:      "namespace",
		Title:     "Namespaces",
		Highlight: true,
		Clickable: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			"clusterInfo",
			reactiveHeight,
			migrateTopFunc,
		),
	}

	Option = &gui.View{
		Name: "option",
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			return 0, maxHeight - 2, maxWidth / 3, maxHeight
		},
		NoFrame: true,
		FgColor: gocui.ColorBlue,
	}

	Pod = &gui.View{
		Name:      "pod",
		Title:     "Pods",
		Highlight: true,
		Clickable: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			"deployment",
			reactiveHeight,
			migrateTopFunc,
		),
	}

	Service = &gui.View{
		Name:      "service",
		Title:     "Services",
		Highlight: true,
		Clickable: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			"namespace",
			reactiveHeight,
			migrateTopFunc,
		),
	}
)
