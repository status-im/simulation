package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/divan/graph-experiments/graph"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/status-im/simulator/simulation"
	"github.com/status-im/simulator/simulation/naivep2p"
	"github.com/status-im/simulator/simulation/whisperv6"
)

func main() {
	var (
		simType       = flag.String("type", "whisperv6", "Type of simulators (naivep2p, whisperv6)")
		ttl           = flag.Int("ttl", 10, "Message TTL for simulation")
		naiveP2PN     = flag.Int("naivep2p.N", 3, "Number of peers to propagate (0..N of peers)")
		naiveP2PDelay = flag.Duration("naivep2p.delay", 10*time.Millisecond, "Delay for each step")
		input         = flag.String("i", "network.json", "Input filename for pregenerated data to be used with simulation")
		output        = flag.String("o", "propagation.json", "Output filename for p2p sending data")
		gethlogLevel  = flag.String("loglevel", "crit", "Geth log level for whisper simulator (crti, error, warn, info, debug, trace)")
	)
	flag.Parse()

	data, err := graph.NewGraphFromJSON(*input)
	if err != nil {
		log.Fatal("Opening input file failed: ", err)
	}

	fd, err := os.Create(*output)
	if err != nil {
		log.Fatal("Opening output file failed: ", err)
	}
	defer fd.Close()

	var sim simulation.Simulator
	switch *simType {
	case "naivep2p":
		sim = naivep2p.NewSimulator(data, *naiveP2PN, *naiveP2PDelay)
	case "whisperv6":
		lvl, err := gethlog.LvlFromString(*gethlogLevel)
		if err != nil {
			lvl = gethlog.LvlCrit
		}
		gethlog.Root().SetHandler(gethlog.LvlFilterHandler(lvl, gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true))))
		sim = whisperv6.NewSimulator(data)
	default:
		log.Fatal("Unknown simulation type: ", *simType)
	}
	defer sim.Stop()

	// Start simulation by sending single message
	log.Printf("Starting message sending %s simulation for graph with %d nodes...", *simType, len(data.Nodes()))
	sendData := sim.SendMessage(0, *ttl)
	err = json.NewEncoder(fd).Encode(sendData)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Written %s propagation data into %s", *simType, *output)
}
