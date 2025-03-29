package path

import (
	"fmt"
	"src/network"
	"src/network/topology"
)
//主要的function

func BestPath(Network *network.Network) {
	for nth, flow := range Network.TSNFlow_Set.TSNFlows {
		// fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination, Network.Graph_Set.TSNGraphs[nth])
		//fmt.Printf("Flow: Source=%d, Destination=%v", flow.Source, flow.Destination)
		path,kpath := saveShortestPathsToGraph(flow.Source, flow.Destination, Network.TSNGraph_Set.TSNGraphs[nth])
		if path != nil {
			fmt.Println("Best Path:", path)
			fmt.Printf("KPath: Source=%d, Target=%d, NodeCount=%d\n",
			kpath.Source, kpath.Target, len(kpath.Paths[0].Nodes))
		} else {
			fmt.Println("No path found.")
		}
		p := ConvertIDsToPath(path,Network.TSNGraph_Set.TSNGraphs[nth])
		// 建立一個 *KPath
		k := new_KPath(nth, flow.Source, flow.Destination)
		// 把剛生成的 p 加到 k.Paths 裡
		k.Paths = append(k.Paths, p)

		// 你可以存到一個 KPath_Set 或印出來
		// fmt.Printf("KPath: Source=%d, Target=%d, NodeIDs=%v\n", k.Source, k.Target, path)

	}

	// for nth, flow := range Network.TSNFlow_Set.AVBFlows {
	// 	// fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination[0], Network.Graph_Set.AVBGraphs[nth])
	// 	//fmt.Printf("Flow: Source=%d, Destination=%v", flow.Source, flow.Destination)
	// 	path := saveShortestPathsToGraph(flow.Source, flow.Destination, Network.TSNGraph_Set.AVBGraphs[nth])
	// 	if path != nil {
	// 		fmt.Println("Best Path:", path)
	// 	} else {
	// 		fmt.Println("No path found.")
	// 	}
	// }
	
	// for nth, flow := range Network.CANFlow_Set.ImportantCANFlows {
	// 	//fmt.Printf("Flow: Source=%d, Destination=%v ", flow.Source, flow.Destination)
	// 	// fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination[0], Network.Graph_Set.Important_CANGraphs[nth])
	// 	path := saveShortestPathsToGraph(flow.Source, flow.Destination, Network.CANGraph_Set.Important_CANGraphs[nth])
	// 	if path != nil {
	// 		fmt.Println("Best Path:", path)
	// 	} else {
	// 		fmt.Println("No path found.")
	// 	}
	// }

	// for nth, flow := range Network.CANFlow_Set.UnimportantCANFlows {
	// 	//fmt.Printf("Flow: Source=%d, Destination=%v ", flow.Source, flow.Destination)
	// 	// fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination[0], Network.Graph_Set.Unimportant_CANGraphs[nth])
	// 	path := saveShortestPathsToGraph(flow.Source, flow.Destination, Network.CANGraph_Set.Unimportant_CANGraphs[nth])
	// 	if path != nil {
	// 		fmt.Println("Best Path:", path)
	// 	} else {
	// 		fmt.Println("No path found.")
	// 	}
	// }
	
}

func saveShortestPathsToGraph(source int, target int, t *topology.Topology) ([]int, *KPath) {
	// Check if this path has already been taken
	graph := GetGarph(t)
	graph.ToVertex = target
	graph = Dijkstra(graph, target, source)
	
	if len(graph.Path) > 0 {
		path := &Path{
			Weight: 0,
		}
		for count,id :=range graph.Path[0] {
			fmt.Printf("Source=%d,NodeIDs=%v\n", count, id)
			// 建立一個 *Node，ID 設為 nodeID
			newNode := &Node{
				ID: id,
				// 若要帶入 shape、connections 等，可以在此做處理或查表
			}
			path.Weight +=1
			path.Nodes = append(path.Nodes, newNode)
		
		}
		k := new_KPath(1, source, target)
        // 把剛生成的 p 加入 k.Paths
        k.Paths = append(k.Paths, path)
		return graph.Path[0],k
	}
	
	return nil,nil
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

// 依據節點 ID 切片，生成一個 *Path，其中每個節點都是新的 Node 物件
func ConvertIDsToPath(ids []int,topo *topology.Topology) *Path {
    p := &Path{
        Nodes:  make([]*Node, 0, len(ids)),
        Weight: 0,
    }
    // 逐節點轉換
    for i, id := range ids {
        realNode := topo.GetNodeByID(id)
        if realNode == nil {
            // 找不到對應節點，表示資料或路徑不一致，可直接 return p 或 return nil
            fmt.Printf("Warning: Node with ID=%d not found in topology.\n", id)
            continue
        }

        // 建立 path.Node，帶入 shape 與必要資訊
        node := &Node{
            ID:   realNode.ID,
            Shape: realNode.Shape, // 如果在 realNode 中有 shape
            Connections: make([]*Connection, 0), // 稍後若需要也可填
        }
        p.Nodes = append(p.Nodes, node)
        
        // 如果要順手加「邊的資訊」：可以在這裡 or 之後做
        if i < len(ids)-1 {
            // 尋找 (id -> ids[i+1]) 這條連線
            nextID := ids[i+1]
            conn := findConnectionInNode(realNode, nextID)
            if conn != nil {
                // 新增到 node.Connections or p.Nodes[i].Connections
                node.Connections = append(node.Connections, &Connection{
                    FromNodeID: conn.FromNodeID,
                    ToNodeID:   conn.ToNodeID,
                    Cost:       conn.Cost,
                })
                
                // 如果要用來計算路徑權重
                p.Weight += conn.Cost
            }
        }
    }
    return p
}

func findConnectionInNode(node *topology.Node, toID int) *topology.Connection {
    for _, c := range node.Connections {
        if c.ToNodeID == toID {
            return c
        }
    }
    return nil
}