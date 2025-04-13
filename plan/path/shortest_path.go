package path

import (
	"fmt"
	"src/network"
	"src/network/topology"
	"src/network/flow"
)
//主要的function

func BestPath(Network *network.Network) *Path_set{

	path_set:=new_Path_Set()

	appendPaths(&path_set.TSNPath,  Network.TSNFlow_Set.TSNFlows,             Network.TSNGraph_Set.TSNGraphs)
    appendPaths(&path_set.AVBPath,  Network.TSNFlow_Set.AVBFlows,             Network.TSNGraph_Set.AVBGraphs)
    appendPaths(&path_set.ImportCanPath, Network.CANFlow_Set.ImportantCANFlows, Network.CANGraph_Set.Important_CANGraphs)
    appendPaths(&path_set.UnimportCanPath, Network.CANFlow_Set.UnimportantCANFlows, Network.CANGraph_Set.Unimportant_CANGraphs)

	return path_set
}

func ShortestPathAsPath(topo *topology.Topology, src, dst int) *Path {
    g := GetGarph(topo)
    g.ToVertex = dst
    Dijkstra(g, dst, src)

    if len(g.Path) == 0 {
        return nil
    }
    return idsToPath(g.Path[0], topo)
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

// func (v2v *V2V) GetV2VEdge(terminal int) (*V2VEdge, bool) {
// 	for _, edge := range v2v.V2VEdges {
// 		if edge.FromVertex == terminal {
// 			return edge, false
// 		}
// 	}
// 	return &V2VEdge{FromVertex: terminal}, true
// }

// func (v2vedge *V2VEdge) InV2VEdge(terminal int) bool {
// 	for _, graph := range v2vedge.Graphs {
// 		if graph.ToVertex == terminal {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (v2vedge *V2VEdge) GetV2VPath(terminal int) [][]int {
// 	var path [][]int
// 	for _, graph := range v2vedge.Graphs {
// 		if graph.ToVertex == terminal {
// 			path = graph.Path
// 		}
// 	}
// 	return path
// }

// 依據節點 ID 切片，生成一個 *Path，其中每個節點都是新的 Node 物件
func idsToPath(ids []int,topo *topology.Topology) *Path {
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



func appendPaths(dst *[]*Path, flows []*flow.Flow, topos []*topology.Topology) {
    for i, f := range flows {
        p := ShortestPathAsPath(topos[i], f.Source, f.Destination)
        *dst = append(*dst, p)
    }
}