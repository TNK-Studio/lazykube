package gui

import (
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/golang-collections/collections/stack"
	"github.com/jroimartin/gocui"
	"log"
	"time"
)

type Gui struct {
	State  State
	Render func(gui *Gui) error
	RenderOptions func(gui *Gui) error

	// History of focused views name.
	previousViews *stack.Stack

	g      *gocui.Gui
	views  []*View
	config config.GuiConfig
}

func NewGui(config config.GuiConfig, views ...*View) *Gui {

	gui := &Gui{
		State:         &StateMap{state: make(map[string]interface{}, 0)},
		previousViews: stack.New(),
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
		view.BindGui(gui)
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
	if err := gui.Clear(); err != nil {
		return err
	}
	for _, view := range gui.views {
		err := gui.RenderView(view)
		if err == nil {
			continue
		}

		if err == ErrNotEnoughSpace {
			if err := gui.renderNotEnoughSpaceView(); err != nil {
				return err
			}
			err = nil
		}

		return err
	}

	if gui.Render != nil {
		if err := gui.Render(gui); err != nil {
			return nil
		}
	}

	if err := gui.renderOptions(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) ViewDimensionValidated(x0, y0, x1, y1 int) bool {
	if x0 >= x1 || y0 >= y1 {
		return false
	}

	return true
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

func (gui *Gui) RenderView(view *View) error {
	x0, y0, x1, y1 := view.GetDimensions()
	if !gui.ViewDimensionValidated(x0, y0, x1, y1) {
		view.v = nil
		return ErrNotEnoughSpace
	}
	return gui.renderView(view, x0, y0, x1, y1)
}

func (gui *Gui) unRenderNotEnoughSpaceView() error {
	v, _ := gui.g.View(NotEnoughSpace.Name)
	if v != nil {
		NotEnoughSpace.v = nil
		return gui.g.DeleteView(NotEnoughSpace.Name)
	}
	return nil
}

func (gui *Gui) Clear() error {
	if err := gui.unRenderNotEnoughSpaceView(); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) renderNotEnoughSpaceView() error {
	NotEnoughSpace.BindGui(gui)
	x0, y0, x1, y1 := NotEnoughSpace.GetDimensions()
	if !gui.ViewDimensionValidated(x0, y0, x1, y1) {
		return nil
	}
	return gui.renderView(NotEnoughSpace, x0, y0, x1, y1)
}

func (gui *Gui) renderView(view *View, x0, y0, x1, y1 int) error {
	if v, err := gui.g.SetView(
		view.Name,
		x0, y0, x1, y1,
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
	return nil
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

func (gui *Gui) popPreviousView() string {
	if gui.previousViews.Len() > 0 {
		return gui.previousViews.Pop().(string)
	}

	return ""
}

func (gui *Gui) peekPreviousView() string {
	if gui.previousViews.Len() > 0 {
		return gui.previousViews.Peek().(string)
	}

	return ""
}

func (gui *Gui) pushPreviousView(name string) {
	gui.previousViews.Push(name)
}

func (gui *Gui) FocusView(name string) error {
	if _, err := gui.g.SetCurrentView(name); err != nil {
		return err
	}
	if _, err := gui.g.SetViewOnTop(name); err != nil {
		return err
	}
	return nil
}

// pass in oldView = nil if you don't want to be able to return to your old view
// TODO: move some of this logic into our onFocusLost and onFocus hooks
func (gui *Gui) SwitchFocus(oldViewName, newViewName string, returning bool) error {
	// we assume we'll never want to return focus to a popup panel i.e.
	// we should never stack popup panels
	//if oldView != nil && !gui.isPopupPanel(oldView.Name()) && !returning {
	//	gui.pushPreviousView(oldView.Name())
	//}
	if oldViewName != "" && !returning {
		gui.pushPreviousView(oldViewName)
	}

	//gui.Log.Info("setting highlight to true for view " + newView.Name())
	//gui.Log.Info("new focused view is " + newView.Name())
	if err := gui.FocusView(newViewName); err != nil {
		return err
	}

	//g.Cursor = newView.Editable

	//return gui.newLineFocused(newView)
	return nil
}


func (gui *Gui) renderOptions() error {
	currentView := gui.CurrentView()
	if currentView != nil && currentView.RenderOptions != nil {
		return currentView.RenderOptions(gui, currentView)
	}

	if gui.RenderOptions != nil {
		return gui.RenderOptions(gui)
	}
	return nil
}