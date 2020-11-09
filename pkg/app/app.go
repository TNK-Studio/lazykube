package app

import (
	"github.com/TNK-Studio/lazykube/pkg/app/panel"
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

type App struct {
	ClusterInfo *gui.View
	Namespace   *gui.View
	Service     *gui.View
	Deployment  *gui.View
	Pod         *gui.View
	Detail      *gui.View
	Gui         *gui.Gui
}

func NewApp() *App {
	app := &App{
		ClusterInfo: panel.ClusterInfo,
		Namespace:   panel.Namespace,
		Service:     panel.Service,
		Deployment:  panel.Deployment,
		Pod:         panel.Pod,
		Detail:      panel.Detail,
	}

	//Todo: add app config
	conf := config.GuiConfig{
		Highlight:  true,
		SelFgColor: gocui.ColorGreen,
	}
	app.Gui = gui.NewGui(
		conf,
		app.ClusterInfo,
		app.Namespace,
		app.ClusterInfo,
		app.Namespace,
		app.Service,
		app.Deployment,
		app.Pod,
		app.Detail,
	)
	return app
}

func (app *App) Run() {
	app.Gui.Run()
}

func (app *App) Render(gui *gui.Gui) error {
	return nil
}
