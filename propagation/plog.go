package propagation

// Log describes propagation data collected during simulation.
type Log struct {
	Timestamps []int   // timestamps in milliseconds starting from T0
	Links      [][]int // indices of links for each step, len should be equal to len of Timestamps
	Nodes      [][]int // indices of nodes involved in each step, should match Timestamps
}

// NewLog inits a new empty plog structure with known number of timestamps. It
// allocates memory upfront.
func NewLog(n int) *Log {
	return &Log{
		Timestamps: make([]int, 0, n),
		Links:      make([][]int, 0, n),
		Nodes:      make([][]int, 0, n),
	}
}

// AddStep adds a single timestamp record to the propagation log.
// This is a preferred way of adding data to log, as it insures timestamp
// matching between all fields.
func (l *Log) AddStep(ts int, nodes, links []int) {
	l.Timestamps = append(l.Timestamps, ts)
	l.Nodes = append(l.Nodes, nodes)
	l.Links = append(l.Links, links)
}

// Less implements sort.Interface.
func (l *Log) Less(i, j int) bool {
	return l.Timestamps[i] < l.Timestamps[j]
}

// Swap implements sort.Interface.
func (l *Log) Swap(i, j int) {
	l.Timestamps[i], l.Timestamps[j] = l.Timestamps[j], l.Timestamps[i]
	l.Nodes[i], l.Nodes[j] = l.Nodes[j], l.Nodes[i]
	l.Links[i], l.Links[j] = l.Links[j], l.Links[i]
}

// Len implements sort.Interface.
func (l *Log) Len() int {
	return len(l.Timestamps)
}
