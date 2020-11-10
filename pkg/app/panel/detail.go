package panel

import "github.com/TNK-Studio/lazykube/pkg/gui"

var (
	Detail *gui.View
)

func init() {
	Detail = &gui.View{
		Name:  "detail",
		Title: "",
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
}
