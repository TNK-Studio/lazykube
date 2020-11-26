package gui

import (
	"errors"
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"
	"sort"
)

type (
	// Gui Gui
	Gui struct {
		views   []*View
		Actions []*Action
		// History of focused views name.
		previousViews      TowHeadQueue
		OnSizeChange       func(gui *Gui) error
		OnRender           func(gui *Gui) error
		OnRenderOptions    func(gui *Gui) error
		previousViewsLimit int
		g                  *gocui.Gui
		state              State
		preHeight          int
		preWidth           int
		Config             config.GuiConfig
		renderTimes        int
	}
)

// NewGui NewGui
func NewGui(config config.GuiConfig, views ...*View) *Gui {
	gui := &Gui{
		state:              NewStateMap(),
		previousViews:      NewQueue(),
		previousViewsLimit: 20,
		Config:             config,
	}
	gui.views = make([]*View, 0)
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Logger.Panicf("%+v", err)
	}

	gui.g = g
	gui.Configure()

	gui.g.SetManagerFunc(gui.layout)

	gui.BindAction("", Quit)

	for _, view := range views {
		view.BindGui(gui)
		gui.views = append(gui.views, view)
	}

	return gui
}

// ReRenderAll ReRenderAll
func (gui *Gui) ReRenderAll() {
	gui.renderTimes += 1
	for _, view := range gui.views {
		view.ReRender()
	}
}

