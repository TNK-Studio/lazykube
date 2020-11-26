package gui

import (
	"fmt"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/jroimartin/gocui"
)

var (
	NotEnoughSpace *View
)

func init() {
	NotEnoughSpace = &View{
		Name:  "notEnoughSpace",
		Title: "Not enough space to render.",
		DimensionFunc: func(gui *Gui, view *View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			return 0, 0, maxWidth - 1, maxHeight - 1
		},
		OnRender: func(gui *Gui, view *View) error {
			gui.Config.Cursor = false
			gui.Configure()
			view.ReRender()
			return nil
		},
	}
}

type (
	// View View
	View struct {
		Actions               []ActionInterface
		Name                  string
		Title                 string
		SelectedLine          string
		OnClick               ViewHandler
		OnLineClick           func(gui *Gui, view *View, cy int, lineString string) error
		OnRender              ViewHandler
		OnRenderOptions       ViewHandler
		OnFocus               ViewHandler
		OnFocusLost           ViewHandler
		OnCursorChange        func(gui *Gui, view *View, x, y int) error
		OnEditedChange        func(gui *Gui, view *View, key gocui.Key, ch rune, mod gocui.Modifier)
		OnSelectedLineChange  func(gui *Gui, view *View, selectedLine string) error
		DimensionFunc         DimensionFunc
		UpperLeftPointXFunc   ViewPointFunc
		UpperLeftPointYFunc   ViewPointFunc
		LowerRightPointXFunc  ViewPointFunc
		LowerRightPointYFunc  ViewPointFunc
		ZIndex                int
		x0                    int
		y0                    int
		x1                    int
		y1                    int
		gui                   *Gui
		v                     *gocui.View
		state                 State
		FgColor               gocui.Attribute
		BgColor               gocui.Attribute
		SelBgColor            gocui.Attribute
		SelFgColor            gocui.Attribute
		Clickable             bool
		Editable              bool
		Wrap                  bool
		Autoscroll            bool
		IgnoreCarriageReturns bool
		Highlight             bool
		NoFrame               bool
		MouseDisable          bool
		// When the "CanNotReturn" parameter is true, it will not be placed in previousViews where the view was clicked.
		CanNotReturn bool
		renderTimes  int
		AlwaysOnTop  bool
	}

	ViewHandler func(gui *Gui, view *View) error
)

// InitView InitView
func (view *View) InitView() {
	if view.state == nil {
		view.state = NewStateMap()
	}
	if view.v != nil {
		view.v.Title = view.Title
		view.v.Wrap = view.Wrap
		view.v.Editable = view.Editable
		view.v.Autoscroll = view.Autoscroll
		view.v.Highlight = view.Highlight
		view.v.Frame = !view.NoFrame
		view.v.FgColor = view.FgColor
		view.v.BgColor = view.BgColor
		view.v.SelBgColor = view.SelBgColor
		view.v.SelFgColor = view.SelFgColor
		view.v.MouseDisable = view.MouseDisable
		view.v.Editor = NewViewEditor(view.gui, view)
		view.v.OnCursorChange = view.onCursorChange
	}
}

// BindGui BindGui
func (view *View) BindGui(gui *Gui) {
	view.gui = gui
}

// InitDimension InitDimension
func (view *View) InitDimension() {
	if !view.IsBindingGui() {
		log.Logger.Warningf("Please run 'InitDimension' after binding Gui.")
		return
	}

	if view.DimensionFunc == nil {
		return
	}

	view.x0, view.y0, view.x1, view.y1 = view.DimensionFunc(view.gui, view)
}

// UpperLeftPointX UpperLeftPointX
func (view *View) UpperLeftPointX() int {
	if view.IsBindingGui() && view.UpperLeftPointXFunc != nil {
		return view.UpperLeftPointXFunc(view.gui, view)
	}
	return view.x0
}

// UpperLeftPointY UpperLeftPointY
func (view *View) UpperLeftPointY() int {
	if view.IsBindingGui() && view.UpperLeftPointYFunc != nil {
		return view.UpperLeftPointYFunc(view.gui, view)
	}
	return view.y0
}

// LowerRightPointX LowerRightPointX
func (view *View) LowerRightPointX() int {
	if view.IsBindingGui() && view.LowerRightPointXFunc != nil {
		return view.LowerRightPointXFunc(view.gui, view)
	}
	return view.x1
}

// LowerRightPointY LowerRightPointY
func (view *View) LowerRightPointY() int {
	if view.IsBindingGui() && view.LowerRightPointYFunc != nil {
		return view.LowerRightPointYFunc(view.gui, view)
	}
	return view.y1
}

