package gui

import (
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/jroimartin/gocui"
	"log"
	"time"
)

type Gui struct {
	State  State
	Render func(gui *Gui) error

	g      *gocui.Gui
	views  []*View
	config config.GuiConfig
}

func NewGui(config config.GuiConfig, views ...*View) *Gui {

	gui := &Gui{
		State: &StateMap{state: make(map[string]interface{}, 0)},
	}
	gui.views = make([]*View, 0)
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panic(err)
	}

	// Todo: add debug option to wait
	time.Sleep(1000 * time.Millisecond)

	gui.g = g
	gui.Configure(config)

	gui.g.SetManagerFunc(gui.layout)

	gui.BindAction("", Quit)

	for _, view := range views {
		view.gui = gui
		gui.views = append(gui.views, view)
	}

	return gui
}

func (gui *Gui) Configure(config config.GuiConfig) {
	gui.g.Highlight = config.Highlight
	gui.g.Cursor = config.Cursor
	gui.g.SelFgColor = config.SelFgColor

	gui.config = config
}

func (gui *Gui) Size() (int, int) {
	return gui.g.Size()
}

func (gui *Gui) MaxWidth() int {
	maxWidth, _ := gui.g.Size()
	return maxWidth
}

func (gui *Gui) MaxHeight() int {
	_, maxHeight := gui.g.Size()
	return maxHeight
}

func (gui *Gui) GetViews() []*View {
	return gui.views
}

func (gui *Gui) BindAction(viewName string, action *Action) {

	if viewName != "" {
		_, err := gui.g.View("v2")
		if err != nil {
			log.Panicln(err)
		}
	}

	if err := gui.g.SetKeybinding(
		viewName,
		action.Key,
		action.Mod,
		action.Handler(gui),
	); err != nil {
		log.Panic(err)
	}
}

func (gui *Gui) layout(g *gocui.Gui) error {
	for _, view := range gui.views {
		view.InitDimension()
		if v, err := g.SetView(
			view.Name,
			view.UpperLeftPointX(),
			view.UpperLeftPointY(),
			view.LowerRightPointX(),
			view.LowerRightPointY(),
		); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			if v != nil {
				view.v = v
				view.InitView()
				if view.Render != nil {
					if err := view.Render(gui, view); err != nil {
						return err
					}
				}
			}
		}
	}

	if gui.Render != nil {
		if err := gui.Render(gui); err != nil {
			return nil
		}
	}

	return nil
}

func (gui *Gui) Run() {
	if err := gui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (gui *Gui) Close() {
	gui.g.Close()
}

func (gui *Gui) GetView(name string) (*View, error) {
	if err := gui.ViewExisted(name); err != nil {
		return nil, err
	}

	return gui.getView(name), nil
}

func (gui *Gui) CurrentView() *View {
	v := gui.g.CurrentView()
	if v == nil {
		return nil
	}
	return gui.getView(v.Name())
}

func (gui *Gui) AddView(view *View) {
	// Todo: Check if view existed
	gui.views = append(gui.views, view)
}

func (gui *Gui) DeleteView(name string) error {
	if err := gui.ViewExisted(name); err != nil {
		return err
	}

	if err := gui.g.DeleteView(name); err != nil {
		return err
	}

	for index, view := range gui.views {
		if view.Name == name {
			gui.views = append(gui.views[:index], gui.views[index+1:]...)
		}
	}

	return nil
}

func (gui *Gui) ViewExisted(name string) error {
	_, err := gui.g.View(name)
	if err != nil {
		return err
	}
	return nil
}

func (gui *Gui) RenderString(viewName, s string) error {
	gui.Update(func(g *gocui.Gui) error {
		view, err := gui.GetView(viewName)
		if err != nil {
			return nil // return gracefully if view has been deleted
		}

		if err := view.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := view.SetCursor(0, 0); err != nil {
			return err
		}

		if view != nil {
			return view.SetViewContent(s)
		}

		return nil

	})
	return nil
}

func (gui *Gui) Update(f func(*gocui.Gui) error) {
	gui.g.Update(f)
}

func (gui *Gui) SetCurrentView(name string) (*View, error) {
	if _, err := gui.g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return gui.getView(name), nil
}

func (gui *Gui) SetViewOnTop(name string) (*View, error) {
	if _, err := gui.g.SetViewOnTop(name); err != nil {
		return nil, err
	}
	return gui.getView(name), nil
}

func (gui *Gui) getView(name string) *View {
	for _, view := range gui.views {
		if view.Name == name {
			return view
		}
	}
	return nil
}
