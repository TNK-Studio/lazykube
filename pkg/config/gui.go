package config

import "github.com/jroimartin/gocui"

// GuiConfig GuiConfig
type GuiConfig struct {
	Highlight  bool
	Cursor     bool
	FgColor    gocui.Attribute
	BgColor    gocui.Attribute
	SelBgColor gocui.Attribute
	SelFgColor gocui.Attribute
	Mouse      bool
	InputEsc   bool
}
