package app

import (
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

const (
	filterInputViewName = "filterInput"
	filteredViewName    = "filtered"
)

var (
	toFilteredView = &guilib.Action{
		Name: "toFiltered",
		Keys: []interface{}{
			gocui.KeyTab,
			gocui.KeyArrowDown,
		},
		Handler: func(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
			return func(*gocui.Gui, *gocui.View) error {
				return gui.FocusView(filteredViewName, false)
			}
		},
		Mod: gocui.ModNone,
	}

	toFilterInputView = &guilib.Action{
		Name: "toFilterInput",
		Keys: []interface{}{
			gocui.KeyTab,
		},
		Handler: func(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
			return func(*gocui.Gui, *gocui.View) error {
				return gui.FocusView(filterInputViewName, false)
			}
		},
		Mod: gocui.ModNone,
	}
)

func newFilterDialog(title string, gui *guilib.Gui, resourceView *guilib.View) error {
	filterInput := &guilib.View{
		Name:         filterInputViewName,
		Title:        title,
		CanNotReturn: true,
		AlwaysOnTop:  true,
		Clickable:    true,
		Editable:     true,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			quarterWidth, quarterHeight := maxWidth/4, maxHeight/4
			x0 := quarterWidth
			x1 := quarterWidth * 3
			y0 := quarterHeight
			y1 := quarterHeight + 3
			return x0, y0, x1, y1
		},
		Actions: []*guilib.Action{
			toFilteredView,
		},
		OnFocusLost: func(gui *guilib.Gui, view *guilib.View) error {
			gui.Config.Cursor = false
			gui.Configure()
			return clearFilterDialog(gui, view)
		},
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			gui.Config.Cursor = true
			gui.Configure()
			return nil
		},
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			cx, _ := view.Cursor()
			cx = 0
			bufferLines := view.ViewBufferLines()
			filterStr := make([]rune, 0)
			if len(bufferLines) > 0 {
				filterStr = []rune(bufferLines[0])
				cx = len(filterStr) + 1
			}

			if err := view.SetCursor(cx, 0); err != nil {
				return err
			}
			view.ReRender()
			return nil
		},
	}
	filtered := &guilib.View{
		Name:         filteredViewName,
		CanNotReturn: true,
		AlwaysOnTop:  true,
		Clickable:    true,
		Actions: []*guilib.Action{
			toFilterInputView,
		},
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			quarterWidth, quarterHeight := maxWidth/4, maxHeight/4
			x0 := quarterWidth
			x1 := quarterWidth * 3
			y0 := quarterHeight + 2
			y1 := quarterHeight * 3
			return x0, y0, x1, y1
		},
		OnFocusLost: clearFilterDialog,
	}

	if err := gui.AddView(filterInput); err != nil {
		return err
	}
	if err := gui.AddView(filtered); err != nil {
		return err
	}

	if _, err := gui.SetViewOnTop(filterInput.Name); err != nil {
		return err
	}
	if _, err := gui.SetViewOnTop(filtered.Name); err != nil {
		return err
	}
	if err := gui.FocusView(filterInput.Name, true); err != nil {
		return err
	}
	return nil
}

func clearFilterDialog(gui *guilib.Gui, view *guilib.View) error {
	currentView := gui.CurrentView()

	if currentView != nil && (currentView.Name == filterInputViewName || currentView.Name == filteredViewName) {
		return nil
	}

	if err := gui.DeleteView(filterInputViewName); err != nil {
		return err
	}
	if err := gui.DeleteView(filteredViewName); err != nil {
		return err
	}
	return nil
}
