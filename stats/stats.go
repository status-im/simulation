package stats

import (
	"fmt"
	"log"

	"github.com/divan/graph-experiments/graph"
	"github.com/status-im/simulator/simulation"
)

// Stats represents stats data for given simulation log.
type Stats struct {
	NodeHits map[string]int
}

// PrintVerbose prints detailed terminal-friendly stats to
// the console.
func (s *Stats) PrintVerbose() {
	fmt.Println("Stats:")
	for k, v := range s.NodeHits {
		fmt.Printf("%s: ", k)
		for i := 0; i < v; i++ {
			fmt.Printf(".")
		}
		fmt.Println()
	}
}

// Analyze analyzes given propagation log and returns filled Stats object.
func Analyze(g *graph.Graph, plog *simulation.Log) *Stats {
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

	return &Stats{
		NodeHits: nodeHits,
	}
}
