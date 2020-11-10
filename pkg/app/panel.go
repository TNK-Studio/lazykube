package app

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	ClusterInfo *gui.View
	Deployment *gui.View
	Detail *gui.View
	Namespace *gui.View
	Option *gui.View
	Pod *gui.View
	Service *gui.View
)

func init() {
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
		Name:    "deployment",
		Title:   "Deployments",
		FgColor: gocui.ColorDefault,
		Clickable: true,
		DimensionFunc: gui.BeneathView(
			"service",
			reactiveHeight,
			1,
		),
	}

	Detail = &gui.View{
		Name:  "detail",
		Title: "",
		Clickable: true,
		UpperLeftPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			return leftSideWidth(gui.MaxWidth()) + 1
		},
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			return gui.MaxWidth() - 1
		},
		LowerRightPointYFunc: func(gui *gui.Gui, view *gui.View) int {
			return gui.MaxHeight() - 2
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
			1,
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
			1,
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
			1,
		),
	}
}
