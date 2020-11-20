package app

import "github.com/jroimartin/gocui"

const (
	// All common actions name
	nextCyclicViewAction       = "nextCyclicView"
	backToPreviousViewAction   = "backToPreviousView"
	toNavigationAction         = "toNavigation"
	previousLineAction         = "previousLine"
	nextLineAction             = "nextLine"
	previousPageAction         = "previousPage"
	nextPageAction             = "nextPage"
	scrollUpAction             = "scrollUp"
	scrollDownAction           = "scrollDown"
	scrollTopAction            = "scrollTop"
	scrollBottomAction         = "scrollBottom"
	filterActionName           = "filterAction"
	editResourceActionName     = "Edit Resource"
	moreActionsName            = "moreActions"
	toFilteredViewAction       = "toFiltered"
	toFilterInputAction        = "toFilterInput"
	filteredNextLineAction     = "filteredNextLine"
	filteredPreviousLineAction = "filteredPreviousLine"
	confirmFilterInputAction   = "confirmFilterInput"
)

var (
	// All common actions key map.
	keyMap = map[string][]interface{}{
		nextCyclicViewAction:     {gocui.KeyTab},
		backToPreviousViewAction: {gocui.KeyEsc},
		toNavigationAction:       {gocui.KeyEnter, gocui.KeyArrowRight, 'l'},
		previousLineAction:       {gocui.KeyArrowUp, 'h'},
		nextLineAction:           {gocui.KeyArrowDown, 'j'},
		previousPageAction:       {gocui.KeyPgup},
		nextPageAction:           {gocui.KeyPgdn},
		scrollUpAction:           {gocui.MouseWheelUp},
		scrollDownAction:         {gocui.MouseWheelDown},
		scrollTopAction:          {gocui.KeyHome},
		scrollBottomAction:       {gocui.KeyEnd},
		filterActionName:         {gocui.KeyF4, 'f'},
		editResourceActionName:   {'e'},
		moreActionsName:          {gocui.KeyF3, 'm'},
		toFilteredViewAction:     {gocui.KeyTab, gocui.KeyArrowDown},
		toFilterInputAction:      {gocui.KeyTab},
		filteredNextLineAction:   {gocui.KeyArrowDown},
		confirmFilterInputAction: {gocui.KeyEnter},
	}
)
