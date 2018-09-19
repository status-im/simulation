package stats

import (
	"testing"

	"github.com/divan/graphx/graph"
	"github.com/status-im/simulation/propagation"
)

// node implements string-only graph.Node
type node string

func (n node) ID() string { return string(n) }

func TestAnalyze(t *testing.T) {
	g := graph.NewGraph()
	g.AddNode(node("0"))
	g.AddNode(node("1"))
	g.AddNode(node("2"))
	g.AddNode(node("3"))

	g.AddLink("0", "1")
	g.AddLink("1", "2")
	g.AddLink("2", "0")
	g.AddLink("0", "3")

	// example propagation log
	// three timestamps: 10, 20 and 30 ms
	// with first node hit 1 time, second and third - 3 times
	plog := &propagation.Log{
		Timestamps: []int{10, 20, 30},
		Nodes: [][]int{
			[]int{0, 1, 2},
			[]int{1, 2},
			[]int{2, 1, 3},
		},
	}

	stats := Analyze(g, plog)

	expected := []struct {
		name string
		hits int
	}{
		{"0", 1},
		{"1", 3},
		{"2", 3},
		{"3", 1},
	}

	for _, node := range expected {
		got := stats.NodeHits[node.name]
		if got != node.hits {
			t.Fatalf("Expected node '%s' to be hit %d times, but got %d", node.name, node.hits, got)
		}

	}
}
