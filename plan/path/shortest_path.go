package path

import (
	"fmt"
	"src/network"
	"src/network/topology"
)
//主要的function

func BestPath(Network *network.Network) {
	for nth , flow := range Network.CANFlow_Set.ImportantCANFlows {
		fmt.Println(flow.Source, flow.Destinations, Network.Graph_Set.Important_CANGraphs[nth])
		saveShortestPathsToBP(flow.Source, flow.Destinations[0], Network.Graph_Set.Important_CANGraphs[nth])
		
	
	}
}



func saveShortestPathsToBP(vertex1 int, vertex2 int, t *topology.Topology) {
	// Check if this path has already been taken
	graph := GetGarph(t)
	graph.ToVertex = vertex1
	graph = Dijkstra(graph, vertex2, vertex1)
	fmt.Println(graph.Path)
}


func GetGarph(topology *topology.Topology) *Graph {
	graph := &Graph{}
	//Talker
	for _, t := range topology.Talker {
		gt := &Vertex{}
		gt.ID = t.ID
		gt.AddEdge(t.Connections)
		graph.Vertexs = append(graph.Vertexs, gt)
	}

	//Switch
	for _, s := range topology.Switch {
		gs := &Vertex{}
		gs.ID = s.ID
		gs.AddEdge(s.Connections)
		graph.Vertexs = append(graph.Vertexs, gs)
	}

	//Listener
	for _, l := range topology.Listener {
		gl := &Vertex{}
		gl.ID = l.ID
		gl.AddEdge(l.Connections)
		graph.Vertexs = append(graph.Vertexs, gl)
	}

	return graph
}


func (vertex *Vertex) AddEdge(connections []*topology.Connection) {
	for _, c := range connections {
		edge := &Edge{
			Strat: c.FromNodeID,
			End:   c.ToNodeID,
			Cost:  1,
		}
		vertex.Edges = append(vertex.Edges, edge)
	}
}