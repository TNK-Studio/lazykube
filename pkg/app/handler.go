package app

import (
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/pkg/errors"
	"math"
)

func nextCyclicViewHandler(gui *guilib.Gui, _ *guilib.View) error {
	currentView := gui.CurrentView()
	if currentView == nil {
		return nil
	}

	var nextViewName string
	for index, viewName := range cyclicViews {
		if currentView.Name == viewName {
			nextIndex := index + 1
			if nextIndex >= len(cyclicViews) {
				nextIndex = 0
			}
			nextViewName = cyclicViews[nextIndex]
			log.Logger.Debugf("nextCyclicViewHandler - nextViewName: %s", nextViewName)
			break
		}
	}
	if nextViewName == "" {
		return nil
	}
	gui.ReRenderViews(navigationViewName, detailViewName)
	return gui.FocusView(nextViewName, true)
}

func backToPreviousViewHandler(gui *guilib.Gui, _ *guilib.View) error {
	gui.ReRenderViews(navigationViewName, detailViewName)
	if gui.HasPreviousView() {
		return gui.ReturnPreviousView()
	}

	return gui.FocusView(clusterInfoViewName, false)
}

func toNavigationHandler(gui *guilib.Gui, _ *guilib.View) error {
	return gui.FocusView(navigationViewName, true)
}

func navigationArrowRightHandler(gui *guilib.Gui, _ *guilib.View) error {
	gui.ReRenderViews(navigationViewName, detailViewName)
	options := viewNavigationMap[activeView.Name]
	if navigationIndex+1 >= len(options) {
		return nil
	}
	switchNavigation(navigationIndex + 1)
	return nil
}

func navigationArrowLeftHandler(gui *guilib.Gui, _ *guilib.View) error {
	gui.ReRenderViews(navigationViewName, detailViewName)
	if navigationIndex-1 < 0 {
		return gui.ReturnPreviousView()
	}
	switchNavigation(navigationIndex - 1)
	return nil
}

func nextPageHandler(_ *guilib.Gui, view *guilib.View) error {
	view.Autoscroll = false
	ox, oy := view.Origin()
	_, height := view.Size()
	newOy := int(math.Min(float64(len(view.ViewBufferLines())), float64(oy+height)))
	return view.SetOrigin(ox, newOy)
}

func previousPageHandler(_ *guilib.Gui, view *guilib.View) error {
	view.Autoscroll = false
	ox, oy := view.Origin()
	_, height := view.Size()
	newOy := int(math.Max(0, float64(oy-height)))
	return view.SetOrigin(ox, newOy)
}

func scrollUpHandler(_ *guilib.Gui, view *guilib.View) error {
	view.Autoscroll = false
	ox, oy := view.Origin()
	newOy := int(math.Max(0, float64(oy-2)))
	return view.SetOrigin(ox, newOy)
}

func scrollDownHandler(_ *guilib.Gui, view *guilib.View) error {
	view.Autoscroll = false
	ox, oy := view.Origin()

	reservedLines := 0
	_, sizeY := view.Size()
	reservedLines = sizeY

	totalLines := len(view.ViewBufferLines())
	if oy+reservedLines >= totalLines {
		view.Autoscroll = true
		return nil
	}

	return view.SetOrigin(ox, oy+2)
}

func scrollTopHandler(_ *guilib.Gui, view *guilib.View) error {
	view.Autoscroll = false
	ox, _ := view.Origin()
	return view.SetOrigin(ox, 0)
}

func scrollBottomHandler(_ *guilib.Gui, view *guilib.View) error {
	totalLines := len(view.ViewBufferLines())
	if totalLines == 0 {
		return nil
	}
	_, vy := view.Size()
	if totalLines <= vy {
		return nil
	}

	ox, _ := view.Origin()
	view.Autoscroll = true
	return view.SetOrigin(ox, totalLines-1)
}

func previousLineHandler(gui *guilib.Gui, view *guilib.View) error {
	currentView := gui.CurrentView()
	if currentView == nil {
		return nil
	}

	_, height := view.Size()
	cx, cy := view.Cursor()
	ox, oy := view.Origin()

	if cy-1 <= 0 && oy-1 > 0 {
		err := view.SetOrigin(ox, int(math.Max(0, float64(oy-height+1))))
		if err != nil {
			return err
		}

		err = view.SetCursor(cx, height-1)
		if err != nil {
			return err
		}
		return nil
	}

	view.MoveCursor(0, -1, false)
	return nil
}

