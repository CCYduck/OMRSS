package path

import (
	"fmt"
	"src/network"
	"src/network/topology"
)
//主要的function

func BestPath(Network *network.Network) {
	
	for nth, flow := range Network.CANFlow_Set.ImportantCANFlows {
		fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination[0], Network.Graph_Set.Important_CANGraphs[nth])
		path := saveShortestPathsToBP(flow.Source, flow.Destination[0], Network.Graph_Set.Important_CANGraphs[nth])
		if path != nil {
			fmt.Println("Best Path:", path)
		} else {
			fmt.Println("No path found.")
		}
	}
}

func saveShortestPathsToBP(source int, target int, t *topology.Topology) []int {
	// Check if this path has already been taken
	graph := GetGarph(t)
	graph.ToVertex = target
	graph = Dijkstra(graph, target, source)
	//path := ReconstructPath(graph, source, target)
	//fmt.Println(graph.Path)
	//fmt.Println("Reconstructed path:", path)
	//return graph
	if len(graph.Path) > 0 {
		return graph.Path[0]
	}
	return nil
}

func ReconstructPath(graph *Graph, source, target int) []int {
	var revPath []int
	current := target
	for current != -1 {
		revPath = append(revPath, current)
		v := graph.FindVertex(current)
		if v == nil {
			break
		}
		current = v.Path
	}
	// 檢查是否正確到達 source
	if len(revPath) == 0 || revPath[len(revPath)-1] != source {
		return []int{}
	}
	// 反轉路徑，使其從 source 到 target
	for i, j := 0, len(revPath)-1; i < j; i, j = i+1, j-1 {
		revPath[i], revPath[j] = revPath[j], revPath[i]
	}
	return revPath
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

func (v2v *V2V) GetV2VEdge(terminal int) (*V2VEdge, bool) {
	for _, edge := range v2v.V2VEdges {
		if edge.FromVertex == terminal {
			return edge, false
		}
	}
	return &V2VEdge{FromVertex: terminal}, true
}

func (v2vedge *V2VEdge) InV2VEdge(terminal int) bool {
	for _, graph := range v2vedge.Graphs {
		if graph.ToVertex == terminal {
			return true
		}
	}
	return false
}

func (v2vedge *V2VEdge) GetV2VPath(terminal int) [][]int {
	var path [][]int
	for _, graph := range v2vedge.Graphs {
		if graph.ToVertex == terminal {
			path = graph.Path
		}
	}
	return path
}