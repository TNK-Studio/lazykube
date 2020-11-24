package app

import (
	"fmt"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

const (
	// All common actions name
	nextCyclicViewAction       = "nextCyclicView"
	backToPreviousViewAction   = "backToPreviousView"
	toNavigationAction         = "toNavigation"
	navigationArrowLeft        = "navigationArrowLeft"
	navigationArrowRight       = "navigationArrowRight"
	navigationDown             = "navigationDown"
	detailToNavigation         = "detailToNavigation"
	detailArrowUp              = "detailArrowUp"
	detailArrowDown            = "detailArrowDown"
	previousLineAction         = "previousLine"
	nextLineAction             = "nextLine"
	previousPageAction         = "previousPage"
	nextPageAction             = "nextPage"
	scrollUpAction             = "scrollUp"
	scrollDownAction           = "scrollDown"
	scrollTopAction            = "scrollTop"
	scrollBottomAction         = "scrollBottom"
	filterResourceActionName   = "filterResourceAction"
	moreActionsName            = "moreActions"
	toFilteredViewAction       = "toFiltered"
	toFilterInputAction        = "toFilterInput"
	filteredNextLineAction     = "filteredNextLine"
	filteredPreviousLineAction = "filteredPreviousLine"
	confirmFilterInputAction   = "confirmFilterInput"
	switchConfirmDialogOpt     = "switchConfirmDialogOpt"
	confirmDialogEnter         = "confirmDialogEnter"
	optionsDialogEnter         = "optionsDialogEnter"
	inputDialogEnter           = "inputDialogEnter"

	// More actions
	editResourceActionName              = "Edit Resource"
	rolloutRestartActionName            = "Rollout Restart"
	addCustomResourcePanelActionName    = "Add custom resource panel"
	deleteCustomResourcePanelActionName = "Delete custom resource panel"
	containerExecCommandActionName      = "Execute the command"
	changePodLogsContainerActionName    = "Change pod logs container"
)

var (
	// All common actions key map.
	keyMap = map[string][]interface{}{
		nextCyclicViewAction:                {gocui.KeyTab},
		backToPreviousViewAction:            {gocui.KeyEsc},
		toNavigationAction:                  {gocui.KeyEnter, gocui.KeyArrowRight, 'l'},
		navigationArrowLeft:                 {gocui.KeyArrowLeft, 'k'},
		navigationArrowRight:                {gocui.KeyArrowRight, 'l'},
		navigationDown:                      {gocui.KeyArrowDown, 'j', gocui.KeyTab},
		detailToNavigation:                  {gocui.KeyTab},
		detailArrowUp:                       {gocui.KeyArrowUp, 'h'},
		detailArrowDown:                     {gocui.KeyArrowDown, 'j'},
		previousLineAction:                  {gocui.KeyArrowUp, 'h'},
		nextLineAction:                      {gocui.KeyArrowDown, 'j'},
		previousPageAction:                  {gocui.KeyPgup},
		nextPageAction:                      {gocui.KeyPgdn},
		scrollUpAction:                      {gocui.MouseWheelUp},
		scrollDownAction:                    {gocui.MouseWheelDown},
		scrollTopAction:                     {gocui.KeyHome},
		scrollBottomAction:                  {gocui.KeyEnd},
		filterResourceActionName:            {gocui.KeyF4, 'f'},
		editResourceActionName:              {'e'},
		rolloutRestartActionName:            {'r'},
		moreActionsName:                     {gocui.KeyF3, 'm'},
		toFilteredViewAction:                {gocui.KeyTab, gocui.KeyArrowDown},
		toFilterInputAction:                 {gocui.KeyTab},
		filteredNextLineAction:              {gocui.KeyArrowDown},
		confirmFilterInputAction:            {gocui.KeyEnter},
		switchConfirmDialogOpt:              {gocui.KeyTab, gocui.KeyArrowRight, gocui.KeyArrowLeft, 'k', 'l'},
		confirmDialogEnter:                  {gocui.KeyEnter},
		addCustomResourcePanelActionName:    {'+'},
		deleteCustomResourcePanelActionName: {'-'},
		containerExecCommandActionName:      {'x'},
		optionsDialogEnter:                  {gocui.KeyEnter},
		inputDialogEnter:                    {gocui.KeyEnter},
		changePodLogsContainerActionName:    {'c'},
	}
)

func keyMapDescription(keys []interface{}, description string) string {
	keysName := make([]string, 0)
	for _, key := range keys {
		keysName = append(keysName, utils.GetKey(key))
	}
	return fmt.Sprintf("%-7s %s", strings.Join(keysName, " "), description)
}
