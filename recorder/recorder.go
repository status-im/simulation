package simulations

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

// Recorder records send/receive events for
// generating further propagation log.
type Recorder struct {
	Log []*LogEntry

	start   time.Time
	nodeMap map[discover.NodeID]int
}

// NewRecorder inits new recorder.
func NewRecorder(nodeMap map[discover.NodeID]int) *Recorder {
	return &Recorder{
		start:   time.Now(),
		nodeMap: nodeMap,
	}
}

func (r *Recorder) Reset() {
	r.start = time.Now()
}

func (r *Recorder) Send(from, to discover.NodeID) {
	fromIdx, ok := r.nodeMap[from]
	if !ok {
		panic("node not found")
	}
	toIdx, ok := r.nodeMap[to]
	if !ok {
		panic("node not found")
	}
	log.Error("NewLogEntry", "start", r.start, "from", fromIdx, "to", toIdx)
	e := NewLogEntry(r.start, fromIdx, toIdx)
	r.Log = append(r.Log, e)
}

func (r *Recorder) Receive(from, to discover.NodeID) {
	log.Info("DidReceive", "from", from, "to", to)
}
