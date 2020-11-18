package gui

import (
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
)

var (
	Quit = &Action{
		Name: "Quit",
		Key:  gocui.KeyCtrlC,
		Handler: func(gui *Gui) func(*gocui.Gui, *gocui.View) error {
			return func(*gocui.Gui, *gocui.View) error {
				return gocui.ErrQuit
			}
		},
		Mod: gocui.ModNone,
	}

	ClickView = &Action{
		Name:    "clickView",
		Key:     gocui.MouseLeft,
		Handler: ViewClickHandler,
		Mod:     gocui.ModNone,
	}
)

type Action struct {
	Name            string
	Key             interface{}
	Keys            []interface{}
	ReRenderAllView bool
	Handler         func(gui *Gui) func(*gocui.Gui, *gocui.View) error
	Mod             gocui.Modifier
}

type ActionHandler func(gui *Gui) func(*gocui.Gui, *gocui.View) error

func ViewClickHandler(gui *Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		viewName := v.Name()
		log.Logger.Debugf("ViewClickHandler - view '%s' on click.", viewName)

		currentView := gui.CurrentView()
		canReturn := true

		if currentView == nil || currentView.Name != viewName {
			canReturn = true

			if currentView != nil {
				canReturn = !currentView.CanNotReturn
			}

			if err := gui.FocusView(viewName, canReturn); err != nil {
				return err
			}
		}

		view, err := gui.GetView(viewName)
		if err != nil {
			if err == gocui.ErrUnknownView {
				log.Logger.Warningf("ViewClickHandler - gui.GetView(%s) error %+v", view, err)
				return nil
			}
			return err
		}

		cx, cy := view.Cursor()
		log.Logger.Debugf("ViewClickHandler - cx %d cy %d", cx, cy)
		line, err := view.Line(cy)
		if err != nil {
			log.Logger.Warningf("ViewClickHandler - view.Line(%d) error %s", cy, err)
		} else {
			log.Logger.Debugf("ViewClickHandler - view.Line(%d) line %s", cy, line)
			if view.OnLineClick != nil {
				if err := view.OnLineClick(gui, view, cy, line); err != nil {
					return err
				}
			}
		}

		if view.OnClick != nil {
			return view.OnClick(gui, view)
		}

		return nil
	}
}
