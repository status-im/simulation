package stats

import (
	"fmt"
	"time"

	"github.com/status-im/simulation/propagation"
)

// Stats represents stats data for given simulation log.
type Stats struct {
	NodeHits     map[int]int
	NodeCoverage Coverage
	LinkCoverage Coverage
	Time         time.Duration
}

// PrintVerbose prints detailed terminal-friendly stats to
// the console.
func (s *Stats) PrintVerbose() {
	fmt.Println("Stats:")
	fmt.Println("Time elapsed:", s.Time)
	fmt.Println("Nodes coverage:", s.NodeCoverage)
	fmt.Println("Links coverage:", s.LinkCoverage)
}

// Analyze analyzes given propagation log and returns filled Stats object.
func Analyze(plog *propagation.Log, nodeCount, linkCount int) *Stats {
	t := analyzeTiming(plog)
	nodeHits := analyzeNodeHits(plog)
	nodeCoverage := analyzeNodeCoverage(nodeHits, nodeCount)
	linkCoverage := analyzeLinkCoverage(plog, linkCount)

	return &Stats{
		NodeHits:     nodeHits,
		NodeCoverage: nodeCoverage,
		LinkCoverage: linkCoverage,
		Time:         t,
	}
}

func analyzeNodeHits(plog *propagation.Log) map[int]int {
	nodeHits := make(map[int]int)

	for _, nodes := range plog.Nodes {
		for _, j := range nodes {
			nodeHits[j]++
		}
	}

	return nodeHits
}

func analyzeNodeCoverage(nodeHits map[int]int, total int) Coverage {
	actual := len(nodeHits)
	return NewCoverage(actual, total)
}

func analyzeLinkCoverage(plog *propagation.Log, total int) Coverage {
	linkHits := make(map[int]struct{})
	for _, links := range plog.Indices {
		for _, j := range links {
			linkHits[j] = struct{}{}
		}
	}

	actual := len(linkHits)
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
