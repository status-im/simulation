package stats

import (
	"fmt"
	"sort"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
	"github.com/joliv/spark"
)

// Histogram is a simple histogram data holding structure.
type Histogram struct {
	data []float64
}

// NewHistogram creates and calculates a histogram from raw counts slice.
func NewHistogram(x []float64, nBins int) *Histogram {
	// x should be sorted
	sort.Slice(x, func(i, j int) bool { return x[i] < x[j] })

	dividers := []float64{}

	// automatically calculate dividers
	dividers = make([]float64, nBins+1)
	min := floats.Min(x)
	max := floats.Max(x)
	max += 1 // increase the max divider so max value of x is contained within the last bucket
	floats.Span(dividers, min, max)
	data := stat.Histogram(nil, dividers, x, nil)
	return &Histogram{
		data: data,
	}
}

// String implements Stringer for Histogram.
func (h *Histogram) String() string {
	return fmt.Sprintf("%v\n%v", h.data, spark.Line(h.data))
}
