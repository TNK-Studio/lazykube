package app

import (
	"fmt"
	guilib "github.com/TNK-Studio/lazykube/pkg/gui"
	"github.com/TNK-Studio/lazykube/pkg/kubecli"
	"github.com/TNK-Studio/lazykube/pkg/log"
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/gookit/color"
	v1 "k8s.io/api/core/v1"
	"math"
	"time"
)

const (
	cpuPlotStateKey        = "cpuPlot"
	memoryPlotStateKey     = "memoryPlot"
	plotLastRenderStateKey = "plotLastRender"
)

func podMetricsPlotRender(gui *guilib.Gui, view *guilib.View) error {
	view.ReRender()
	if !canRenderPlot(gui, view) {
		return nil
	}
	view.Clear()
	var err error

	podView, err := gui.GetView(podViewName)
	if err != nil {
		return err
	}
	resource := "pod"
	namespace, resourceName, err := getResourceNamespaceAndName(
		gui,
		podView,
	)
	if err != nil {
		return err
	}

	if notResourceSelected(resourceName) {
		showPleaseSelected(view, resource)
		return nil
	}

	metrics, err := kubecli.Cli.GetPodMetrics(namespace, resourceName, false, nil)
	if err != nil {
		log.Logger.Warningf("podMetricsDataGetter - kubecli.Cli.GetPodMetrics('%s', '%s', false, nil) error %s", namespace, resourceName, err)
	}
	fmt.Fprintln(view)
	cpuPlot := getPlot(
		gui,
		view,
		cpuPlotStateKey,
		"CPU: %0.0fm (%v)",
		namespace,
		resourceName,
		func() []float64 {
			data := make([]float64, 0)
			if metrics == nil {
				return data
			}
			for _, m := range metrics {
				data = append(data, float64(m[v1.ResourceCPU]))
			}
			return data
		},
		v1.ResourceCPU,
		color.Blue.Sprintf,
	)
	cpuPlot.Render(view)
	fmt.Fprintln(view)
	memoryPlot := getPlot(
		gui,
		view,
		memoryPlotStateKey,
		"Memory: %0.0fMi (%v)",
		namespace,
		resourceName,
		func() []float64 {
			data := make([]float64, 0)
			if metrics == nil {
				return data
			}
			for _, m := range metrics {
				data = append(data, float64(m[v1.ResourceMemory]))
			}
			return data
		},
		v1.ResourceMemory,
		color.Green.Sprintf,
	)
	memoryPlot.Render(view)
	return nil
}

func canRenderPlot(gui *guilib.Gui, view *guilib.View) bool {
	val, _ := view.State.Get(plotLastRenderStateKey)
	now := time.Now()
	if val == nil {
		view.State.Set(plotLastRenderStateKey, now)
		return true
	}

	since := val.(time.Time)
	if since.Add(1 * time.Second).Before(now) {
		view.State.Set(plotLastRenderStateKey, now)
		return true
	}

	return false
}

func getPlot(gui *guilib.Gui, view *guilib.View, plotStateKey, captionFormat, namespace, name string, dataGetter func() []float64, resourceName v1.ResourceName, colorSprintf func(format string, args ...interface{}) string) *guilib.Plot {
	var plot *guilib.Plot
	plotName := fmt.Sprintf("%s - %s", namespace, name)
	val, _ := view.State.Get(plotStateKey)
	newCPUPlot := false
	if val == nil {
		newCPUPlot = true
	} else {
		plot = val.(*guilib.Plot)
		if plot.Name != plotName {
			newCPUPlot = true
		}
	}

	if newCPUPlot {
		plot = guilib.NewPlot(
			plotName,
			dataGetter,
			podPlotHeight(gui, view),
			podPlotWidth(gui, view),
			podPlotMax(gui, view),
			podPlotMin(gui, view),
			podPlotCaption(captionFormat),
			func(graph string) string {
				return colorSprintf(graph)
			},
		)
		view.State.Set(plotStateKey, plot)
	}
	return plot
}

func podPlotHeight(gui *guilib.Gui, view *guilib.View) func(*guilib.Plot) int {
	return func(*guilib.Plot) int {
		_, MaxHeight := view.Size()

		height := MaxHeight/2 - 2
		if height <= 0 {
			return 0
		}

		return height
	}
}

func podPlotWidth(gui *guilib.Gui, view *guilib.View) func(*guilib.Plot) int {
	return func(*guilib.Plot) int {
		MaxWidth, _ := view.Size()
		return MaxWidth - 10
	}
}

func podPlotMax(gui *guilib.Gui, view *guilib.View) func(*guilib.Plot) float64 {
	return func(plot *guilib.Plot) float64 {
		return utils.MaxFloat64(plot.Data()) * 2
	}
}

func podPlotMin(gui *guilib.Gui, view *guilib.View) func(*guilib.Plot) float64 {
	return func(plot *guilib.Plot) float64 {
		return math.Min(0, utils.MinFloat64(plot.Data()))
	}
}

func podPlotCaption(format string) func(*guilib.Plot) string {
	return func(plot *guilib.Plot) string {
		length := len(plot.Data())
		if length == 0 {
			return "No data. "
		}
		return fmt.Sprintf(format, plot.Data()[length-1], time.Since(plot.Since().Round(time.Second)))
	}
}
