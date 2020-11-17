package app

import (
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	filterInputViewName = "filterInput"
	filteredViewName    = "filtered"
	filteredNoResource  = "No Resource."
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
	resourceList := resourceView.ViewBufferLines()
	if len(resourceList) == 0 {
		return nil
	}

	confirmFilterInput := &guilib.Action{
		Name: "confirmFilterInput",
		Key:  gocui.KeyEnter,
		Handler: func(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
			return func(*gocui.Gui, *gocui.View) error {
				filteredView, err := gui.GetView(filteredViewName)

				if err != nil {
					return err
				}

				_, cy := filteredView.Cursor()
				filtered, _ := filteredView.Line(cy)
				if filtered == "" || filtered == filteredNoResource {
					return nil
				}

				y := resourceView.WhichLine(filtered)
				if y < 0 {
					if err := resourceView.ResetCursorOrigin(); err != nil {
						return err
					}
				} else {
					if err := setViewSelectedLine(gui, resourceView, filtered); err != nil {
						return err
					}
					if err := resourceView.SetOrigin(0, y); err != nil {
						return err
					}
				}
				if err := closeFilterDialog(gui); err != nil {
					return err
				}
				if err := gui.ReturnPreviousView(); err != nil {
					return err
				}
				gui.ReRender()
				return nil
			}
		},
		Mod: gocui.ModNone,
	}

	filterInput := &guilib.View{
		Name:         filterInputViewName,
		Title:        title,
		CanNotReturn: true,
		AlwaysOnTop:  true,
		Clickable:    true,
		Editable:     true,
		MouseDisable: true,
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
			confirmFilterInput,
		},
		OnFocusLost: func(gui *guilib.Gui, view *guilib.View) error {
			gui.Config.Cursor = false
			gui.Configure()
			return filterDialogFocusLost(gui, view)
		},
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			gui.Config.Cursor = true
			gui.Configure()
			return nil
		},
		OnEditedChange: func(gui *guilib.Gui, view *guilib.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			var value string
			bufferLines := view.ViewBufferLines()
			if len(bufferLines) > 0 {
				value = view.ViewBufferLines()[0]
			}
			filteredView, err := gui.GetView(filteredViewName)

			if err != nil {
				log.Logger.Warningf("OnEditedChange - gui.GetView(%s) error %s", filteredViewName, err)
				return
			}

			filteredView.Clear()
			if err := filteredView.ResetCursorOrigin(); err != nil {
				log.Logger.Warningf("OnEditedChange filteredView.ResetCursorOrigin() error %s", err)
				return
			}

			if value == "" {
				fmt.Fprint(filteredView, strings.Join(resourceList[1:], "\n"))
				return
			}

			filtered := make([]string, 0)
			value = strings.ToLower(value)
			for _, resource := range resourceList[1:] {
				if strings.Index(strings.ToLower(resource), value) > -1 {
					filtered = append(filtered, resource)
				}
			}

			if len(filtered) == 0 {
				fmt.Fprint(filteredView, filteredNoResource)
				return
			}

			fmt.Fprint(filteredView, strings.Join(filtered, "\n"))
		},
	}
	filtered := &guilib.View{
		Name:         filteredViewName,
		Clickable:    true,
		CanNotReturn: true,
		AlwaysOnTop:  true,
		Highlight:    true,
		SelFgColor:   gocui.ColorBlack,
		SelBgColor:   gocui.ColorWhite,
		Actions: []*guilib.Action{
			toFilterInputView,
			confirmFilterInput,
		},
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			fmt.Fprint(view, strings.Join(resourceList[1:], "\n"))
			return nil
		},
		OnFocusLost: filterDialogFocusLost,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			quarterWidth, quarterHeight := maxWidth/4, maxHeight/4
			x0 := quarterWidth
			x1 := quarterWidth * 3
			y0 := quarterHeight + 2
			y1 := quarterHeight * 3
			return x0, y0, x1, y1
		},
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

func filterDialogFocusLost(gui *guilib.Gui, view *guilib.View) error {
	currentView := gui.CurrentView()

	if currentView != nil && (currentView.Name == filterInputViewName || currentView.Name == filteredViewName) {
		return nil
	}

	if err := closeFilterDialog(gui); err != nil {
		return err
	}
	return nil
}

func closeFilterDialog(gui *guilib.Gui) error {
	if err := gui.DeleteView(filterInputViewName); err != nil {
		return err
	}
	if err := gui.DeleteView(filteredViewName); err != nil {
		return err
	}
	gui.Config.Cursor = false
	gui.Configure()
	return nil
}
