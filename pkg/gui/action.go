package gui

import "github.com/jroimartin/gocui"

var (
	Quit = &Action{
		Name: "Quit",
		Key: gocui.KeyCtrlC,
		Handler: func(gui *Gui) func(*gocui.Gui, *gocui.View) error {
			return func(*gocui.Gui, *gocui.View) error {
				return gocui.ErrQuit
			}
		},
		Mod: gocui.ModNone,
	}
)

type Action struct {
	Name string
	Key     interface{}
	Handler func(gui *Gui) func(*gocui.Gui, *gocui.View) error
	Mod     gocui.Modifier
}
