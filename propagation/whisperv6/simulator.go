package whisperv6

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/divan/graphx/graph"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/status-im/simulation/propagation"
)

// Simulator simulates WhisperV6 message propagation through the
// given p2p network. Implements Simulator interface.
type Simulator struct {
	data     *graph.Graph
	network  *simulations.Network
	whispers map[enode.ID]*whisper.Whisper
}

var ErrLinkExists = errors.New("link exists")

// NewSimulator intializes simulator for the given graph data.
// It uses defaults for PoW settings.
func NewSimulator(data *graph.Graph) *Simulator {
	rand.Seed(time.Now().UnixNano())

	cfg := &whisper.Config{
		MaxMessageSize:     whisper.DefaultMaxMessageSize,
		MinimumAcceptedPOW: 0.001,
	}

	whispers := make(map[enode.ID]*whisper.Whisper, len(data.Nodes()))
	services := map[string]adapters.ServiceFunc{
		"shh": func(ctx *adapters.ServiceContext) (node.Service, error) {
			return whispers[ctx.Config.ID], nil
		},
	}

	adapter := adapters.NewSimAdapter(services)
	network := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "shh",
	})

	nodes := data.Nodes()
	nodeCount := len(nodes)
	sim := &Simulator{
		data:    data,
		network: network,
	}

	log.Println("Creating nodes...")
	for i := 0; i < nodeCount; i++ {
		node, err := sim.network.NewNodeWithConfig(nodeConfig(i))
		if err != nil {
			log.Fatal("[ERROR] Can't start node: ", err)
		}
		// it's important to init whisper service here, as it
		// be initialized for each peer
		service := whisper.New(cfg)
		whispers[node.ID()] = service
	}

	log.Println("Starting nodes...")
	if err := network.StartAll(); err != nil {
		log.Fatal("[ERROR] Can't start nodes: ", err)
	}

	// subscribing to network events
	events := make(chan *simulations.Event)
	sub := sim.network.Events().Subscribe(events)
	defer sub.Unsubscribe()

	count := 0
	go func() {
		log.Println("Connecting nodes...")
		for _, link := range data.Links() {
			err := sim.connectNodes(link.FromIdx(), link.ToIdx())
			if err != nil && err != ErrLinkExists {
				log.Fatalf("[ERROR] Can't connect nodes %s and %s: %s", link.From(), link.To(), err)
			} else if err == nil {
				count++
			}
		}
	}()

	// wait for all nodes to establish connections
	var connected int
	var subErr error
	for connected < count && subErr == nil {
		select {
		case event := <-events:
			if event.Type == simulations.EventTypeConn {
				if event.Conn.Up {
					connected++
				}
			}
		case e := <-sub.Err():
			subErr = e
			log.Fatal("Failed to connect nodes", subErr)
		}
	}
	log.Println("All connections established")

	return sim
}

// Stop stops simulator and frees all resources if any.
func (s *Simulator) Stop() error {
	log.Println("Shutting down simulation nodes...")
	s.network.Shutdown()

	return nil
}

