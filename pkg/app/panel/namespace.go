package panel

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	Namespace *gui.View
)

func init() {
	Namespace = &gui.View{
		Name:      "namespace",
		Title:     "Namespaces",
		Highlight: true,
		FgColor:   gocui.ColorDefault,
		DimensionFunc: gui.BeneathView(
			"clusterInfo",
			reactiveHeight,
			1,
		),
	}
}