func (gui *Gui) layout(*gocui.Gui) error {
	height, width := gui.Size()
	if gui.preHeight != height || gui.preWidth != width {
		gui.preHeight = height
		gui.preWidth = width
		if gui.OnSizeChange != nil {
			if err := gui.OnSizeChange(gui); err != nil {
				return err
			}
		}
		gui.ReRenderAll()
		gui.SortViewsByZIndex()
		gui.SetAlwaysOnTopViews()
	}

	if err := gui.Clear(); err != nil {
		return err
	}
	for _, view := range gui.views {
		if err := gui.updateSelectedViewLine(view); err != nil {
			return err
		}

		err := gui.RenderView(view)
		if err == nil {
			continue
		}

		if errors.Is(err, ErrNotEnoughSpace) {
			if err := gui.renderNotEnoughSpaceView(); err != nil {
				return err
			}
			err = nil
		}

		return err
	}

	if gui.renderTimes < 0 {
		if err := gui.onRender(); err != nil {
			return err
		}
	}

	if err := gui.renderOptions(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) onRender() error {
	if gui.OnRender != nil && gui.renderTimes < 0 {
		if err := gui.OnRender(gui); err != nil {
			return nil
		}
	}

	currentView := gui.CurrentView()
	if currentView != nil {
		if _, err := gui.SetViewOnTop(currentView.Name); err != nil {
			return err
		}
	}
	gui.renderTimes += 1
	return nil
}

func (gui *Gui) updateSelectedViewLine(view *View) error {
	if !view.Rendered() {
		return nil
	}

	_, cy := view.Cursor()
	selectedLine, _ := view.Line(cy)
	if selectedLine != view.SelectedLine {
		view.SelectedLine = selectedLine
		if view.OnSelectedLineChange != nil {
			if err := view.OnSelectedLineChange(gui, view, selectedLine); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetAlwaysOnTopViews SetAlwaysOnTopViews
func (gui *Gui) SetAlwaysOnTopViews() {
	views := gui.views
	sort.Sort(ViewsZIndexSorter(views))

	for _, view := range views {
		if !view.AlwaysOnTop {
			continue
		}
		if _, err := gui.SetViewOnTop(view.Name); err != nil {
			continue
		}
	}

	return
}

// SortViewsByZIndex SortViewsByZIndex
func (gui *Gui) SortViewsByZIndex() {
	views := gui.views
	sort.Sort(ViewsZIndexSorter(views))

	for _, view := range views {
		if view.AlwaysOnTop {
			break
		}
		if _, err := gui.SetViewOnTop(view.Name); err != nil {
			continue
		}
	}

	return
}

// Configure Configure
func (gui *Gui) Configure() {
	gui.g.Highlight = gui.Config.Highlight
	gui.g.Cursor = gui.Config.Cursor
	gui.g.SelFgColor = gui.Config.SelFgColor
	gui.g.SelBgColor = gui.Config.SelBgColor
	gui.g.FgColor = gui.Config.FgColor
	gui.g.BgColor = gui.Config.BgColor
	gui.g.Mouse = gui.Config.Mouse
	gui.g.InputEsc = gui.Config.InputEsc

	inputMode := termbox.InputAlt
	if gui.g.InputEsc {
		inputMode = termbox.InputEsc
	}
	if gui.g.Mouse {
		inputMode |= termbox.InputMouse
	}
	termbox.SetInputMode(inputMode)

	// Must to set cursor otherwise hide cursor not work.
	termbox.SetCursor(0, 0)
	if !gui.g.Cursor {
		termbox.HideCursor()
	}
}

// Size Size
func (gui *Gui) Size() (int, int) {
	return gui.g.Size()
}

// MaxWidth MaxWidth
func (gui *Gui) MaxWidth() int {
	maxWidth, _ := gui.g.Size()
	return maxWidth
}

// MaxHeight MaxHeight
func (gui *Gui) MaxHeight() int {
	_, maxHeight := gui.g.Size()
	return maxHeight
}

// GetViews GetViews
func (gui *Gui) GetViews() []*View {
	return gui.views
}

// SetKeybinding SetKeybinding
func (gui *Gui) SetKeybinding(viewName string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
	if err := gui.g.SetKeybinding(
		viewName,
		key,
		mod,
		handler,
	); err != nil {
		log.Logger.Panicf("%+v", err)
	}
}

// BindAction BindAction
func (gui *Gui) BindAction(viewName string, action ActionInterface) {
	var handler ViewHandler
	if action.ReRenderAll() {
		handler = func(gui *Gui, view *View) error {
			if err := action.HandlerFunc(gui, view); err != nil {
				return err
			}
			gui.ReRenderAll()
			return nil
		}
	} else {
		handler = action.HandlerFunc
	}

	wrappedHandler := actionHandlerWrapper(gui, handler)
	if action.BindKey() != nil {
		gui.SetKeybinding(viewName,
			action.BindKey(),
			action.Modifier(),
			wrappedHandler,
		)
	}

	if action.BindKeys() != nil {
		for _, k := range action.BindKeys() {
			gui.SetKeybinding(viewName,
				k,
				action.Modifier(),
				wrappedHandler,
			)
		}
	}
}

// ViewDimensionValidated ViewDimensionValidated
func (gui *Gui) ViewDimensionValidated(x0, y0, x1, y1 int) bool {
	if x0 >= x1 || y0 >= y1 {
		return false
	}

	return true
}

// Run Run
func (gui *Gui) Run() {
	if gui.Actions != nil {
		for _, act := range gui.Actions {
			gui.BindAction("", act)
		}
	}

	for _, view := range gui.views {
		if view.Clickable {
			gui.BindAction(view.Name, ClickView)
		}
		if view.Actions != nil {
			for _, act := range view.Actions {
				gui.BindAction(view.Name, act)
			}
		}
	}

	if err := gui.g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Logger.Panicf("MainLoop - %+v", err)
	}
}

// Close Close
func (gui *Gui) Close() {
	gui.g.Close()
}

// GetView GetView
func (gui *Gui) GetView(name string) (*View, error) {
	if err := gui.ViewExisted(name); err != nil {
		return nil, err
	}

	return gui.getView(name), nil
}

// RenderView RenderView
func (gui *Gui) RenderView(view *View) error {
	x0, y0, x1, y1 := view.GetDimensions()
	if !gui.ViewDimensionValidated(x0, y0, x1, y1) {
		log.Logger.Warningf("View '%s' has not enough space to render. x0: %d, y0: %d, x1: %d, y1: %d", view.Name, x0, y0, x1, y1)
		return ErrNotEnoughSpace
	}
	return gui.renderView(view, x0, y0, x1, y1)
}

func (gui *Gui) unRenderNotEnoughSpaceView() error {
	v, _ := gui.g.View(NotEnoughSpace.Name)
	if v != nil {
		if err := gui.g.DeleteView(NotEnoughSpace.Name); err != nil {
			if errors.Is(err, gocui.ErrUnknownView) {
				return nil
			}
		}
	}
	return nil
}

// Clear Clear
func (gui *Gui) Clear() error {
	if err := gui.unRenderNotEnoughSpaceView(); err != nil {
		return err
	}
	if err := termbox.Clear(termbox.Attribute(gui.g.FgColor), termbox.Attribute(gui.g.BgColor)); err != nil {
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

// SetView SetView
func (gui *Gui) SetView(view *View, x0, y0, x1, y1 int) (*View, error) {
	if v, err := gui.g.SetView(
		view.Name,
		x0, y0, x1, y1,
	); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return nil, err
		}

		if v == nil {
			return nil, err
		}

		view.v = v
		view.x0, view.y0, view.x1, view.y1 = x0, y0, x1, y1
		view.InitView()
		return view, gocui.ErrUnknownView
	}
	return view, nil
}

func (gui *Gui) renderView(view *View, x0, y0, x1, y1 int) error {
	if _, err := gui.SetView(
		view,
		x0, y0, x1, y1,
	); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}

	if view != nil {
		if err := view.render(); err != nil {
			return err
		}
	}

	return nil
}

// ViewColors ViewColors
func (gui *Gui) ViewColors(view *View) (gocui.Attribute, gocui.Attribute) {
	if gui.Config.Highlight && view == gui.CurrentView() {
		return gui.Config.SelFgColor, gui.Config.SelBgColor
	}
	return gui.Config.FgColor, gui.Config.BgColor
}

// CurrentView CurrentView
func (gui *Gui) CurrentView() *View {
	v := gui.g.CurrentView()
	if v == nil {
		return nil
	}
	return gui.getView(v.Name())
}

// AddView AddView
func (gui *Gui) AddView(view *View) error {
	// Todo: Check if view existed
	gui.views = append(gui.views, view)
	view.gui = gui

	if view.Clickable {
		gui.BindAction(view.Name, ClickView)
	}

	if view.Actions != nil {
		for _, act := range view.Actions {
			gui.BindAction(view.Name, act)
		}
	}

	view.InitView()

	err := gui.RenderView(view)
	if errors.Is(err, ErrNotEnoughSpace) {
		if err := gui.renderNotEnoughSpaceView(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// DeleteView DeleteView
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

	gui.g.DeleteKeybindings(name)

	return nil
}

// ViewExisted ViewExisted
func (gui *Gui) ViewExisted(name string) error {
	_, err := gui.g.View(name)
	if err != nil {
		return err
	}
	return nil
}

// RenderString RenderString
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

// Update Update
func (gui *Gui) Update(f func(*gocui.Gui) error) {
	gui.g.Update(f)
}

// SetCurrentView SetCurrentView
func (gui *Gui) SetCurrentView(name string) (*View, error) {
	if _, err := gui.g.SetCurrentView(name); err != nil {
		return nil, err
	}
	view := gui.getView(name)
	return view, nil
}

// SetViewOnTop SetViewOnTop
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
	if !gui.previousViews.IsEmpty() {
		viewName := gui.previousViews.Pop().(string)
		log.Logger.Debugf("popPreviousView pop '%s', previousViews '%+v'", viewName, gui.previousViews)
		return viewName
	}

	return ""
}

func (gui *Gui) PeekPreviousView() string {
	if !gui.previousViews.IsEmpty() {
		return gui.previousViews.Peek().(string)
	}

	return ""
}

func (gui *Gui) pushPreviousView(name string) {
	if name == "" || name == gui.PeekPreviousView() {
		return
	}
	gui.previousViews.Push(name)
	if gui.previousViews.Len() > gui.previousViewsLimit {
		tail := gui.previousViews.PopTail()
		log.Logger.Debugf("pushPreviousView - previousViews over limit, pop tail '%s'", tail)
	}

	log.Logger.Debugf("pushPreviousView push '%s', previousViews '%+v'", name, gui.previousViews)
}

func (gui *Gui) FocusView(name string, canReturn bool) error {
	log.Logger.Debugf("FocusView - name: %s canReturn: %+v", name, canReturn)
	previousView := gui.CurrentView()

	gui.SortViewsByZIndex()

	if err := gui.focusView(name); err != nil {
		return err
	}

	gui.SetAlwaysOnTopViews()

	currentView := gui.CurrentView()

	if previousView != nil {
		if canReturn && !currentView.CanNotReturn && (currentView == nil || (previousView.Name != currentView.Name)) {
			gui.pushPreviousView(previousView.Name)
		}
		if previousView.Name != name {
			if err := currentView.focus(); err != nil {
				return err
			}
		}
	} else if currentView.OnFocus != nil {
		if err := currentView.focus(); err != nil {
			return err
		}
	}

	if previousView != nil && previousView.Name != name {
		if err := previousView.focusLost(); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) focusView(name string) error {
	if _, err := gui.SetCurrentView(name); err != nil {
		return err
	}
	if _, err := gui.SetViewOnTop(name); err != nil {
		return err
	}
	return nil
}

// HasPreviousView HasPreviousView
func (gui *Gui) HasPreviousView() bool {
	return !gui.previousViews.IsEmpty()
}

// ReturnPreviousView ReturnPreviousView
func (gui *Gui) ReturnPreviousView() error {
	previousViewName := gui.popPreviousView()
	previousView, err := gui.GetView(previousViewName)
	if err != nil {
		if errors.Is(err, gocui.ErrUnknownView) {
			log.Logger.Warningf("ReturnPreviousView view '%s' not found", previousViewName)
			return nil
		}
		return err
	}
	log.Logger.Debugf("ReturnPreviousView - gui.focusView(%s)", previousView.Name)
	return gui.FocusView(previousView.Name, false)
}

func (gui *Gui) renderOptions() error {
	currentView := gui.CurrentView()
	if gui.OnRenderOptions != nil {
		if err := gui.OnRenderOptions(gui); err != nil {
			return nil
		}
	}

	if currentView != nil {
		if err := currentView.renderOptions(); err != nil {
			return err
		}
	}
	return nil
}

// SetRune SetRune
func (gui *Gui) SetRune(x, y int, ch rune, fgColor, bgColor gocui.Attribute) error {
	return gui.g.SetRune(x, y, ch, fgColor, bgColor)
}

// ReRenderViews ReRenderViews
func (gui *Gui) ReRenderViews(viewNames ...string) {
	for _, name := range viewNames {
		view, err := gui.GetView(name)
		if err != nil {
			log.Logger.Warningf("ReRenderViews - view '%s' error %s", name, err)
			continue
		}

		view.ReRender()
	}
}

// ClearViews ClearViews
func (gui *Gui) ClearViews(viewNames ...string) {
	for _, name := range viewNames {
		view, err := gui.GetView(name)
		if err != nil {
			log.Logger.Warningf("ClearViews - view '%s' error %s", name, err)
			continue
		}

		view.Clear()
	}
}

// ForceFlush ForceFlush
func (gui *Gui) ForceFlush() error {
	termbox.Close()
	if err := termbox.Init(); err != nil {
		return err
	}
	inputMode := termbox.InputAlt
	if gui.g.InputEsc {
		inputMode = termbox.InputEsc
	}
	if gui.g.Mouse {
		inputMode |= termbox.InputMouse
	}
	termbox.SetInputMode(inputMode)

	// Must to set cursor otherwise hide cursor not work.
	termbox.SetCursor(0, 0)
	if !gui.g.Cursor {
		termbox.HideCursor()
	}

	return termbox.Flush()
}

func (gui *Gui) ReInitTermBox() error {
	termbox.Close()
	if err := termbox.Init(); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) SetState(key string, value interface{}, reRenderAll bool, reRenderViews ...string) error {
	err := gui.state.Set(key, value)
	if err != nil {
		return err
	}
	if reRenderAll {
		gui.ReRenderAll()
		return nil
	}

	if reRenderViews != nil {
		gui.ReRenderViews(reRenderViews...)
		return nil
	}
	return nil
}

func (gui *Gui) GetState(key string) (interface{}, error) {
	return gui.state.Get(key)
}
