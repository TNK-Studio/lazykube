package main

import (
	"github.com/TNK-Studio/lazykube/pkg/app"
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

func main() {
	conf := config.GuiConfig{
		Highlight:  true,
		SelFgColor: gocui.ColorGreen,
		Mouse: true,
	}

	g := gui.NewGui(
		conf,
		app.Detail,
		app.ClusterInfo,
		app.Namespace,
		app.Service,
		app.Deployment,
		app.Pod,
	)
	defer g.Close()

	//g.BindAction(app.Service.Name, gui.ClickView)
	g.SetKeybinding(
		app.Service.Name,
		gocui.MouseLeft,
		gocui.ModNone,
		func(gui *gocui.Gui, v *gocui.View) error {
			if _, err := gui.SetCurrentView(app.Service.Name); err != nil {
				return err
			}
			if _, err := gui.SetViewOnTop(app.Service.Name); err != nil {
				return err
			}
			return nil
		},
	)
	g.Run()
}
