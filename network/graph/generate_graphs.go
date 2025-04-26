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
	for _, flow := range flows.Encapsulate {	
		for _,can2tsnflow := range flow.CAN2TSNFlows {

			if !graphs.checkListenerAndTalker(can2tsnflow.Source, can2tsnflow.Destination){
				t := topology.TopologyDeepCopy()                            // Duplicate of Topology		
				t.AddN2S2N_For_Path(can2tsnflow.Source, can2tsnflow.Destination, bytes_rate) // Undirected Graph
				graphs.CAN2TSNGraphs = append(graphs.CAN2TSNGraphs, t)
			}
			
		}
	}
	return graphs
}

func (graph *Graphs)checkListenerAndTalker(source int, destination int)bool{

	for _,g := range graph.CAN2TSNGraphs{
		if g.GetListenerAndTalker(source,destination){
			return true
		}	
	}
	return false
}

func (graph *Graphs)GetGarphBySD(source int, destination int) *topology.Topology{

	for _,g := range graph.CAN2TSNGraphs{
		if g.GetListenerAndTalker(source,destination){
			return g
		}	
	}
	return nil
}
