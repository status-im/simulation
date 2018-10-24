package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/divan/graphx/formats"
)

// SimulationRequests defines a POST request payload for simulation backend.
type SimulationRequest struct {
	SenderIdx int             `json:"senderIdx"` // index of the sender node (index of data.Nodes, in fact)
	TTL       int             `json:"ttl"`       // ttl in seconds
	MsgSize   int             `json:"msg_size"`  // msg size in bytes
	Network   json.RawMessage `json:"network"`   // current network graph
}

// simulationHandler serves request to start simulation. It expectes network graph
// in the request body, syncronously runs a new simulation on this network and
// sends back simulation log in JSON format.
//
// TODO(divan): in the future, simulation will probably take longer, so it'll have to upgrade
// connection to Websocket and talk to frontend this way.
func simulationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req SimulationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("[ERROR] Bad payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buf := bytes.NewBuffer(req.Network)
	network, err := formats.FromD3JSONReader(buf)
	if err != nil {
		log.Println("[ERROR] Bad payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Loaded graph with %d nodes", network.NumNodes())
	sim := NewSimulation(network)
	sim.Start(req.SenderIdx, req.TTL, req.MsgSize)
	defer sim.Stop()

	log.Println("Sending propagation log")
	sim.WriteOutput(w)
}
