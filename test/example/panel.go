package main

import (
	"github.com/TNK-Studio/lazykube/pkg/app/panel"
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

func main() {
	conf := config.GuiConfig{
		Highlight:  true,
		SelFgColor: gocui.ColorGreen,
	}

	g := gui.NewGui(
		conf,
		panel.Detail,
		panel.ClusterInfo,
		panel.Namespace,
		panel.Service,
		panel.Deployment,
		panel.Pod,
	)
	defer g.Close()
	g.Run()
}
