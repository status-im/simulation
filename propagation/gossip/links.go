package gossip

import (
	"github.com/divan/graphx/graph"
)

// LinkIndex stores link information in form of indexes, rather than nodes IP.
type LinkIndex struct {
	From int
	To   int
}

// PrecalculatePeers creates map with peers indexes for faster lookup.
func PrecalculatePeers(data *graph.Graph) map[int][]int {
	links := data.Links()

	ret := make(map[int][]int)
	for _, link := range links {
		if link.From() == link.To() {
			continue
		}
		if _, ok := ret[link.FromIdx()]; !ok {
			ret[link.FromIdx()] = make([]int, 0)
		}
		if _, ok := ret[link.ToIdx()]; !ok {
			ret[link.ToIdx()] = make([]int, 0)
		}

		peers := ret[link.FromIdx()]
		peers = append(peers, link.ToIdx())
		ret[link.FromIdx()] = peers

		peers = ret[link.ToIdx()]
		peers = append(peers, link.FromIdx())
		ret[link.ToIdx()] = peers
	}
	return ret
}
