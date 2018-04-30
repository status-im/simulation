package naivep2p

import (
	"github.com/divan/graph-experiments/graph"
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
		if link.From == link.To {
			continue
		}
		if _, ok := ret[link.From]; !ok {
			ret[link.From] = make([]int, 0)
		}
		if _, ok := ret[link.To]; !ok {
			ret[link.To] = make([]int, 0)
		}

		peers := ret[link.From]
		peers = append(peers, link.To)
		ret[link.From] = peers

		peers = ret[link.To]
		peers = append(peers, link.From)
		ret[link.To] = peers
	}
	return ret
}
