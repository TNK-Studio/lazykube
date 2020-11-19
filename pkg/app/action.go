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
		Name:    "nextCyclicView",
		Keys:    []interface{}{gocui.KeyTab},
		Handler: nextCyclicViewHandler,
		Mod:     gocui.ModNone,
	}

	backToPreviousView = &guilib.Action{
		Name:    "backToPreviousView",
		Key:     gocui.KeyEsc,
		Handler: backToPreviousViewHandler,
		Mod:     gocui.ModNone,
	}

	toNavigation = &guilib.Action{
		Name: "toNavigation",
		Keys: []interface{}{
			gocui.KeyEnter,
			gocui.KeyArrowRight,
		},
		Handler: toNavigationHandler,
		Mod:     gocui.ModNone,
	}

	previousLine = &guilib.Action{
		Name:    "previousLine",
		Key:     gocui.KeyArrowUp,
		Handler: previousLineHandler,
		Mod:     gocui.ModNone,
	}

	nextLine = &guilib.Action{
		Name:    "nextLine",
		Key:     gocui.KeyArrowDown,
		Handler: nextLineHandler,
		Mod:     gocui.ModNone,
	}

	appActions = []*guilib.Action{
		backToPreviousView,
		{
			Name: "previousPage",
			Keys: []interface{}{
				gocui.KeyPgup,
			},
			Handler: previousPageHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name: "nextPage",
			Keys: []interface{}{
				gocui.KeyPgdn,
			},
			Handler: nextPageHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name: "scrollUp",
			Keys: []interface{}{
				gocui.MouseWheelUp,
			},
			Handler: scrollUpHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name: "scrollDown",
			Keys: []interface{}{
				gocui.MouseWheelDown,
			},
			Handler: scrollDownHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name:    "scrollTop",
			Key:     gocui.KeyHome,
			Handler: scrollTopHandler,
			Mod:     gocui.ModNone,
		},
		{
			Name:    "scrollBottom",
			Key:     gocui.KeyEnd,
			Handler: scrollBottomHandler,
			Mod:     gocui.ModNone,
		},
	}

	filterAction = &guilib.Action{
		Name: "filterAction",
		Keys: []interface{}{
			gocui.KeyF4,
			'f',
		},
		Handler: func(gui *guilib.Gui, v *guilib.View) error {
			resourceName := getViewResourceName(v.Name)
			if resourceName == "" {
				return nil
			}
			if err := newFilterDialog(fmt.Sprintf("Input to filter %s", resourceName), gui, v.Name); err != nil {
				return err
			}
			return nil
		},
		Mod: gocui.ModNone,
	}

	editResourceAction = &guilib.Action{
		Name:    "Edit Resource",
		Key:     'e',
		Handler: editResourceHandler,
		Mod:     gocui.ModNone,
	}

	editResourceMoreAction = &moreAction{
		NeedSelectPanel:    false,
		NeedSelectResource: false,
		Action:             *editResourceAction,
	}

	moreActionsMap = map[string][]*moreAction{
		namespaceViewName: {
			editResourceMoreAction,
		},
		serviceViewName: {
			editResourceMoreAction,
		},
		deploymentViewName: {
			editResourceMoreAction,
		},
		podViewName: {
			editResourceMoreAction,
		},
	}
)

type (
	moreAction struct {
		NeedSelectPanel    bool
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
	detailView.SetOrigin(0, 0)
	gui.ReRenderViews(namespaceViewName, serviceViewName, deploymentViewName, podViewName, navigationViewName, detailViewName)
}

func newMoreActions(moreActions []*moreAction) *guilib.Action {
	return &guilib.Action{
		Name: "moreActions",
		Keys: []interface{}{
			gocui.KeyF3,
			'm',
		},
		Handler: func(gui *guilib.Gui, view *guilib.View) error {
			if err := newMoreActionDialog("More Actions", gui, view, moreActions); err != nil {
				return err
			}
			return nil
		},
		Mod: gocui.ModNone,
	}
}
