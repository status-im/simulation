package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/divan/graphx/formats"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/status-im/simulation/stats"
)

func main() {
	var (
		input        = flag.String("i", "network.json", "Input filename for pregenerated data to be used with simulation")
		output       = flag.String("o", "propagation.json", "Output filename for p2p sending data")
		gethlogLevel = flag.String("loglevel", "crit", "Geth log level for whisper simulator (crti, error, warn, info, debug, trace)")
		server       = flag.Bool("server", false, "Start as server to be used with whisperviz")
		serverAddr   = flag.String("h", "localhost:8084", "Address to bind to in server mode")
	)
	flag.Parse()

	setGethLogLevel(*gethlogLevel)

	if *server {
		log.Println("Starting simulator server on", *serverAddr)
		http.HandleFunc("/", simulationHandler)
		log.Fatal(http.ListenAndServe(*serverAddr, nil))
		return
	}

	data, err := formats.FromD3JSON(*input)
	if err != nil {
		log.Fatal("Opening input file failed: ", err)
	}

	sim := NewSimulation(data)
	log.Printf("Starting message sending simulation for graph with %d nodes...", len(data.Nodes()))
	sim.Start()
	defer sim.Stop()
	sim.WriteOutputToFile(*output)

	defer sim.Stop()

	// stats
	ss := stats.Analyze(data, sim.plog)
	ss.PrintVerbose()

	log.Printf("Written propagation data into %s", *output)
}

func setGethLogLevel(level string) {
	lvl, err := gethlog.LvlFromString(level)
	if err != nil {
		lvl = gethlog.LvlCrit
	}
	gethlog.Root().SetHandler(gethlog.LvlFilterHandler(lvl, gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true))))
}
