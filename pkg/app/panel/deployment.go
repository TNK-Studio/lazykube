package panel

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	Deployment *gui.View
)

func init() {
	Deployment = &gui.View{
		Name:    "deployment",
		Title:   "Deployments",
		FgColor: gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			"service",
			reactiveHeight,
			1,
		),
	}
}
