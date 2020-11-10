package panel

import "github.com/TNK-Studio/lazykube/pkg/gui"

var (
	ClusterInfo *gui.View
)

func init() {
	ClusterInfo = &gui.View{
		Name:      "clusterInfo",
		Title:     "Cluster Info",
		Highlight: true,
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			return leftSideWidth(gui.MaxWidth())
		},
		LowerRightPointYFunc: reactiveHeight,
	}
}
