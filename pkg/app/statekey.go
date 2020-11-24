package app

const (
	viewLastRenderTimeStateKey    = "viewLastRenderTime" // value type: time.Time
	cpuPlotStateKey               = "cpuPlot"            // value type: *gui.Plot
	memoryPlotStateKey            = "memoryPlot"         // value type: *gui.Plot
	moreActionTriggerViewStateKey = "triggerView"        // value type: *gui.View
	filterInputValueStateKey      = "filterInputValue"   // value type: string
	confirmValueStateKey          = "confirmValue"       // value type: string
	logSinceTimeStateKey          = "logSinceTime"       // value type: time.Time
	podContainersStateKey         = "podContainers"      // value type: []string
	logContainerStateKey          = "logContainer"       // value type: string
)
