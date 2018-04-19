package stats

import "testing"

func TestCoverage(t *testing.T) {
	var tests = []struct {
		name          string
		actual, total int
		expected      Coverage
	}{
		{"0%", 0, 100, Coverage{0, 100, 0.0}},
		{"50%", 50, 100, Coverage{50, 100, 50.0}},
		{"100%", 100, 100, Coverage{100, 100, 100.0}},
	}

	for _, test := range tests {
		c := NewCoverage(test.actual, test.total)
		if c != test.expected {
			t.Fatalf("Expected test '%s' to get %v, but got %v", test.name, test.expected, c)
		}
	}
}
