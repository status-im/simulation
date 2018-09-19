package whisperv6

import (
	"fmt"
	"time"
)

// logEntry defines the reporting log entry for one
// p2p message sending.
type logEntry struct {
	From int
	To   int
	Ts   time.Duration
}

// String implements Stringer interface for logEntry.
func (l logEntry) String() string {
	return fmt.Sprintf("%s: %d -> %d", l.Ts.String(), l.From, l.To)
}

// newlogEntry creates new log entry.
func newlogEntry(start time.Time, from, to int) *logEntry {
	return &logEntry{
		Ts:   time.Since(start) / time.Millisecond,
		From: from,
		To:   to,
	}
}

/*
// newlogEntry creates new log entry.
func newlogEntry(start time.Time, from, to int) *logEntry {
	return &logEntry{
		Ts:   time.Since(start) / time.Millisecond,
		From: from,
		To:   to,
	}
}
*/
