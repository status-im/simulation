package stats

import (
	"fmt"
	"log"

	"v.io/x/ref/lib/stats/histogram"

	"github.com/divan/graph-experiments/graph"
	"github.com/status-im/simulator/simulation"
)

// Stats represents stats data for given simulation log.
type Stats struct {
	NodeHits     map[string]int
	NodeCoverage Coverage
	LinkCoverage Coverage
	Hist         *histogram.Histogram
}

// PrintVerbose prints detailed terminal-friendly stats to
// the console.
func (s *Stats) PrintVerbose() {
	fmt.Println("Stats:")
	fmt.Println("Nodes coverage:", s.NodeCoverage)
	fmt.Println("Links coverage:", s.LinkCoverage)
	fmt.Println("Node hits:", s.Hist.Value())
}

// Analyze analyzes given propagation log and returns filled Stats object.
func Analyze(g *graph.Graph, plog *simulation.Log) *Stats {
	nodeHits, hist := analyzeNodeHits(g, plog)
	nodeCoverage := analyzeNodeCoverage(g, nodeHits)
	linkCoverage := analyzeLinkCoverage(g, plog)

	return &Stats{
		NodeHits:     nodeHits,
		NodeCoverage: nodeCoverage,
		LinkCoverage: linkCoverage,
		Hist:         hist,
	}
}

func analyzeNodeHits(g *graph.Graph, plog *simulation.Log) (map[string]int, *histogram.Histogram) {
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

	hist := NewHistogram(1, 1, 1)
	for _, v := range nodeHits {
		hist.Add(v)
	}

	return nodeHits, hist
}

func analyzeNodeCoverage(g *graph.Graph, nodeHits map[string]int) Coverage {
	actual := len(nodeHits)
	total := len(g.Nodes())
	return NewCoverage(actual, total)
}

func analyzeLinkCoverage(g *graph.Graph, plog *simulation.Log) Coverage {
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
