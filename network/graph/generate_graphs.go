package graph

import (
	"src/network/flow"
	"src/network/topology"
)

func Generate_OSRO_Graphs(topology *topology.Topology, flows *flow.Flows, bytes_rate float64) *Graphs {
	// Constructing Graph structures
	graphs := new_Graphs()

	// Generating TSN Graphs
	for _, flow := range flows.TSNFlows {
		t := topology.TopologyDeepCopy()                               // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destination, bytes_rate) // Undirected Graph
		graphs.TSNGraphs = append(graphs.TSNGraphs, t)
	}

	// Generating AVB Graphs
	for _, flow := range flows.AVBFlows {
		t := topology.TopologyDeepCopy()                               // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destination, bytes_rate) // Undirected Graph
		graphs.AVBGraphs = append(graphs.AVBGraphs, t)
	}

	// Generating ImportantCAN Graphs
	for _, flow := range flows.ImportantCANFlows {
		t := topology.TopologyDeepCopy()                               // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destination, bytes_rate) // Undirected Graph
		graphs.Important_CANGraphs = append(graphs.Important_CANGraphs, t)

	}

	// Generating UnimportantCAN Graphs
	for _, flow := range flows.UnimportantCANFlows {
		t := topology.TopologyDeepCopy()                               // Duplicate of Topology
		t.AddN2S2N_For_Path(flow.Source, flow.Destination, bytes_rate) // Undirected Graph
		graphs.Unimportant_CANGraphs = append(graphs.Unimportant_CANGraphs, t)
	}

	return graphs
}
