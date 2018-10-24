package main

import (
	"log"
	"net/http"

	"github.com/divan/graphx/formats"
)

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

	data, err := formats.FromD3JSONReader(r.Body)
	if err != nil {
		log.Println("[ERROR] Bad payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Loaded graph with %d nodes", data.NumNodes())
	sim := NewSimulation(data)
	sim.Start()
	defer sim.Stop()

	log.Println("Sending propagation log")
	sim.WriteOutput(w)
}