func nextLineHandler(gui *guilib.Gui, view *guilib.View) error {
	currentView := gui.CurrentView()
	if currentView == nil {
		return nil
	}

	_, height := view.Size()
	cx, cy := view.Cursor()

	if cy+1 >= height-1 {
		ox, oy := view.Origin()
		err := view.SetOrigin(ox, oy+height-1)
		if err != nil {
			return err
		}

		err = view.SetCursor(cx, 0)
		if err != nil {
			return err
		}

		return nil
	}

	view.MoveCursor(0, 1, false)
	return nil
}

func viewSelectedLineChangeHandler(gui *guilib.Gui, view *guilib.View, _ string) error {
	gui.ReRenderViews(view.Name, navigationViewName, detailViewName)
	gui.ClearViews(detailViewName)
	return nil
}

func getResourceNamespaceAndName(gui *guilib.Gui, resourceView *guilib.View) (string, string, error) {
	namespaceView, err := gui.GetView(namespaceViewName)
	if err != nil {
		return "", "", err
	}

	namespace := formatSelectedNamespace(namespaceView.SelectedLine)
	selected := resourceView.SelectedLine

	if selected == "" {
		return "", "", err
	}

	if !notResourceSelected(namespace) {
		resourceName := formatResourceName(selected, 0)
		if notResourceSelected(resourceName) {
			return "", "", err
		}
		return namespace, resourceName, nil
	}

	namespace = formatResourceName(selected, 0)
	resourceName := formatResourceName(selected, 1)
	if notResourceSelected(resourceName) {
		return "", "", err
	}

	if namespace == "" {
		namespace = kubecli.Cli.Namespace()
	}

	return namespace, resourceName, nil
}

func editResourceHandler(gui *guilib.Gui, view *guilib.View) error {
	view, resource, namespace, resourceName, err := resourceMoreActionHandlerHelper(gui, view)
	if errors.Is(err, resourceNotFoundErr) || errors.Is(err, noResourceSelectedErr) {
		// Todo: show error on panel
		return nil
	}

	cli(namespace).Edit(newStdStream(), resource, resourceName).Run()
	if err := gui.ForceFlush(); err != nil {
		return err
	}
	gui.ReRenderAll()
	return nil
}

func rolloutRestartHandler(gui *guilib.Gui, view *guilib.View) error {
	view, resource, namespace, resourceName, err := resourceMoreActionHandlerHelper(gui, view)
	if errors.Is(err, resourceNotFoundErr) || errors.Is(err, noResourceSelectedErr) {
		// Todo: show error on panel
		return nil
	}

	cli(namespace).RolloutRestart(viewStreams(view), resource, resourceName).Run()
	view.ReRender()
	return nil
}

func resourceMoreActionHandlerHelper(gui *guilib.Gui, view *guilib.View) (resourceView *guilib.View, resource string, namespace string, resourceName string, err error) {
	if view.Name == moreActionsViewName {
		var err error
		view, err = getMoreActionTriggerView(view)
		if err != nil {
			return nil, "", "", "", err
		}
	}
	resource = getViewResourceName(view.Name)
	if resource == "" {
		return nil, "", "", "", resourceNotFoundErr
	}
	namespace, resourceName, err = getResourceNamespaceAndName(gui, view)
	if err != nil {
		return nil, "", "", "", err
	}
	if notResourceSelected(resourceName) {
		return nil, "", "", "", noResourceSelectedErr
	}
	return view, resource, namespace, resourceName, nil
}

func newConfirmDialogHandler(title, relatedViewName string, handler guilib.ViewHandler) guilib.ViewHandler {
	return func(gui *guilib.Gui, view *guilib.View) error {
		if err := showConfirmActionDialog(gui, title, relatedViewName, handler); err != nil {
			return err
		}
		return nil
	}
}

func confirmDialogOptionHandler(gui *guilib.Gui, view *guilib.View, relatedViewName, option string, handler guilib.ViewHandler) error {
	if option == cancelDialogOpt {
		if err := gui.DeleteView(view.Name); err != nil {
			return err
		}
		if err := gui.ReturnPreviousView(); err != nil {
			return err
		}
		return nil
	}

	if option == confirmDialogOpt {
		relatedView, err := gui.GetView(relatedViewName)
		if err != nil {
			return err
		}
		if err := handler(gui, relatedView); err != nil {
			return err
		}
		if err := gui.DeleteView(view.Name); err != nil {
			return err
		}
		if err := gui.ReturnPreviousView(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func addCustomResourcePanelHandler(gui *guilib.Gui, view *guilib.View) error {
	if err := showFilterDialog(
		gui,
		"Filter resource by name.",
		func(string) error { return nil },
		func() ([]string, error) { return []string{}, nil },
		"Resource not found.",
	); err != nil {
		return err
	}
	return nil
}
