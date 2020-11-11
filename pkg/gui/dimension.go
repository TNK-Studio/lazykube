package gui

func BeneathView(aboveViewName string, heightFunc func(*Gui, *View) int, marginTopFunc func(*Gui, *View) int) func(gui *Gui, view *View) (int, int, int, int) {
	return func(gui *Gui, view *View) (int, int, int, int) {
		aboveX0, _, aboveX1, aboveY1, err := gui.g.ViewPosition(aboveViewName)
		if err != nil {
			return 0, 0, 0, 0
		}

		y0 := aboveY1 + marginTopFunc(gui, view)

		y1 := y0 + heightFunc(gui, view)
		if y1 < 0 {
			y1 = 0
		}

		return aboveX0, y0, aboveX1, y1
	}
}
