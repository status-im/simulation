package stats

import "fmt"

// Coverage stores stats about coverage of large amount of units (nodes/links/etc).
type Coverage struct {
	Actual     int
	Total      int
	Percentage float64
}

// NewCoverage creates new Coverage out of the actual and total numbers.
func NewCoverage(actual, total int) Coverage {
	return Coverage{
		Actual:     actual,
		Total:      total,
		Percentage: 100.0 / float64(total) * float64(actual),
	}
}

// String implements Stringer interface for Coverage.
func (c Coverage) String() string {
	return fmt.Sprintf("%.0f%% (%d/%d)", c.Percentage, c.Actual, c.Total)
}
