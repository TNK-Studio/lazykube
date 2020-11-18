package app

import (
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/jroimartin/gocui"
	"math"
)

func nextCyclicViewHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {

		currentView := gui.CurrentView()
		if currentView == nil {
			return nil
		}

		for index, viewName := range cyclicViews {
			if currentView.Name == viewName {
				nextIndex := index + 1
				if nextIndex >= len(cyclicViews) {
					nextIndex = 0
				}
				nextViewName := cyclicViews[nextIndex]
				log.Logger.Debugf("nextCyclicViewHandler - nextViewName: %s", nextViewName)
				return gui.FocusView(nextViewName, true)
			}
		}
		return nil
	}
}

func backToPreviousViewHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		if gui.HasPreviousView() {
			return gui.ReturnPreviousView()
		}

		return gui.FocusView(clusterInfoViewName, false)
	}
}

func toNavigationHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		return gui.FocusView(navigationViewName, true)
	}
}

func navigationArrowRightHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		options := viewNavigationMap[activeView.Name]
		if navigationIndex+1 >= len(options) {
			return nil
		}
		switchNavigation(navigationIndex + 1)
		gui.ReRenderViews(navigationViewName, detailViewName)
		return nil
	}
}

func navigationArrowLeftHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, view *gocui.View) error {
		if navigationIndex-1 < 0 {
			return gui.ReturnPreviousView()
		}
		switchNavigation(navigationIndex - 1)
		gui.ReRenderViews(navigationViewName, detailViewName)
		return nil
	}
}

func nextPageHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, oy := v.Origin()
		_, height := v.Size()
		newOy := int(math.Min(float64(len(v.ViewBufferLines())), float64(oy+height)))
		return v.SetOrigin(ox, newOy)
	}
}

func previousPageHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, oy := v.Origin()
		_, height := v.Size()
		newOy := int(math.Max(0, float64(oy-height)))
		return v.SetOrigin(ox, newOy)
	}
}

func scrollUpHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, oy := v.Origin()
		newOy := int(math.Max(0, float64(oy-2)))
		return v.SetOrigin(ox, newOy)
	}
}

func scrollDownHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, oy := v.Origin()

		reservedLines := 0
		_, sizeY := v.Size()
		reservedLines = sizeY

		totalLines := len(v.ViewBufferLines())
		if oy+reservedLines >= totalLines {
			v.Autoscroll = true
			return nil
		}

		return v.SetOrigin(ox, oy+2)
	}
}

func scrollTopHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = false
		ox, _ := v.Origin()
		return v.SetOrigin(ox, 0)
	}
}

func scrollBottomHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		totalLines := len(v.ViewBufferLines())
		if totalLines == 0 {
			return nil
		}
		_, vy := v.Size()
		if totalLines <= vy {
			return nil
		}

		ox, _ := v.Origin()
		v.Autoscroll = true
		return v.SetOrigin(ox, totalLines-1)
	}
}

func previousLineHandler(gui *guilib.Gui) func(gui *gocui.Gui, view *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		currentView := gui.CurrentView()
		if currentView == nil {
			return nil
		}

		_, height := v.Size()
		cx, cy := v.Cursor()
		ox, oy := v.Origin()

		if cy-1 <= 0 && oy-1 > 0 {
			v.SetOrigin(ox, int(math.Max(0, float64(oy-height+1))))
			v.SetCursor(cx, height-1)
			return nil
		}

		v.MoveCursor(0, -1, false)
		return nil
	}
}

func nextLineHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		currentView := gui.CurrentView()
		if currentView == nil {
			return nil
		}

		_, height := v.Size()
		cx, cy := v.Cursor()

		if cy+1 >= height-1 {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy+height-1)
			v.SetCursor(cx, 0)
			return nil
		}

		v.MoveCursor(0, 1, false)
		return nil
	}
}

func viewSelectedLineChangeHandler(gui *guilib.Gui, view *guilib.View, selectedLine string) error {
	if err := setViewSelectedLine(gui, view, selectedLine); err != nil {
		return err
	}
	gui.ReRenderViews(view.Name, navigationViewName, detailViewName)
	return nil
}

func editResourceHandler(gui *guilib.Gui) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}
}
