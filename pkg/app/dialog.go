package app

import (
	"errors"
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/gookit/color"
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

	optionsDialogViewName = "optionsDialog"
	inputDialogViewName   = "inputDialog"
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

// Show dialog functions

func showFilterDialog(gui *guilib.Gui, title string, confirmHandler func(string) error, dataFunc func() ([]string, error), noResultMsg string) error {
	var filterInput, filtered *guilib.View
	// If views existed.
	filterInput, _ = gui.GetView(filterInputViewName)
	if filterInput != nil {
		return nil
	}
	filtered, _ = gui.GetView(filteredViewName)
	if filtered != nil {
		return nil
	}

	filterInput, filtered = newFilterDialog(title, confirmHandler, dataFunc, noResultMsg)
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

func showMoreActionDialog(gui *guilib.Gui, view *guilib.View, title string, moreActions []*moreAction) error {
	var moreActionView *guilib.View
	// If more action view existed.
	moreActionView, _ = gui.GetView(moreActionsViewName)
	if moreActionView != nil {
		return nil
	}

	moreActionView = newMoreActionDialog(title, moreActions)
	if err := gui.AddView(moreActionView); err != nil {
		return err
	}

	if err := moreActionView.SetState(moreActionTriggerViewStateKey, view); err != nil {
		return err
	}

	if err := gui.FocusView(moreActionView.Name, true); err != nil {
		return err
	}
	return nil
}

func showConfirmActionDialog(gui *guilib.Gui, title, relatedViewName string, handler guilib.ViewHandler) error {
	var confirmDialog *guilib.View
	// If view existed.
	confirmDialog, _ = gui.GetView(confirmDialogViewName)
	if confirmDialog != nil {
		return nil
	}

	confirmDialog = newConfirmActionDialog(title, relatedViewName, handler)
	if err := gui.AddView(confirmDialog); err != nil {
		return err
	}
	return nil
}

func showOptionsDialog(gui *guilib.Gui, title string, zIndex int, confirmHandler func(string) error, optionsFunc func() []string) error {
	var optionsDialog *guilib.View
	// If view existed.
	optionsDialog, _ = gui.GetView(optionsDialogViewName)
	if optionsDialog != nil {
		return nil
	}

	optionsDialog = newOptionsDialog(title, zIndex, confirmHandler, optionsFunc)
	if err := gui.AddView(optionsDialog); err != nil {
		return err
	}
	if err := gui.FocusView(optionsDialogViewName, false); err != nil {
		return err
	}
	return nil
}

func showInputDialog(gui *guilib.Gui, title string, zIndex int, confirmHandler func(string) error, defaultValue string) error {
	var inputDialog *guilib.View
	// If view existed.
	inputDialog, _ = gui.GetView(inputDialogViewName)
	if inputDialog != nil {
		return nil
	}

	inputDialog = newInputDialog(title, zIndex, confirmHandler, defaultValue)
	if err := gui.AddView(inputDialog); err != nil {
		return err
	}
	if err := gui.FocusView(inputDialogViewName, false); err != nil {
		return err
	}
	return nil
}

// New dialog functions

func newFilterDialog(title string, confirmHandler func(string) error, dataFunc func() ([]string, error), noResultMsg string) (*guilib.View, *guilib.View) {
	confirmAction := newConfirmFilterInput(confirmHandler)
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
			confirmAction,
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
			// Todo: fix character "_"
			var value string
			bufferLines := view.ViewBufferLines()
			if len(bufferLines) > 0 {
				value = view.ViewBufferLines()[0]
			}

			if err := view.SetState(filterInputValueStateKey, value); err != nil {
				log.Logger.Warningf("OnEditedChange - view.SetState(filterInputValueStateKey,%s) error %s", value, err)
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
			confirmAction,
			filteredNextLine,
			filteredPreviousLine,
		}),
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			value := ""
			val, _ := filterInput.GetState(filterInputValueStateKey)
			if val != nil {
				value = val.(string)
			}

			view.Clear()
			if err := view.ResetCursorOrigin(); err != nil {
				log.Logger.Warningf("OnRender - view %s view.ResetCursorOrigin() error %s", filterInputViewName, err)
				return err
			}

			data, err := dataFunc()
			if err != nil {
				return err
			}

			if len(data) == 0 {
				return nil
			}

			if value == "" {
				_, err := fmt.Fprint(view, strings.Join(data[1:], "\n"))
				if err != nil {
					return err
				}

				return nil
			}

			filtered := make([]string, 0)
			value = strings.TrimSpace(strings.ToLower(value))
			for _, resource := range data[1:] {
				if strings.Contains(strings.ToLower(resource), value) {
					filtered = append(filtered, resource)
				}
			}

			if len(filtered) == 0 {
				_, err := fmt.Fprint(view, noResultMsg)
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
			return confirmHandler(lineString)
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
	return filterInput, filtered
}

func newMoreActionDialog(title string, moreActions []*moreAction) *guilib.View {
	moreActionDialog := &guilib.View{
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
	}

	moreActionDialog.Actions = make([]guilib.ActionInterface, 0)
	for _, each := range moreActions {
		action := each.Action
		action.Handler = moreActionHandlerWrapper(action.Handler)
		moreActionDialog.Actions = append(
			moreActionDialog.Actions,
			&moreAction{
				NeedSelectResource: each.NeedSelectResource,
				Action:             action,
			},
		)
	}
	return moreActionDialog
}

func newConfirmActionDialog(title, relatedViewName string, handler guilib.ViewHandler) *guilib.View {
	return &guilib.View{
		Name:         confirmDialogViewName,
		Title:        title,
		Clickable:    true,
		CanNotReturn: true,
		AlwaysOnTop:  true,
		ZIndex:       10,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			halfWidth, halfHeight := maxWidth/2, maxHeight/2
			eighthWidth, eighthHeight := maxWidth/8, maxHeight/8
			return halfWidth - eighthWidth*2, halfHeight - eighthHeight, halfWidth + eighthHeight*2, halfHeight + eighthHeight
		},
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			view.Clear()
			canReturn := true
			currentView := gui.CurrentView()
			if currentView != nil && currentView.Name == moreActionsViewName {
				canReturn = false
			}
			if err := gui.FocusView(confirmDialogViewName, canReturn); err != nil {
				return err
			}

			var value string
			val, err := view.GetState("value")
			if err != nil {
				_ = view.SetState("value", confirmDialogOpt)
				value = confirmDialogOpt
			} else {
				value = val.(string)
			}

			optionsStr := strings.Join(confirmDialogOptions, optSeparator)
			length := len([]rune(optionsStr))
			width, height := view.Size()
			if width < length {
				return guilib.ErrNotEnoughSpace
			}

			offset := (width - length) / 2
			for i := 0; i <= height/3*2; i++ {
				if _, err := fmt.Fprintln(view); err != nil {
					return err
				}
			}

			optionsStr = strings.Replace(optionsStr, value, color.Green.Sprint(value), 1)
			for i := 0; i < offset; i++ {
				optionsStr = " " + optionsStr
			}
			log.Logger.Debugf("confirmDialogView - optionsStr %s", optionsStr)
			if _, err := fmt.Fprint(view, optionsStr); err != nil {
				return err
			}
			return nil
		},
		OnFocusLost: func(gui *guilib.Gui, view *guilib.View) error {
			if err := gui.DeleteView(view.Name); err != nil {
				return err
			}
			if err := gui.ReturnPreviousView(); err != nil {
				return err
			}
			return nil
		},
		OnLineClick: func(gui *guilib.Gui, view *guilib.View, cy int, lineString string) error {
			if strings.ReplaceAll(lineString, " ", "") == "" {
				return nil
			}
			cx, _ := view.Cursor()
			optionsStr := strings.Join(confirmDialogOptions, optSeparator)
			length := len([]rune(optionsStr))
			width, _ := view.Size()
			offset := (width - length) / 2
			optIndex, selected := utils.ClickOption(confirmDialogOptions, optSeparator, cx, offset)
			if optIndex < 0 {
				return nil
			}
			return confirmDialogOptionHandler(gui, view, relatedViewName, selected, handler)
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			{
				Keys: keyMap[switchConfirmDialogOpt],
				Name: switchConfirmDialogOpt,
				Handler: func(gui *guilib.Gui, view *guilib.View) error {
					var value string
					val, err := view.GetState("value")
					if err != nil {
						_ = view.SetState("value", confirmDialogOpt)
						value = confirmDialogOpt
					} else {
						value = val.(string)
					}

					if value == confirmDialogOpt {
						_ = view.SetState("value", cancelDialogOpt)
					} else {
						_ = view.SetState("value", confirmDialogOpt)
					}
					view.ReRender()
					return nil
				},
				ReRenderAllView: false,
				Mod:             gocui.ModNone,
			},
			{
				Keys: keyMap[confirmDialogEnter],
				Name: confirmDialogEnter,
				Handler: func(gui *guilib.Gui, view *guilib.View) error {
					var value string
					val, err := view.GetState("value")
					if err != nil {
						_ = view.SetState("value", confirmDialogOpt)
						value = confirmDialogOpt
					} else {
						value = val.(string)
					}

					return confirmDialogOptionHandler(gui, view, relatedViewName, value, handler)
				},
				Mod: gocui.ModNone,
			},
		}),
	}
}

