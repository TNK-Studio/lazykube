package panel

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	Service *gui.View
)

func init() {
	Service = &gui.View{
		Name:      "service",
		Title:     "Services",
		Highlight: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			"namespace",
			reactiveHeight,
			1,
		),
	}
}