// SendMessage sends single message and tracks propagation. Implements propagation.Simulator.
func (s *Simulator) SendMessage(startNodeIdx, ttl int) *propagation.Log {
	node := s.network.Nodes[startNodeIdx]

	// the easiest way to send a message through the node is
	// by using its public RPC methods - ssh_post.
	client, err := node.Client()
	if err != nil {
		log.Fatal("Failed getting client", err)
	}

	log.Printf(" Sending Whisper message from %s...\n", node.ID().String())

	var symkeyID string
	symKey := make([]byte, aesKeyLength)
	rand.Read(symKey)

	err = client.Call(&symkeyID, "shh_addSymKey", hexutil.Bytes(symKey))
	if err != nil {
		log.Fatal("Failed adding new symmetric key: ", err)
	}

	// subscribing to network events
	events := make(chan *simulations.Event)
	sub := s.network.Events().Subscribe(events)
	defer sub.Unsubscribe()

	start := time.Now()

	msg := generateMessage(ttl, symkeyID)
	var ignored string
	err = client.Call(&ignored, "shh_post", msg)
	if err != nil {
		log.Fatal("Failed sending new post message: ", err)
	}

	// pre-cache node indexes
	var ncache = make(map[enode.ID]int)
	for i := range s.network.Nodes {
		ncache[s.network.Nodes[i].ID()] = i
	}

	timer := time.NewTimer(time.Duration(ttl) * time.Second)
	defer timer.Stop()
	var (
		subErr          error
		done, hasEvents bool
		plog            []*logEntry
	)
	for subErr == nil && !done {
		select {
		case event := <-events:
			if event.Type == simulations.EventTypeMsg {
				msg := event.Msg
				if msg.Code == 1 && msg.Protocol == "shh" && msg.Received == false {
					from := ncache[msg.One]
					to := ncache[msg.Other]
					entry := newlogEntry(start, from, to)
					plog = append(plog, entry)

					hasEvents = true
				}
			}
		case <-timer.C:
			done = true
		case e := <-sub.Err():
			subErr = e
		}
	}
	if subErr != nil {
		log.Fatal("[ERROR] Failed to collect propagation info", subErr)
	}
	if !hasEvents {
		log.Fatal("[ERROR] Didn't get any events, something wrong with simulator.")
	}

	return s.logEntries2PropagationLog(plog)
}

// logEntries2PropagationLog converts raw slice of LogEntries to PropagationLog,
// aggregating by timestamps and converting nodes indices to link indices.
// We expect that timestamps already bucketed into Nms groups.
func (s *Simulator) logEntries2PropagationLog(entries []*logEntry) *propagation.Log {
	tss := make(map[time.Duration][]int)
	tsnodes := make(map[time.Duration][]int)
	for _, entry := range entries {
		idx, err := s.data.LinkByIndices(entry.From, entry.To)
		if err != nil {
			log.Println("[EE] Wrong link", entry)
			continue
		}

		// fill links map
		if _, ok := tss[entry.Ts]; !ok {
			tss[entry.Ts] = make([]int, 0)
		}

		values := tss[entry.Ts]
		values = append(values, idx)
		tss[entry.Ts] = values

		// fill tsnodes map
		if _, ok := tsnodes[entry.Ts]; !ok {
			tsnodes[entry.Ts] = make([]int, 0)
		}
		nnodes := tsnodes[entry.Ts]
		nnodes = append(nnodes, entry.From, entry.To)
		tsnodes[entry.Ts] = nnodes
	}

	var ret = &propagation.Log{
		Timestamps: make([]int, 0, len(tss)),
		Indices:    make([][]int, 0, len(tss)),
		Nodes:      make([][]int, 0, len(tss)),
	}

	for ts, links := range tss {
		ret.Timestamps = append(ret.Timestamps, int(ts))
		ret.Indices = append(ret.Indices, links)
		ret.Nodes = append(ret.Nodes, tsnodes[ts])
	}

	return ret
}

// nodeConfig generates config for simulated node with random key.
func nodeConfig(idx int) *adapters.NodeConfig {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("[ERROR] Can't generate key: ", err)
	}
	id := pubkeyToID(&key.PublicKey)
	return &adapters.NodeConfig{
		ID:              id,
		PrivateKey:      key,
		Name:            nodeIdxToName(idx),
		EnableMsgEvents: true,
	}
}

func pubkeyToID(key *ecdsa.PublicKey) enode.ID {
	return enode.PubkeyToIDV4(key)
}

func nodeIdxToName(id int) string {
	return fmt.Sprintf("Node %d", id)
}

func (sim *Simulator) connectNodes(from, to int) error {
	// TODO(divan): check if we have IDs in from/to strings
	node1 := sim.network.Nodes[from]
	if node1 == nil {
		return fmt.Errorf("node with ID '%v' not found", from)
	}
	node2 := sim.network.Nodes[to]
	if node2 == nil {
		return fmt.Errorf("node with ID '%v' not found", to)
	}
	// if connection already exists, skip it, as network.Connect will fail
	if sim.network.GetConn(node1.ID(), node2.ID()) != nil {
		return ErrLinkExists
	}
	return sim.network.Connect(node1.ID(), node2.ID())
}
