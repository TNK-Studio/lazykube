package app

import (
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
)

var (
	cyclicViews = []string{clusterInfoViewName, namespaceViewName, serviceViewName, deploymentViewName, podViewName}

	nextCyclicView = &guilib.Action{
		Name:    nextCyclicViewAction,
		Keys:    keyMap[nextCyclicViewAction],
		Handler: nextCyclicViewHandler,
		Mod:     gocui.ModNone,
	}

	backToPreviousView = &guilib.Action{
		Name:    backToPreviousViewAction,
		Keys:    keyMap[backToPreviousViewAction],
		Handler: backToPreviousViewHandler,
		Mod:     gocui.ModNone,
	}

	toNavigation = &guilib.Action{
		Name:    toNavigationAction,
		Keys:    keyMap[toNavigationAction],
		Handler: toNavigationHandler,
		Mod:     gocui.ModNone,
	}

	previousLine = &guilib.Action{
		Name:    previousLineAction,
		Keys:    keyMap[previousLineAction],
		Handler: previousLineHandler,
		Mod:     gocui.ModNone,
	}

	nextLine = &guilib.Action{
		Name:    nextLineAction,
		Keys:    keyMap[nextLineAction],
		Handler: nextLineHandler,
		Mod:     gocui.ModNone,
	}

	filterResource = &guilib.Action{
		Name: filterResourceActionName,
		Keys: keyMap[filterResourceActionName],
		Handler: func(gui *guilib.Gui, v *guilib.View) error {
			resourceName := getViewResourceName(v.Name)
			if resourceName == "" {
				return nil
			}
			resourceViewName := resourceViewName(resourceName)
			if err := showFilterDialog(
				gui,
				fmt.Sprintf("Input to filter %s", resourceName),
				func(filtered string) error {
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
				func() ([]string, error) {
					var data []string
					resourceView, err := gui.GetView(resourceViewName)
					if err != nil {
						return nil, err
					}

					data = resourceView.ViewBufferLines()
					return data, nil
				},
				filteredNoResource,
			); err != nil {
				return err
			}
			return nil
		},
		Mod: gocui.ModNone,
	}

	editResourceAction = &guilib.Action{
		Name:    editResourceActionName,
		Keys:    keyMap[editResourceActionName],
		Handler: editResourceHandler,
		Mod:     gocui.ModNone,
	}

	rolloutRestartAction = &guilib.Action{
		Keys:    keyMap[rolloutRestartActionName],
		Name:    rolloutRestartActionName,
		Handler: rolloutRestartHandler,
		Mod:     gocui.ModNone,
	}

	editResourceMoreAction = &moreAction{
		NeedSelectResource: true,
		Action:             *editResourceAction,
	}

	addCustomResourcePanelAction = &guilib.Action{
		Keys:    keyMap[addCustomResourcePanelActionName],
		Name:    addCustomResourcePanelActionName,
		Handler: addCustomResourcePanelHandler,
		Mod:     gocui.ModNone,
	}

	deleteCustomResourcePanelAction = &guilib.Action{
		Keys:    keyMap[deleteCustomResourcePanelActionName],
		Name:    deleteCustomResourcePanelActionName,
		Handler: deleteCustomResourcePanelHandler,
		Mod:     gocui.ModNone,
	}

	containerExecCommandAction = &guilib.Action{
		Keys:    keyMap[containerExecCommandActionName],
		Name:    containerExecCommandActionName,
		Handler: nil,
		Mod:     gocui.ModNone,
	}

	addCustomResourcePanelMoreAction = &moreAction{
		NeedSelectResource: false,
		Action:             *addCustomResourcePanelAction,
	}
	deleteCustomResourcePanelMoreAction = &moreAction{
		NeedSelectResource: false,
		Action:             *deleteCustomResourcePanelAction,
	}

	commonResourceMoreActions = []*moreAction{
		addCustomResourcePanelMoreAction,
		editResourceMoreAction,
	}

	moreActionsMap = map[string][]*moreAction{
		clusterInfoViewName: {
			addCustomResourcePanelMoreAction,
		},
		namespaceViewName: append(
			commonResourceMoreActions,
		),
		serviceViewName: append(
			commonResourceMoreActions,
		),
		deploymentViewName: append(
			commonResourceMoreActions,
			&moreAction{
				NeedSelectResource: true,
				Action:             *newConfirmDialogAction(deploymentViewName, rolloutRestartAction),
			},
		),
		podViewName: append(
			commonResourceMoreActions,
		),
		navigationViewName: {
			addCustomResourcePanelMoreAction,
		},
		detailViewName: {
			addCustomResourcePanelMoreAction,
		},
	}

	appActions = []*guilib.Action{
		backToPreviousView,
		{
			Name:    previousPageAction,
			Keys:    keyMap[previousPageAction],
			Handler: previousPageHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name:    nextPageAction,
			Keys:    keyMap[nextPageAction],
			Handler: nextPageHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name:    scrollUpAction,
			Keys:    keyMap[scrollUpAction],
			Handler: scrollUpHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name:    scrollDownAction,
			Keys:    keyMap[scrollDownAction],
			Handler: scrollDownHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name:    scrollTopAction,
			Keys:    keyMap[scrollTopAction],
			Handler: scrollTopHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name:    scrollBottomAction,
			Keys:    keyMap[scrollBottomAction],
			Handler: scrollBottomHandler,
			Mod:     gocui.ModNone,
		},
	}
)

type (
	moreAction struct {
		NeedSelectResource bool
		guilib.Action
	}
)

func toMoreActionArr(actions []*moreAction) []guilib.ActionInterface {
	arr := make([]guilib.ActionInterface, 0)
	for _, act := range actions {
		arr = append(arr, act)
	}
	return arr
}

func switchNamespace(gui *guilib.Gui, selectedNamespaceLine string) {
	kubecli.Cli.SetNamespace(selectedNamespaceLine)
	for _, viewName := range []string{serviceViewName, deploymentViewName, podViewName} {
		view, err := gui.GetView(viewName)
		if err != nil {
			return
		}
		err = view.SetOrigin(0, 0)
		if err != nil {
			log.Logger.Warningf("switchNamespace - error %s", err)
		}
	}

	detailView, err := gui.GetView(detailViewName)
	if err != nil {
		return
	}
	detailView.Autoscroll = false
	err = detailView.SetOrigin(0, 0)
	if err != nil {
		log.Logger.Warningf("switchNamespace - detailView.SetOrigin(0, 0) error %s", err)
	}
	gui.ReRenderViews(resizeableViews...)
	gui.ReRenderViews(navigationViewName, detailViewName)
}

func newMoreActions(moreActions []*moreAction) *guilib.Action {
	return &guilib.Action{
		Name: moreActionsName,
		Keys: keyMap[moreActionsName],
		Handler: func(gui *guilib.Gui, view *guilib.View) error {
			if err := showMoreActionDialog(gui, view, "More Actions", moreActions); err != nil {
				return err
			}
			return nil
		},
		Mod: gocui.ModNone,
	}
}

func newConfirmDialogAction(relatedViewName string, action *guilib.Action) *guilib.Action {
	confirmTitle := fmt.Sprintf("Confirm to '%s' ?", action.Name)
	return &guilib.Action{
		Keys:            action.Keys,
		Name:            action.Name,
		Key:             action.Key,
		Handler:         newConfirmDialogHandler(confirmTitle, relatedViewName, action.Handler),
		ReRenderAllView: action.ReRenderAllView,
		Mod:             action.Mod,
	}
}

func newConfirmFilterInput(confirmHandler func(string) error) *guilib.Action {
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

			return confirmHandler(filtered)
		},
		Mod: gocui.ModNone,
	}
	return confirmFilterInput
}
