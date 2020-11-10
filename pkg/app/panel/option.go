package panel

import (
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	Option *gui.View
)

func init() {
	Option = &gui.View{
		Name: "option",
		DimensionFunc: func(gui *gui.Gui, view *gui.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			return 0, maxHeight - 2, maxWidth / 3, maxHeight
		},
		NoFrame: true,
		FgColor: gocui.ColorBlue,
	}
}