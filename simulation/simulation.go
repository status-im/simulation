package simulation

// Simulator defines the simulators for message propagation within the graph.
type Simulator interface {
	SendMessage(idx, ttl int) *Log
	Stop() error
}

// Log represnts log of p2p message propagation
// with relative timestamps (starting from T0).
type Log struct {
	Timestamps []int   // timestamps in milliseconds starting from T0
	Indices    [][]int // indices of links for each step, len should be equal to len of Timestamps field
	Nodes      [][]int // indices of nodes involved in each step
}
