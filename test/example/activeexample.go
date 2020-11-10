package main

import (
	"fmt"
	"github.com/TNK-Studio/lazykube/pkg/config"
	"github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/jroimartin/gocui"
)

var (
	active = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func main() {
	c1 := &gui.View{
		Name:     "c1",
		Title:    "View 1 (editable)",
		Editable: true,
		Wrap:     true,
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			maxWidth, _ := gui.Size()
			return maxWidth/2 - 1
		},
		LowerRightPointYFunc: func(gui *gui.Gui, view *gui.View) int {
			_, maxHeight := gui.Size()
			return maxHeight/2 - 1
		},
		Render: func(gui *gui.Gui, view *gui.View) error {
			if _, err := gui.SetCurrentView(view.Name); err != nil {
				return err
			}

			if _, err := gui.SetViewOnTop(view.Name); err != nil {
				return err
			}
			return nil
		},
	}

	c2 := &gui.View{
		Name:       "c2",
		Title:      "View 2",
		Wrap:       true,
		Autoscroll: true,
		UpperLeftPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			maxWidth, _ := gui.Size()
			return maxWidth/2 - 1
		},
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			maxWidth, _ := gui.Size()
			return maxWidth - 1
		},
		LowerRightPointYFunc: func(gui *gui.Gui, view *gui.View) int {
			_, maxHeight := gui.Size()
			return maxHeight/2 - 1
		},
	}

	c3 := &gui.View{
		Name:       "c3",
		Title:      "View 3",
		Wrap:       true,
		Autoscroll: true,
		UpperLeftPointYFunc: func(gui *gui.Gui, view *gui.View) int {
			_, maxHeight := gui.Size()
			return maxHeight/2 - 1
		},
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			maxWidth, _ := gui.Size()
			return maxWidth/2 - 1
		},
		LowerRightPointYFunc: func(gui *gui.Gui, view *gui.View) int {
			_, maxHeight := gui.Size()
			return maxHeight - 1
		},
	}

	c4 := &gui.View{
		Name:       "c4",
		Title:      "View 4 (editable)",
		Editable:   true,
		Wrap:       true,
		Autoscroll: true,
		UpperLeftPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			maxWidth, _ := gui.Size()
			return maxWidth / 2
		},
		UpperLeftPointYFunc: func(gui *gui.Gui, view *gui.View) int {
			_, maxHeight := gui.Size()
			return maxHeight / 2
		},
		LowerRightPointXFunc: func(gui *gui.Gui, view *gui.View) int {
			maxWidth, _ := gui.Size()
			return maxWidth - 1
		},
		LowerRightPointYFunc: func(gui *gui.Gui, view *gui.View) int {
			_, maxHeight := gui.Size()
			return maxHeight - 1
		},
	}

	nextView := &gui.Action{
		Name: "NextView",
		Key:  gocui.KeyTab,
		Handler: func(gui *gui.Gui) func(*gocui.Gui, *gocui.View) error {
			return func(g *gocui.Gui, v *gocui.View) error {
				nextIndex := (active + 1) % len(gui.GetViews())
				name := gui.GetViews()[nextIndex].Name

				out, err := g.View(c2.Name)
				if err != nil {
					return err
				}

				if v == nil {
					return nil
				}

				fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

				if _, err := g.SetCurrentView(name); err != nil {
					return err
				}

				if _, err := g.SetViewOnTop(name); err != nil {
					return err
				}

				if nextIndex == 0 || nextIndex == 3 {
					g.Cursor = true
				} else {
					g.Cursor = false
				}

				active = nextIndex
				return nil
			}
		},
		Mod: gocui.ModNone,
	}

	conf := config.GuiConfig{
		Highlight:  true,
		Cursor:     true,
		SelFgColor: gocui.ColorGreen,
	}

	g := gui.NewGui(
		conf,
		c1, c2, c3, c4,
	)
	defer g.Close()
	g.BindAction("", nextView)
	g.Run()
}
