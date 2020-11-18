package app

import (
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	filterInputViewName      = "filterInput"
	filterInputValueStateKey = "value"
	filteredViewName         = "filtered"
	filteredNoResource       = "No Resource."
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

	filteredNextLine = &guilib.Action{
		Name: "filteredNextLine",
		Keys: []interface{}{
			gocui.KeyArrowDown,
		},
		Handler: nextLineHandler,
		Mod:     gocui.ModNone,
	}

	filteredPreviousLine = &guilib.Action{
		Name: "filteredPreviousLine",
		Keys: []interface{}{
			gocui.KeyArrowUp,
		},
		Handler: func(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
			return func(g *gocui.Gui, v *gocui.View) error {
				_, oy := v.Origin()
				_, cy := v.Cursor()
				if oy == 0 && cy-1 < 0 {
					return gui.FocusView(filterInputViewName, false)
				}
				return previousLineHandler(gui)(g, v)
			}
		},
		Mod: gocui.ModNone,
	}
)

func NewConfirmFilterInput(resourceViewName string) *guilib.Action {
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

				resourceView, err := gui.GetView(resourceViewName)
				if err != nil {
					return err
				}

				y := resourceView.WhichLine(filtered)
				if y < 0 {
					if err := resourceView.ResetCursorOrigin(); err != nil {
						return err
					}
				} else {
					if err := resourceView.SetOrigin(0, y); err != nil {
						return err
					}
					if err := resourceView.SetCursor(0, 0); err != nil {
						return err
					}
				}
				if err := closeFilterDialog(gui); err != nil {
					return err
				}
				if err := gui.ReturnPreviousView(); err != nil {
					return err
				}
				return nil
			}
		},
		Mod: gocui.ModNone,
	}
	return confirmFilterInput
}

func newFilterDialog(title string, gui *guilib.Gui, resourceViewName string) error {

	confirmFilterInput := NewConfirmFilterInput(resourceViewName)

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
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			gui.Config.Cursor = true
			gui.Configure()
			return nil
		},
		OnRenderOptions: filterDialogRenderOption,
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

			if err := view.State.Set(filterInputValueStateKey, value); err != nil {
				log.Logger.Warningf("OnEditedChange - view.State.Set(filterInputValueStateKey,%s) error %s", value, err)
				return
			}

			filteredView, err := gui.GetView(filteredViewName)
			if err != nil {
				log.Logger.Warningf("filteredView - gui.GetView(filteredViewName) error %s", err)
				return
			}
			filteredView.ReRender()
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
			filteredNextLine,
			filteredPreviousLine,
		},
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			value := ""
			val, _ := filterInput.State.Get(filterInputValueStateKey)
			if val != nil {
				value = val.(string)
			}

			view.Clear()
			if err := view.ResetCursorOrigin(); err != nil {
				log.Logger.Warningf("OnRender - view %s view.ResetCursorOrigin() error %s", filterInputViewName, err)
				return err
			}

			resourceView, err := gui.GetView(resourceViewName)
			if err != nil {
				return err
			}

			resourceList := resourceView.ViewBufferLines()
			if len(resourceList) == 0 {
				return nil
			}

			if value == "" {
				fmt.Fprint(view, strings.Join(resourceList[1:], "\n"))
				return nil
			}

			filtered := make([]string, 0)
			value = strings.ToLower(value)
			for _, resource := range resourceList[1:] {
				if strings.Index(strings.ToLower(resource), value) > -1 {
					filtered = append(filtered, resource)
				}
			}

			if len(filtered) == 0 {
				fmt.Fprint(view, filteredNoResource)
				return nil
			}

			fmt.Fprint(view, strings.Join(filtered, "\n"))
			return nil
		},
		OnRenderOptions: filterDialogRenderOption,
		OnFocusLost:     filterDialogFocusLost,
		OnLineClick: func(gui *guilib.Gui, view *guilib.View, cy int, lineString string) error {
			return nil
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

func filterDialogRenderOption(gui *guilib.Gui, view *guilib.View) error {
	return gui.RenderString(
		optionViewName,
		utils.OptionsMapToString(
			map[string]string{
				"← → ↑ ↓":   "navigate",
				"Ctrl+c":    "exit",
				"Esc":       "close dialog",
				"PgUp/PgDn": "scroll",
				"Home/End":  "top/bottom",
				"Tab":       "next panel",
				"Enter":     "confirm",
			}),
	)
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
