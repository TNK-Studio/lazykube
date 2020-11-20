package gui

type ViewsZIndexSorter []*View

// Len Len
func (views ViewsZIndexSorter) Len() int { return len(views) }

// Less Less
func (views ViewsZIndexSorter) Less(i, j int) bool {
	if views[i].AlwaysOnTop && !views[j].AlwaysOnTop {
		return false
	}
	if !views[i].AlwaysOnTop && views[j].AlwaysOnTop {
		return true
	}

	return views[i].ZIndex < views[j].ZIndex
}

// Swap Swap
func (views ViewsZIndexSorter) Swap(i, j int) { views[i], views[j] = views[j], views[i] }
