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
	}
}

type View struct {
	State State
	Name  string
	Title string

	Clickable   bool
	OnClick     func(gui *Gui, view *View) error
	OnLineClick func(gui *Gui, view *View, cy int, lineString string) error

	Editable              bool
	Wrap                  bool
	Autoscroll            bool
	IgnoreCarriageReturns bool
	Highlight             bool
	NoFrame               bool
	FgColor               gocui.Attribute
	BgColor               gocui.Attribute
	SelBgColor            gocui.Attribute
	SelFgColor            gocui.Attribute

	// When the "CanNotReturn" parameter is true, it will not be placed in previousViews
	CanNotReturn bool

	reRendered bool

	OnRender        func(gui *Gui, view *View) error
	OnRenderOptions func(gui *Gui, view *View) error

	OnFocus     func(gui *Gui, view *View) error
	OnFocusLost func(gui *Gui, view *View) error

	DimensionFunc DimensionFunc

	UpperLeftPointXFunc  ViewPointFunc
	UpperLeftPointYFunc  ViewPointFunc
	LowerRightPointXFunc ViewPointFunc
	LowerRightPointYFunc ViewPointFunc

	x0, y0, x1, y1 int

	gui *Gui
	v   *gocui.View
}

func (view *View) InitView() {
	if view.State == nil {
		view.State = NewStateMap()
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
	}
}

func (view *View) BindGui(gui *Gui) {
	view.gui = gui
}

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

func (view *View) UpperLeftPointX() int {
	if view.IsBindingGui() && view.UpperLeftPointXFunc != nil {
		return view.UpperLeftPointXFunc(view.gui, view)
	}
	return view.x0
}

func (view *View) UpperLeftPointY() int {
	if view.IsBindingGui() && view.UpperLeftPointYFunc != nil {
		return view.UpperLeftPointYFunc(view.gui, view)
	}
	return view.y0
}

func (view *View) LowerRightPointX() int {
	if view.IsBindingGui() && view.LowerRightPointXFunc != nil {
		return view.LowerRightPointXFunc(view.gui, view)
	}
	return view.x1
}

func (view *View) LowerRightPointY() int {
	if view.IsBindingGui() && view.LowerRightPointYFunc != nil {
		return view.LowerRightPointYFunc(view.gui, view)
	}
	return view.y1
}

func (view *View) GetDimensions() (int, int, int, int) {
	view.InitDimension()
	x0, y0, x1, y1 := view.UpperLeftPointX(), view.UpperLeftPointY(), view.LowerRightPointX(), view.LowerRightPointY()
	return x0, y0, x1, y1
}

func (view *View) IsBindingGui() bool {
	if view.gui != nil && view.gui.g != nil {
		return true
	}

	return false
}

func (view *View) Rendered() bool {
	if view.v != nil {
		return true
	}

	return false
}

func (view *View) SetViewContent(s string) error {
	view.v.Clear()
	if _, err := fmt.Fprint(view.v, utils.CleanString(s)); err != nil {
		return err
	}
	return nil
}

func (view *View) SetOrigin(x, y int) error {
	if view.Rendered() {
		return view.v.SetOrigin(x, y)
	}
	return nil
}

func (view *View) SetCursor(x, y int) error {
	if view.Rendered() {
		return view.v.SetCursor(x, y)
	}
	return nil
}

func (view *View) Write(p []byte) (n int, err error) {
	return view.v.Write(p)
}

func (view *View) Clear() {
	view.v.Clear()
}

func (view *View) Cursor() (int, int) {
	return view.v.Cursor()
}

func (view *View) ViewBufferLines() []string {
	return view.v.ViewBufferLines()
}

func (view *View) Line(y int) (string, error) {
	return view.v.Line(y)
}

func (view *View) MoveCursor(dx, dy int, writeMode bool) {
	view.v.MoveCursor(dx, dy, writeMode)
}

func (view *View) ReRender() {
	view.reRendered = false
}

func (view *View) render() error {
	if view.reRendered {
		return nil
	}
	view.reRendered = true

	if view.OnRender != nil {
		if err := view.OnRender(view.gui, view); err != nil {
			return nil
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

type DimensionFunc func(gui *Gui, view *View) (int, int, int, int)
type ViewPointFunc func(gui *Gui, view *View) int
