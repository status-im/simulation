package stats

import (
	"fmt"
	"log"
	"time"

	"github.com/divan/graph-experiments/graph"
	"github.com/status-im/simulator/propagation"
)

// Stats represents stats data for given simulation log.
type Stats struct {
	NodeHits     map[string]int
	NodeCoverage Coverage
	LinkCoverage Coverage
	Hist         *Histogram
	Time         time.Duration
}

// PrintVerbose prints detailed terminal-friendly stats to
// the console.
func (s *Stats) PrintVerbose() {
	fmt.Println("Stats:")
	fmt.Println("Time elapsed:", s.Time)
	fmt.Println("Nodes coverage:", s.NodeCoverage)
	fmt.Println("Links coverage:", s.LinkCoverage)
	fmt.Println("Histogram:")
	s.Hist.Print()
}

// Analyze analyzes given propagation log and returns filled Stats object.
func Analyze(g *graph.Graph, plog *propagation.Log) *Stats {
	t := analyzeTiming(plog)
	nodeHits, hist := analyzeNodeHits(g, plog)
	nodeCoverage := analyzeNodeCoverage(g, nodeHits)
	linkCoverage := analyzeLinkCoverage(g, plog)

	return &Stats{
		NodeHits:     nodeHits,
		NodeCoverage: nodeCoverage,
		LinkCoverage: linkCoverage,
		Hist:         hist,
		Time:         t,
	}
}

func analyzeNodeHits(g *graph.Graph, plog *propagation.Log) (map[string]int, *Histogram) {
	nodeHits := make(map[string]int)

	for _, nodes := range plog.Nodes {
		for _, j := range nodes {
			id, err := g.NodeIDByIdx(j)
			if err != nil {
				log.Fatal("Stats:", err)
			}
			nodeHits[id]++
		}
	}

	hist := NewHistogram(HistogramOptions{
		NumBuckets:   16,
		GrowthFactor: 0.2,
		MinValue:     1,
	})
	for key, v := range nodeHits {
		fmt.Println("Node:", key, v)
		err := hist.Add(int64(v))
		if err != nil {
			log.Println(err)
		}
	}

	return nodeHits, hist
}

func analyzeNodeCoverage(g *graph.Graph, nodeHits map[string]int) Coverage {
	actual := len(nodeHits)
	total := len(g.Nodes())
	return NewCoverage(actual, total)
}

func analyzeLinkCoverage(g *graph.Graph, plog *propagation.Log) Coverage {
	linkHits := make(map[int]struct{})
	for _, links := range plog.Indices {
		for _, j := range links {
			linkHits[j] = struct{}{}
		}
	}

	actual := len(linkHits)
	total := len(g.Links())
	return NewCoverage(actual, total)
}

// analyzeTiming returns the amount of time the simulation took.
func analyzeTiming(plog *propagation.Log) time.Duration {
	// log contains timestamps in milliseconds, so the
	// max value will be our number
	var max int
	for _, ts := range plog.Timestamps {
		if ts > max {
			max = ts
		}
	}
	return time.Duration(max) * time.Millisecond
}
