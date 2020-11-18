package gui

import "github.com/jroimartin/gocui"

// NewViewEditor NewViewEditor
func NewViewEditor(gui *Gui, view *View) gocui.Editor {
	return gocui.EditorFunc(ViewEditorFunc(gui, view))
}

// ViewEditorFunc ViewEditorFunc
func ViewEditorFunc(gui *Gui, view *View) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		gocui.DefaultEditor.Edit(v, key, ch, mod)
		if view.OnEditedChange != nil {
			view.OnEditedChange(gui, view, key, ch, mod)
		}
	}
}
