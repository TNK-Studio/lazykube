package gui

import (
	"fmt"
	"github.com/jesseduffield/asciigraph"
	"io"
	"time"
)

// Plot Plot
type Plot struct {
	Name  string
	data  []float64
	since time.Time

	DataGetter     func() []float64
	Height         func(plot *Plot) int
	Width          func(plot *Plot) int
	Max            func(plot *Plot) float64
	Min            func(plot *Plot) float64
	Caption        func(plot *Plot) string
	GraphFormatter func(graph string) string
}

// NewPlot NewPlot
func NewPlot(
	name string,
	dataGetter func() []float64,
	height func(plot *Plot) int,
	width func(plot *Plot) int,
	max func(plot *Plot) float64,
	min func(plot *Plot) float64,
	caption func(plot *Plot) string,
	graphFormatter func(string) string,
) *Plot {
	return &Plot{
		Name:           name,
		data:           make([]float64, 0),
		DataGetter:     dataGetter,
		Height:         height,
		Width:          width,
		Max:            max,
		Min:            min,
		Caption:        caption,
		GraphFormatter: graphFormatter,
		since:          time.Now(),
	}
}

// Graph Graph
func (plot *Plot) Graph() string {
	return plot.formatGraph(asciigraph.Plot(
		plot.data,
		asciigraph.Height(plot.Height(plot)),
		asciigraph.Width(plot.Width(plot)),
		asciigraph.Max(plot.Max(plot)),
		asciigraph.Min(plot.Min(plot)),
		asciigraph.Caption(plot.Caption(plot)),
	))
}

func (plot *Plot) formatGraph(graph string) string {
	if plot.GraphFormatter != nil {
		return plot.GraphFormatter(graph)
	}
	return graph
}

// Data Data
func (plot *Plot) Data() []float64 {
	return plot.data
}

// Since Since
func (plot *Plot) Since() time.Time {
	return plot.since
}

// Render Render
func (plot *Plot) Render(io io.Writer) {
	newData := plot.DataGetter()
	plot.data = append(plot.data, newData...)
	if len(plot.data) == 0 {
		fmt.Fprintf(io, "%s - No data. ", plot.Name)
		return
	}

	fmt.Fprint(io, plot.Graph())
}
