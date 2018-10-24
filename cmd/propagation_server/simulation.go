package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/divan/graphx/graph"
	"github.com/status-im/simulation/propagation"
	"github.com/status-im/simulation/propagation/whisperv6"
)

// Simulation represents single simulation.
type Simulation struct {
	network *graph.Graph
	sim     propagation.Simulator
	plog    *propagation.Log
}

// NewSimulation creates Simulation for the given network.
func NewSimulation(network *graph.Graph) *Simulation {
	sim := whisperv6.NewSimulator(network)

	return &Simulation{
		network: network,
		sim:     sim,
	}
}

// Start starts simulation, creating network and preparing it for message sending.
func (s *Simulation) Start(sender, ttl, size int) {
	s.plog = s.sim.SendMessage(sender, ttl, size)
}

// Stop stops simulation and shuts down network.
func (s *Simulation) Stop() error {
	return s.sim.Stop()
}

// WriteOutput writes propagation log to the given io.Writer.
func (s *Simulation) WriteOutput(w io.Writer) error {
	return json.NewEncoder(w).Encode(s.plog)
}

// WriteOutputToFile writes propagation log to the given io.Writer.
func (s *Simulation) WriteOutputToFile(path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output file: %v", err)
	}
	defer fd.Close()

	return s.WriteOutput(fd)
}
