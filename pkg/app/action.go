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

	actions = []*guilib.Action{
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

	editResource = &guilib.Action{
		Name:    "Edit Resource",
		Key:     'e',
		Handler: editResourceHandler,
		Mod:     0,
	}

	moreActionsMap = map[string][]*guilib.Action{
		deploymentViewName: {
			editResource,
		},
	}
)

type (
	moreActionFunc func(viewName string) []*guilib.Action
)

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

func newFilterAction(viewName string, resourceName string) *guilib.Action {
	return &guilib.Action{
		Name: "filterAction",
		Keys: []interface{}{
			gocui.KeyF4,
			'f',
		},
		Handler: func(gui *guilib.Gui, v *guilib.View) error {
			if err := newFilterDialog(fmt.Sprintf("Input to filter %s", resourceName), gui, viewName); err != nil {
				return err
			}
			return nil
		},
		Mod: gocui.ModNone,
	}
}

func newMoreActions(viewName string, moreActions []*guilib.Action) *guilib.Action {
	return &guilib.Action{
		Name: fmt.Sprintf("%sMoreActions", viewName),
		Keys: []interface{}{
			gocui.KeyF5,
			'm',
		},
		Handler: func(gui *guilib.Gui, v *guilib.View) error {
			if err := newMoreActionDialog("More Actions", gui, moreActions); err != nil {
				return err
			}
			return nil
		},
		Mod: gocui.ModNone,
	}
}

func newEditResourceHandler(resource string) func(*guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(*guilib.Gui) func(*gocui.Gui, *gocui.View) error {
		return func(g *gocui.Gui, v *gocui.View) error {
			return nil
		}
	}
}
