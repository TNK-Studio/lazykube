package app

import (
	"errors"
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

	moreActionsViewName           = "moreActions"
	moreActionTriggerViewStateKey = "triggerView"

	confirmDialogViewName = "confirmDialog"
)

var (
	confirmDialogOpt     = "Confirm"
	cancelDialogOpt      = "Cancel"
	confirmDialogOptions = []string{cancelDialogOpt, confirmDialogOpt}

	toFilteredView = &guilib.Action{
		Name: toFilteredViewAction,
		Keys: keyMap[toFilteredViewAction],
		Handler: func(gui *guilib.Gui, view *guilib.View) error {
			return gui.FocusView(filteredViewName, false)
		},
		Mod: gocui.ModNone,
	}

	toFilterInputView = &guilib.Action{
		Name: toFilterInputAction,
		Keys: keyMap[toFilterInputAction],
		Handler: func(gui *guilib.Gui, view *guilib.View) error {
			return gui.FocusView(filterInputViewName, false)
		},
		Mod: gocui.ModNone,
	}

	filteredNextLine = &guilib.Action{
		Name:    filteredNextLineAction,
		Keys:    keyMap[filteredNextLineAction],
		Handler: nextLineHandler,
		Mod:     gocui.ModNone,
	}

	filteredPreviousLine = &guilib.Action{
		Name: filteredPreviousLineAction,
		Keys: keyMap[filteredPreviousLineAction],
		Handler: func(gui *guilib.Gui, view *guilib.View) error {
			_, oy := view.Origin()
			_, cy := view.Cursor()
			if oy == 0 && cy-1 < 0 {
				return gui.FocusView(filterInputViewName, false)
			}
			return previousLineHandler(gui, view)
		},
		Mod: gocui.ModNone,
	}
)

//nolint:gocognit
func newConfirmFilterInput(resourceViewName string) *guilib.Action {
	confirmFilterInput := &guilib.Action{
		Name: confirmFilterInputAction,
		Keys: keyMap[confirmFilterInputAction],
		Handler: func(gui *guilib.Gui, _ *guilib.View) error {
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
		},
		Mod: gocui.ModNone,
	}
	return confirmFilterInput
}

//nolint:funlen
//nolint:gocognit
//nolint:gocognit
//nolint:gocognit
func newFilterDialog(title string, gui *guilib.Gui, resourceViewName string) error {

	confirmFilterInput := newConfirmFilterInput(resourceViewName)

	filterInput := &guilib.View{
		Name:         filterInputViewName,
		Title:        title,
		CanNotReturn: true,
		AlwaysOnTop:  true,
		ZIndex:       0,
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
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toFilteredView,
			confirmFilterInput,
		}),
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
		ZIndex:       1,
		Clickable:    true,
		CanNotReturn: true,
		AlwaysOnTop:  true,
		Highlight:    true,
		SelFgColor:   gocui.ColorBlack,
		SelBgColor:   gocui.ColorWhite,
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			toFilterInputView,
			confirmFilterInput,
			filteredNextLine,
			filteredPreviousLine,
		}),
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
				_, err := fmt.Fprint(view, strings.Join(resourceList[1:], "\n"))
				if err != nil {
					return err
				}

				return nil
			}

			filtered := make([]string, 0)
			value = strings.ToLower(value)
			for _, resource := range resourceList[1:] {
				if strings.Contains(strings.ToLower(resource), value) {
					filtered = append(filtered, resource)
				}
			}

			if len(filtered) == 0 {
				_, err := fmt.Fprint(view, filteredNoResource)
				if err != nil {
					return err
				}
				return nil
			}

			_, err = fmt.Fprint(view, strings.Join(filtered, "\n"))
			if err != nil {
				return err
			}

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
	filterInput.InitView()
	filtered.InitView()

	if err := gui.AddView(filterInput); err != nil {
		return err
	}
	if err := gui.AddView(filtered); err != nil {
		return err
	}
	if err := gui.FocusView(filterInput.Name, true); err != nil {
		return err
	}
	return nil
}

func filterDialogRenderOption(gui *guilib.Gui, _ *guilib.View) error {
	return gui.RenderString(
		optionViewName,
		utils.OptionsMapToString(
			map[string]string{
				"←→↑↓":      "navigate",
				"Ctrl+c":    "exit",
				"Esc":       "close dialog",
				"PgUp/PgDn": "scroll",
				"Home/End":  "top/bottom",
				"Tab":       "next panel",
				"Enter":     "confirm",
			}),
	)
}

func filterDialogFocusLost(gui *guilib.Gui, _ *guilib.View) error {
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

func newMoreActionDialog(title string, gui *guilib.Gui, view *guilib.View, moreActions []*moreAction) error {
	moreActionView := &guilib.View{
		Title:       title,
		Name:        moreActionsViewName,
		AlwaysOnTop: true,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			quarterWidth, quarterHeight := maxWidth/4, maxHeight/4
			x0 := quarterWidth
			x1 := quarterWidth * 3
			y0 := quarterHeight
			y1 := quarterHeight * 3
			return x0, y0, x1, y1
		},
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			view.Clear()
			moreActionsDescription := make([]string, 0)
			for _, moreAct := range moreActions {
				if moreAct.NeedSelectResource {
					resourceView, err := getMoreActionTriggerView(view)
					if err != nil {
						continue
					}

					_, resourceName, err := getResourceNamespaceAndName(gui, resourceView)
					if err != nil {
						continue
					}

					if notResourceSelected(resourceName) {
						continue
					}
				}
				moreActionsDescription = append(moreActionsDescription, keyMapDescription(moreAct.Keys, moreAct.Name))
			}

			if len(moreActionsDescription) == 0 {
				_, err := fmt.Fprint(view, "No more actions.")
				if err != nil {
					return err
				}
				return nil
			}

			_, err := fmt.Fprint(view, strings.Join(moreActionsDescription, "\n"))
			if err != nil {
				return err
			}
			return nil
		},
		OnFocusLost: func(gui *guilib.Gui, view *guilib.View) error {
			if err := gui.DeleteView(view.Name); err != nil {
				return err
			}
			return nil
		},
		Actions: toMoreActionArr(moreActions),
	}

	if err := gui.AddView(moreActionView); err != nil {
		return err
	}

	if err := moreActionView.State.Set(moreActionTriggerViewStateKey, view); err != nil {
		return err
	}
	// Todo: On view state change. Rerender.
	moreActionView.ReRender()

	if err := gui.FocusView(moreActionView.Name, true); err != nil {
		return err
	}
	return nil
}

func getMoreActionTriggerView(moreActionView *guilib.View) (*guilib.View, error) {
	val, err := moreActionView.State.Get(moreActionTriggerViewStateKey)
	if err != nil {
		return nil, err
	}
	view, ok := val.(*guilib.View)
	if !ok {
		return nil, errors.New("editResourceHandler - more action trigger view not found. ")
	}
	return view, nil
}
