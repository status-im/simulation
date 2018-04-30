package stats

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"sync"
)

// Histogram implements simple histogram for node counters.
type Histogram struct {
	ranges []int
	counts []int

	totalCount int

	TotalDataPoint int
	MinDataPoint   int
	MaxDataPoint   int

	m sync.Mutex
}

// NewHistogram creates a new, ready to use Histogram.  The numBins
// must be >= 2. The binFirst is the width of the first bin.  The
// binGrowthFactor must be > 1.0 or 0.0.
//
// A special case of binGrowthFactor of 0.0 means the the allocated
// bins will have constant, non-growing size or "width".
func NewHistogram(
	numBins int,
	binFirst int,
	binGrowthFactor float64) *Histogram {
	gh := &Histogram{
		ranges:       make([]int, numBins),
		counts:       make([]int, numBins),
		totalCount:   0,
		MinDataPoint: math.MaxInt64,
		MaxDataPoint: 0,
	}

	gh.ranges[0] = 0
	gh.ranges[1] = binFirst

	for i := 2; i < len(gh.ranges); i++ {
		if binGrowthFactor == 0.0 {
			gh.ranges[i] = gh.ranges[i-1] + binFirst
		} else {
			gh.ranges[i] =
				int(math.Ceil(binGrowthFactor * float64(gh.ranges[i-1])))
		}
	}

	return gh
}

// Add increases the count in the bin for the given dataPoint.
func (gh *Histogram) Add(dataPoint int, count int) {
	gh.m.Lock()

	idx := search(gh.ranges, dataPoint)
	if idx >= 0 {
		gh.counts[idx] += count
		gh.totalCount += count

		gh.TotalDataPoint += dataPoint
		if gh.MinDataPoint > dataPoint {
			gh.MinDataPoint = dataPoint
		}
		if gh.MaxDataPoint < dataPoint {
			gh.MaxDataPoint = dataPoint
		}
	}

	gh.m.Unlock()
}

// Finds the last arr index where the arr entry <= dataPoint.
func search(arr []int, dataPoint int) int {
	i, j := 0, len(arr)

	for i < j {
		h := i + (j-i)/2 // Avoids h overflow, where i <= h < j.
		if dataPoint >= arr[h] {
			i = h + 1
		} else {
			j = h
		}
	}

	return i - 1
}

// AddAll adds all the counts from the src histogram into this
// histogram.  The src and this histogram must either have the same
// exact creation parameters.
func (gh *Histogram) AddAll(src *Histogram) {
	src.m.Lock()
	gh.m.Lock()

	for i := 0; i < len(src.counts); i++ {
		gh.counts[i] += src.counts[i]
	}
	gh.totalCount += src.totalCount

	gh.TotalDataPoint += src.TotalDataPoint
	if gh.MinDataPoint > src.MinDataPoint {
		gh.MinDataPoint = src.MinDataPoint
	}
	if gh.MaxDataPoint < src.MaxDataPoint {
		gh.MaxDataPoint = src.MaxDataPoint
	}

	gh.m.Unlock()
	src.m.Unlock()
}

// Graph emits an ascii graph to the optional out buffer, allocating a
// out buffer if none was supplied.  Returns the out buffer.  Each
// line emitted may have an optional prefix.
//
// For example:
//       0+  10=2 10.00% ********
//      10+  10=1 10.00% ****
//      20+  10=3 10.00% ************
func (gh *Histogram) EmitGraph(prefix []byte,
	out *bytes.Buffer) *bytes.Buffer {
	gh.m.Lock()

	ranges := gh.ranges
	rangesN := len(ranges)
	counts := gh.counts
	countsN := len(counts)

	if out == nil {
		out = bytes.NewBuffer(make([]byte, 0, 80*countsN))
	}

	var maxCount int
	for _, c := range counts {
		if maxCount < c {
			maxCount = c
		}
	}
	maxCountF := float64(maxCount)
	totCountF := float64(gh.totalCount)

	widthRange := len(strconv.Itoa(int(ranges[rangesN-1])))
	widthWidth := len(strconv.Itoa(int(ranges[rangesN-1] - ranges[rangesN-2])))
	widthCount := len(strconv.Itoa(int(maxCount)))

	// Each line looks like: "[prefix]START+WIDTH=COUNT PCT% BAR\n"
	f := fmt.Sprintf("%%%dd+%%%dd=%%%dd%% 7.2f%%%%",
		widthRange, widthWidth, widthCount)

	var runCount int // Running total while emitting lines.

	barLen := float64(len(bar))

	for i, c := range counts {
		if prefix != nil {
			out.Write(prefix)
		}

		var width int
		if i < countsN-1 {
			width = int(ranges[i+1] - ranges[i])
		}

		runCount += c
		fmt.Fprintf(out, f, ranges[i], width, c,
			100.0*(float64(runCount)/totCountF))

		if c > 0 {
			out.Write([]byte(" "))
			barWant := int(math.Floor(barLen * (float64(c) / maxCountF)))
			out.Write(bar[0:barWant])
		}

		out.Write([]byte("\n"))
	}

	gh.m.Unlock()

	return out
}

var bar = []byte("******************************")

// CallSync invokes the callback func while the histogram is locked.
func (gh *Histogram) CallSync(f func()) {
	gh.m.Lock()
	f()
	gh.m.Unlock()
}
