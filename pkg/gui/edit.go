package gui

import "github.com/jroimartin/gocui"

func NewViewEditor(gui *Gui, view *View) gocui.Editor {
	return gocui.EditorFunc(ViewEditorFunc(gui, view))
}

func ViewEditorFunc(gui *Gui, view *View) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		gocui.DefaultEditor.Edit(v, key, ch, mod)
		if view.OnEditedChange != nil {
			view.OnEditedChange(gui, view, key, ch, mod)
		}
	}
}
