package panel

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	Pod *gui.View
)

func init() {
	Pod = &gui.View{
		Name:      "pod",
		Title:     "Pods",
		Highlight: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			"deployment",
			reactiveHeight,
			1,
		),
	}
}
