package app

import (
	"github.com/TNK-Studio/lazykube/pkg/app/panel"
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/jroimartin/gocui"
)

type App struct {
	ClusterInfo *gui.View
	Namespace   *gui.View
	Service     *gui.View
	Deployment  *gui.View
	Pod         *gui.View
	Detail      *gui.View
	Option      *gui.View
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
		Option:      panel.Option,
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
		app.Option,
	)
	app.Gui.Render = app.Render
	app.Gui.RenderOptions = app.RenderOptions
	return app
}

func (app *App) Run() {
	app.Gui.Run()
}

func (app *App) Stop() {
	app.Gui.Close()
}

func (app *App) Render(gui *gui.Gui) error {
	if gui.CurrentView() == nil {
		_, err := gui.GetView(app.Namespace.Name)
		if err != nil {
			return nil
		}

		if err:= gui.SwitchFocus("", app.Namespace.Name, false); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) RenderOptions(gui *gui.Gui) error {
	return gui.RenderString(
		app.Option.Name,
		utils.OptionsMapToString(
			map[string]string{
				"PgUp/PgDn": "scroll",
				"← → ↑ ↓":   "navigate",
				"esc/q":     "close",
			}),
	)
}
