# Propagation

## Plog
---
PLog describes messages propagation data in the concise and streamable form, to be used as a main data format for representing simulated message propagation data.

### Format

Currently a single plog record holds information about nodes and edges active per timestamp. Both nodes and edges are identified by their indicies (this should match the original graph indicies).

For example, if network graph has 3 nodes, recorded indicies should be 0,1 and 2 respectively.

Current structure:

```go
package plog

type Data struct {
    Timestamps []int   // timestamps in milliseconds starting from T0
    Indices    [][]int // indices of links for each step, len should match Timestamps 
    Nodes      [][]int // indices of nodes involved in each step, len should match Timestamps 

```

So, if you have nodes 0 and 1 activated on first timestamp step, and nodes 1 and 2 on second, it'll be:

plog.Timestamps = []int{10, 20} // say, 10 and 20 ms timestamps
plog.Nodes = [][]int{[]int{0, 1}, []int{1, 2}} 
