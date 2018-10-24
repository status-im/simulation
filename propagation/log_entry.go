package propagation

import (
	"fmt"
	"log"
	"time"

	"github.com/divan/graphx/graph"
)

// LogEntry defines the reporting log entry for one
// p2p message sending.
type LogEntry struct {
	From int
	To   int
	Ts   int64
}

// String implements Stringer interface for LogEntry.
func (l LogEntry) String() string {
	return fmt.Sprintf("%d: %d -> %d", l.Ts, l.From, l.To)
}

// NewlogEntry creates new log entry.
func NewLogEntry(t, start time.Time, from, to int) *LogEntry {
	delta := t.Sub(start)
	return &LogEntry{
		Ts:   int64(delta / time.Millisecond),
		From: from,
		To:   to,
	}
}

// LogEntries2Log converts raw slice of LogEntries to Log,
// aggregating by timestamps and converting nodes indices to link indices.
// We expect that timestamps already bucketed into Nms groups.
func LogEntries2Log(data *graph.Graph, entries []*LogEntry) *Log {
	tss := make(map[int64][]int)
	tsnodes := make(map[int64][]int)
	for _, entry := range entries {
		idx, err := data.LinkByIndices(entry.From, entry.To)
		if err != nil {
			log.Println("[EE] Wrong link", entry)
			continue
		}

		// fill links map
		if _, ok := tss[entry.Ts]; !ok {
			tss[entry.Ts] = make([]int, 0)
		}

		values := tss[entry.Ts]
		values = append(values, idx)
		tss[entry.Ts] = values

		// fill tsnodes map
		if _, ok := tsnodes[entry.Ts]; !ok {
			tsnodes[entry.Ts] = make([]int, 0)
		}
		nnodes := tsnodes[entry.Ts]
		nnodes = append(nnodes, entry.From, entry.To)
		tsnodes[entry.Ts] = nnodes
	}

	plog := NewLog(len(tss))
	for ts, links := range tss {
		plog.AddStep(int(ts), tsnodes[ts], links)
	}

	return plog
}