func newOptionsDialog(title string, zIndex int, confirmHandler func(string) error, optionsFunc func() []string) *guilib.View {
	return &guilib.View{
		Name:         optionsDialogViewName,
		Title:        title,
		Clickable:    true,
		CanNotReturn: false,
		AlwaysOnTop:  true,
		ZIndex:       zIndex,
		Highlight:    true,
		SelFgColor:   gocui.ColorBlack,
		SelBgColor:   gocui.ColorWhite,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			halfWidth, halfHeight := maxWidth/2, maxHeight/2
			eighthWidth, eighthHeight := maxWidth/8, maxHeight/8
			return halfWidth - eighthWidth*2, halfHeight - eighthHeight, halfWidth + eighthHeight*2, halfHeight + eighthHeight
		},
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			view.Clear()
			options := optionsFunc()
			for _, option := range options {
				if _, err := fmt.Fprintln(view, option); err != nil {
					return err
				}
			}
			return nil
		},
		OnLineClick: func(gui *guilib.Gui, view *guilib.View, cy int, lineString string) error {
			return confirmHandler(lineString)
		},
		OnFocusLost: func(gui *guilib.Gui, view *guilib.View) error {
			if err := gui.DeleteView(view.Name); err != nil {
				return err
			}
			return nil
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			{
				Keys: keyMap[optionsDialogEnter],
				Name: optionsDialogEnter,
				Handler: func(gui *guilib.Gui, view *guilib.View) error {
					return confirmHandler(view.SelectedLine)
				},
				ReRenderAllView: false,
				Mod:             gocui.ModNone,
			},
		}),
	}
}

