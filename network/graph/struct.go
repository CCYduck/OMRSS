package graph

import "src/network/topology"

type Graphs struct {
	TSNGraphs []*topology.Topology
	AVBGraphs []*topology.Topology
	Important_CANGraphs []*topology.Topology
	Unimportant_CANGraphs []*topology.Topology
}

func new_Graphs() *Graphs {
	return &Graphs{}
}
