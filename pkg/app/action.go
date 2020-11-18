package app

import (
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
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

	editResource = &moreAction{
		keyName:     "e",
		description: "edit",
		action: &guilib.Action{
			Name:    "editResource",
			Key:     'e',
			Handler: nil,
			Mod:     0,
		},
	}

	moreActionsMap = map[string][]*moreAction{
		deploymentViewName: {
			editResource,
		},
	}
)

type moreAction struct {
	keyName     string
	description string
	action      *guilib.Action
}

func setViewSelectedLine(gui *guilib.Gui, view *guilib.View, selectedLine string) error {
	formatted := formatResourceName(selectedLine, 0)
	if notResourceSelected(formatted) {
		formatted = ""
	}
	return nil
}

func switchNamespace(gui *guilib.Gui, selectedNamespaceLine string) {
	kubecli.Cli.SetNamespace(selectedNamespaceLine)
	for _, viewName := range []string{serviceViewName, deploymentViewName, podViewName} {
		view, err := gui.GetView(viewName)
		if err != nil {
			return
		}
		view.SetOrigin(0, 0)
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
		Name: fmt.Sprintf("%sFilterAction", viewName),
		Keys: []interface{}{
			gocui.KeyF4,
			'f',
		},
		Handler: func(gui *guilib.Gui) func(g *gocui.Gui, v *gocui.View) error {
			return func(g *gocui.Gui, v *gocui.View) error {
				if err := newFilterDialog(fmt.Sprintf("Input to filter %s", resourceName), gui, viewName); err != nil {
					return err
				}
				return nil
			}
		},
		Mod: gocui.ModNone,
	}
}

func newMoreActions(viewName string, moreActions []*moreAction) *guilib.Action {
	return &guilib.Action{
		Name: fmt.Sprintf("%sMoreActions", viewName),
		Keys: []interface{}{
			gocui.KeyF5,
			'm',
		},
		Handler: func(gui *guilib.Gui) func(g *gocui.Gui, v *gocui.View) error {
			return func(g *gocui.Gui, v *gocui.View) error {
				if err := newMoreActionDialog("More Actions", gui, moreActions); err != nil {
					return err
				}
				return nil
			}
		},
		Mod: gocui.ModNone,
	}
}