// GetDimensions GetDimensions
func (view *View) GetDimensions() (int, int, int, int) {
	view.InitDimension()
	x0, y0, x1, y1 := view.UpperLeftPointX(), view.UpperLeftPointY(), view.LowerRightPointX(), view.LowerRightPointY()
	return x0, y0, x1, y1
}

// IsBindingGui IsBindingGui
func (view *View) IsBindingGui() bool {
	if view.gui != nil && view.gui.g != nil {
		return true
	}

	return false
}

// Rendered Rendered
func (view *View) Rendered() bool {
	return view.v != nil
}

// SetViewContent SetViewContent
func (view *View) SetViewContent(s string) error {
	view.v.Clear()
	if _, err := fmt.Fprint(view.v, utils.CleanString(s)); err != nil {
		return err
	}
	return nil
}

// SetOrigin SetOrigin
func (view *View) SetOrigin(x, y int) error {
	if view.Rendered() {
		return view.v.SetOrigin(x, y)
	}
	return nil
}

// Origin Origin
func (view *View) Origin() (int, int) {
	return view.v.Origin()
}

// SetCursor SetCursor
func (view *View) SetCursor(x, y int) error {
	if view.Rendered() {
		return view.v.SetCursor(x, y)
	}
	return nil
}

func (view *View) Write(p []byte) (n int, err error) {
	return view.v.Write(p)
}

// Clear Clear
func (view *View) Clear() {
	if view.Rendered() {
		view.v.Clear()
	}
}

// Cursor Cursor
func (view *View) Cursor() (int, int) {
	return view.v.Cursor()
}

// ViewBufferLines ViewBufferLines
func (view *View) ViewBufferLines() []string {
	return view.v.ViewBufferLines()
}

// ViewBuffer ViewBuffer
func (view *View) ViewBuffer() string {
	return view.v.ViewBuffer()
}

// Line Line
func (view *View) Line(y int) (string, error) {
	return view.v.Line(y)
}

// WhichLine WhichLine
func (view *View) WhichLine(s string) int {
	y := -1
	for index, line := range view.v.ViewBufferLines() {
		if line == s {
			return index
		}
	}
	return y
}

// MoveCursor MoveCursor
func (view *View) MoveCursor(dx, dy int, writeMode bool) {
	view.v.MoveCursor(dx, dy, writeMode)
}

// ReRender ReRender
func (view *View) ReRender() {
	view.renderTimes++
}

// ReRenderTimes ReRenderTimes
func (view *View) ReRenderTimes(times int) {
	view.renderTimes += times
}

func (view *View) render() error {
	if view.renderTimes < 0 {
		return nil
	}
	view.renderTimes--

	if view.OnRender != nil {
		if err := view.OnRender(view.gui, view); err != nil {
			return err
		}
	}
	return nil
}

func (view *View) renderOptions() error {
	if view.OnRenderOptions != nil {
		if err := view.OnRenderOptions(view.gui, view); err != nil {
			return nil
		}
	}
	return nil
}

func (view *View) focus() error {
	log.Logger.Debugf("view.focus - view name :%s", view.Name)
	if view.OnFocus != nil {
		log.Logger.Debugf("view.OnFocus - view name :%s", view.Name)
		if err := view.OnFocus(view.gui, view); err != nil {
			return nil
		}
	}
	return nil
}

func (view *View) focusLost() error {
	log.Logger.Debugf("view.focusLost - view name :%s", view.Name)
	if view.OnFocusLost != nil {
		log.Logger.Debugf("view.OnFocusLost - view name :%s", view.Name)
		if err := view.OnFocusLost(view.gui, view); err != nil {
			return nil
		}
	}
	return nil
}

func (view *View) Size() (int, int) {
	return view.v.Size()
}

// ResetCursorOrigin ResetCursorOrigin
func (view *View) ResetCursorOrigin() error {
	if err := view.v.SetCursor(0, 0); err != nil {
		return err
	}

	if err := view.v.SetOrigin(0, 0); err != nil {
		return err
	}

	return nil
}

func (view *View) onCursorChange(_ *gocui.View, x, y int) error {
	if view.OnCursorChange != nil {
		if err := view.OnCursorChange(view.gui, view, x, y); err != nil {
			return err
		}
	}
	return nil
}

func (view *View) SetState(key string, value interface{}, reRender bool) error {
	err := view.state.Set(key, value)
	if err != nil {
		return err
	}

	if reRender {
		view.ReRender()
	}
	return nil
}

func (view *View) GetState(key string) (interface{}, error) {
	return view.state.Get(key)
}

// DimensionFunc DimensionFunc
type DimensionFunc func(gui *Gui, view *View) (int, int, int, int)

// ViewPointFunc ViewPointFunc
type ViewPointFunc func(gui *Gui, view *View) int