func newInputDialog(title string, zIndex int, confirmHandler func(string) error, defaultValue string) *guilib.View {
	return &guilib.View{
		Name:         inputDialogViewName,
		Title:        title,
		Clickable:    true,
		CanNotReturn: false,
		AlwaysOnTop:  true,
		Editable:     true,
		MouseDisable: true,
		ZIndex:       zIndex,
		DimensionFunc: func(gui *guilib.Gui, view *guilib.View) (int, int, int, int) {
			maxWidth, maxHeight := gui.Size()
			quarterWidth, quarterHeight := maxWidth/4, maxHeight/4
			x0 := quarterWidth
			x1 := quarterWidth * 3
			y0 := quarterHeight
			y1 := quarterHeight + 3
			return x0, y0, x1, y1
		},
		OnRender: func(gui *guilib.Gui, view *guilib.View) error {
			if view.ViewBuffer() == "" && defaultValue != "" {
				if _, err := fmt.Fprint(view, defaultValue); err != nil {
					return err
				}
				dx := len([]rune(defaultValue))
				view.MoveCursor(dx, 0, true)
			}
			return nil
		},
		OnFocus: func(gui *guilib.Gui, view *guilib.View) error {
			gui.Config.Cursor = true
			gui.Configure()
			return nil
		},
		OnFocusLost: func(gui *guilib.Gui, view *guilib.View) error {
			gui.Config.Cursor = false
			gui.Configure()
			if err := gui.DeleteView(view.Name); err != nil {
				return err
			}
			return nil
		},
		Actions: guilib.ToActionInterfaceArr([]*guilib.Action{
			{
				Keys: keyMap[inputDialogEnter],
				Name: inputDialogEnter,
				Handler: func(gui *guilib.Gui, view *guilib.View) error {
					return confirmHandler(view.SelectedLine)
				},
				ReRenderAllView: false,
				Mod:             gocui.ModNone,
			},
		}),
	}
}

// New dialog function utils.

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

func getMoreActionTriggerView(moreActionView *guilib.View) (*guilib.View, error) {
	val, err := moreActionView.GetState(moreActionTriggerViewStateKey)
	if err != nil {
		return nil, err
	}
	view, ok := val.(*guilib.View)
	if !ok {
		return nil, errors.New("getMoreActionTriggerView - more action trigger view not found. ")
	}
	return view, nil
}

func moreActionHandlerWrapper(handler guilib.ViewHandler) guilib.ViewHandler {
	return func(gui *guilib.Gui, view *guilib.View) error {
		triggerView := view
		var err error
		triggerView, err = getMoreActionTriggerView(view)
		if err != nil {
			return err
		}
		return handler(gui, triggerView)
	}
}
