package whisperv6

import (
	"fmt"
	"time"
)

// LogEntry defines the reporting log entry for one
// p2p message sending.
type LogEntry struct {
	From int
	To   int
	Ts   time.Duration
}

// String implements Stringer interface for LogEntry.
func (l LogEntry) String() string {
	return fmt.Sprintf("%s: %d -> %d", l.Ts.String(), l.From, l.To)
}

// NewLogEntry creates new log entry.
func NewLogEntry(start time.Time, from, to int) *LogEntry {
	return &LogEntry{
		Ts:   time.Since(start) / time.Millisecond,
		From: from,
		To:   to,
	}
}
