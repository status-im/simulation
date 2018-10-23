package propagation

// Log describes propagation data collected during simulation.
type Log struct {
	Timestamps []int   // timestamps in milliseconds starting from T0
	Indices    [][]int // indices of links for each step, len should be equal to len of Timestamps
	Nodes      [][]int // indices of nodes involved in each step, should match Timestamps
}

// NewLog inits a new empty plog structure with known number of timestamps. It
// allocates memory upfront.
func NewLog(n int) *Log {
	return &Log{
		Timestamps: make([]int, 0, n),
		Indices:    make([][]int, 0, n),
		Nodes:      make([][]int, 0, n),
	}
}

// AddStep adds a single timestamp record to the propagation log.
// This is a preferred way of adding data to log, as it insures timestamp
// matching between all fields.
func (l *Log) AddStep(ts int, nodes, links []int) {
	l.Timestamps = append(l.Timestamps, ts)
	l.Nodes = append(l.Nodes, nodes)
	l.Indices = append(l.Indices, links)
}

func (l *Log) Less(i, j int) bool {
	return l.Timestamps[i] < l.Timestamps[j]
}
func (l *Log) Swap(i, j int) {
	l.Timestamps[i], l.Timestamps[j] = l.Timestamps[j], l.Timestamps[i]
	l.Nodes[i], l.Nodes[j] = l.Nodes[j], l.Nodes[i]
	l.Indices[j], l.Indices[j] = l.Indices[j], l.Indices[i]
}
func (l *Log) Len() int {
	return len(l.Timestamps)
}
