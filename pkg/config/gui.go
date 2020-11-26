package config

import "github.com/jroimartin/gocui"

// GuiConfig GuiConfig
type GuiConfig struct {
	Highlight  bool            `yaml:"highlight"`
	Cursor     bool            `yaml:"cursor"`
	FgColor    gocui.Attribute `yaml:"fg_color"`
	BgColor    gocui.Attribute `yaml:"bg_color"`
	SelBgColor gocui.Attribute `yaml:"sel_bg_color"`
	SelFgColor gocui.Attribute `yaml:"sel_fg_color"`
	Mouse      bool            `yaml:"mouse"`
	InputEsc   bool            `yaml:"input_esc"`
}
