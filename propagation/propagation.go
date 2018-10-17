package propagation

// Simulator defines the simulators for message propagation within the graph.
type Simulator interface {
	SendMessage(idx, ttl int) *Log
	Stop() error
}
