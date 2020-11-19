package gui

import (
	"errors"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
)

var (
	// Quit Quit
	Quit = &Action{
		Name: "Quit",
		Key:  gocui.KeyCtrlC,
		Handler: func(*Gui, *View) error {
			return gocui.ErrQuit
		},
		Mod: gocui.ModNone,
	}

	// ClickView ClickView
	ClickView = &Action{
		Name:    "clickView",
		Key:     gocui.MouseLeft,
		Handler: ViewClickHandler,
		Mod:     gocui.ModNone,
	}
)

// ActionInterface ActionInterface
type ActionInterface interface {
	ActionName() string
	HandlerFunc(*Gui, *View) error
	Modifier() gocui.Modifier
	BindKey() interface{}
	BindKeys() []interface{}
	ReRenderAll() bool
}

// Action Action
type Action struct {
	Keys            []interface{}
	Name            string
	Key             interface{}
	Handler         func(*Gui, *View) error
	ReRenderAllView bool
	Mod             gocui.Modifier
}

func (a Action) HandlerFunc(gui *Gui, view *View) error {
	return a.Handler(gui, view)
}

func (a Action) ActionName() string {
	return a.Name
}

func (a Action) Modifier() gocui.Modifier {
	return a.Mod
}

func (a Action) BindKey() interface{} {
	return a.Key
}

func (a Action) BindKeys() []interface{} {
	return a.Keys
}

func (a Action) ReRenderAll() bool {
	return a.ReRenderAllView
}

func ToActionInterfaceArr(actions []*Action) []ActionInterface {
	arr := make([]ActionInterface, 0)
	for _, act := range actions {
		arr = append(arr, act)
	}
	return arr
}

type ActionHandler func(gui *Gui) func(*gocui.Gui, *gocui.View) error

func ViewClickHandler(gui *Gui, view *View) error {
	viewName := view.Name
	log.Logger.Debugf("ViewClickHandler - view '%s' on click.", viewName)

	currentView := gui.CurrentView()

	var canReturn bool
	canReturn = true
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
		if errors.Is(err, gocui.ErrUnknownView) {
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

func actionHandlerWrapper(gui *Gui, handler func(gui *Gui, view *View) error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		view := gui.getView(v.Name())
		return handler(gui, view)
	}
}
