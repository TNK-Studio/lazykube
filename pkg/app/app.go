package app

import (
	"github.com/TNK-Studio/lazykube/pkg/config"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/jroimartin/gocui"
)

type App struct {
	ClusterInfo *guilib.View
	Namespace   *guilib.View
	Service     *guilib.View
	Deployment  *guilib.View
	Pod         *guilib.View
	Navigation  *guilib.View
	Detail      *guilib.View
	Option      *guilib.View
	Gui         *guilib.Gui
}

func NewApp() *App {
	app := &App{
		ClusterInfo: ClusterInfo,
		Namespace:   Namespace,
		Service:     Service,
		Deployment:  Deployment,
		Pod:         Pod,
		Navigation:  Navigation,
		Detail:      Detail,
		Option:      Option,
	}

	//Todo: add app config
	conf := config.GuiConfig{
		Highlight:  true,
		SelFgColor: gocui.ColorGreen,
		FgColor:    gocui.ColorWhite,
		Mouse:      true,
		InputEsc:   true,
	}
	app.Gui = guilib.NewGui(
		conf,
		app.ClusterInfo,
		app.Namespace,
		app.Service,
		app.Deployment,
		app.Pod,
		app.Navigation,
		app.Detail,
		app.Option,
	)
	app.Gui.OnRender = app.OnRender
	app.Gui.OnRenderOptions = app.OnRenderOptions
	app.Gui.Actions = actions
	return app
}

func (app *App) Run() {
	app.Gui.Run()
}

func (app *App) Stop() {
	app.Gui.Close()
}

func (app *App) OnRender(gui *guilib.Gui) error {
	if gui.MaxHeight() < 28 {
		for _, viewName := range functionViews {
			if _, err := gui.SetViewOnTop(viewName); err != nil {
				return err
			}
		}
		currentView := gui.CurrentView()
		if currentView != nil {
			if _, err := gui.SetViewOnTop(currentView.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *App) OnRenderOptions(gui *guilib.Gui) error {
	return gui.RenderString(
		app.Option.Name,
		utils.OptionsMapToString(
			map[string]string{
				"← → ↑ ↓":   "navigate",
				"Ctrl+c":    "exit",
				"Esc":       "back",
				"PgUp/PgDn": "scroll",
				"Home/End":  "top/bottom",
				"Tab":       "next panel",
				"F4":        "filter",
				"m":         "menu",
			}),
	)
}
