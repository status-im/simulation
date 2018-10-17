package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"
)

func main() {
	var (
		gethlogLevel = flag.String("loglevel", "crit", "Geth log level for whisper simulator (crti, error, warn, info, debug, trace)")
		serverAddr   = flag.String("h", "localhost:8084", "Address to bind to in server mode")
	)
	flag.Parse()

	setGethLogLevel(*gethlogLevel)

	log.Println("Starting simulator server on", *serverAddr)
	http.HandleFunc("/", allowCORS(simulationHandler))
	log.Fatal(http.ListenAndServe(*serverAddr, nil))
}

func setGethLogLevel(level string) {
	lvl, err := gethlog.LvlFromString(level)
	if err != nil {
		lvl = gethlog.LvlCrit
	}
	gethlog.Root().SetHandler(gethlog.LvlFilterHandler(lvl, gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true))))
}
