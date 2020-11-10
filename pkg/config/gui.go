package config

import "github.com/jroimartin/gocui"

type GuiConfig struct {
	Highlight bool
	Cursor bool
	SelFgColor gocui.Attribute
}
