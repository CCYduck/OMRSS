package path

import (
	"fmt"
	"src/network"
	"src/network/topology"
	
)

//主要的function
func BestPath(Network *network.Network) *Path_set {

	path_set := new_Path_Set()
	for nth, flow := range Network.Flow_Set.TSNFlows {
		// fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination, Network.Graph_Set.TSNGraphs[nth])
		//fmt.Printf("Flow: Source=%d, Destination=%v", flow.Source, flow.Destination)
		path := saveShortestPathsToGraph(flow.Source, flow.Destination, Network.Graph_Set.TSNGraphs[nth])
		if path != nil {
			// fmt.Println("Best Path:")
			// path.Show_Path()
			
		} else {
			fmt.Println("No path found.")
		}
		path_set.TSNPath = append(path_set.TSNPath, path)
		
		// 你可以存到一個 KPath_Set 或印出來
		// fmt.Printf("KPath: Source=%d, Target=%d, NodeIDs=%v\n", k.Source, k.Target, path)

	}

	for nth, flow := range Network.Flow_Set.AVBFlows {
		// fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination, Network.Graph_Set.TSNGraphs[nth])
		//fmt.Printf("Flow: Source=%d, Destination=%v", flow.Source, flow.Destination)
		path := saveShortestPathsToGraph(flow.Source, flow.Destination, Network.Graph_Set.AVBGraphs[nth])
		if path != nil {
			// fmt.Println("Best Path:")
			// path.Show_Path()
			
		} else {
			fmt.Println("No path found.")
		}
		path_set.AVBPath = append(path_set.AVBPath, path)
		
		// 你可以存到一個 KPath_Set 或印出來
		// fmt.Printf("KPath: Source=%d, Target=%d, NodeIDs=%v\n", k.Source, k.Target, path)

	}

	for _, method := range Network.Flow_Set.Encapsulate {
		type sd struct{ s int; d int}
		userage_path := make(map[sd]*Path)	// 全域收集
		// fmt.Println(method.Method_Name)
		// fmt.Printf("Flow: Source=%d, Destination=%v, Topology=%v\n", flow.Source, flow.Destination, Network.Graph_Set.TSNGraphs[nth])
		//fmt.Printf("Flow: Source=%d, Destination=%v", flow.Source, flow.Destination)
		for _,flow := range method.CAN2TSNFlows{		
			key := sd{flow.Source, flow.Destination}
			// 1. 檢查 key 是否已經存在
			if sp, ok := userage_path[key]; ok {
				// 已經算過，直接使用
				path_set.CAN2TSNPath = append(path_set.CAN2TSNPath, sp)
			}else{
				sp := saveShortestPathsToGraph(flow.Source, flow.Destination, Network.Graph_Set.GetGarphBySD(flow.Source, flow.Destination))
				if sp != nil {sp.Method = method.Method_Name}
				// fmt.Println(method.Method_Name)
				path_set.CAN2TSNPath = append(path_set.CAN2TSNPath, sp)
				userage_path[key] = sp
			}

		}
		
		// 你可以存到一個 KPath_Set 或印出來
		// fmt.Printf("KPath: Source=%d, Target=%d, NodeIDs=%v\n", k.Source, k.Target, path)

	}
	

	return path_set
}



func saveShortestPathsToGraph(source int, target int, t *topology.Topology) *Path {
	// Check if this path has already been taken
	graph := GetGarph(t)
	graph.ToVertex = target
	graph = Dijkstra(graph, target, source)

	if len(graph.Path) > 0 {
		path := new_Path()
		for count, id := range graph.Path[0] {
			// 建立一個 *Node，ID 設為 nodeID
			newNode := &Node{
				ID: id,
				// 若要帶入 shape、connections 等，可以在此做處理或查表
			}

			//fmt.Printf("Source=%d,NodeIDs=%v\n", count, id)
			if count != len(graph.Path[0])-1 {
				newfrontConn := &Connection{
					FromNodeID: id,
					ToNodeID:   graph.Path[0][count+1], // next
					Cost:       0,
				}
				newNode.Connections = append(newNode.Connections, newfrontConn)
			}

			if count != 0 {
				newbackConn := &Connection{
					FromNodeID: id,
					ToNodeID:   graph.Path[0][count-1], // before
					Cost:       0,
				}
				newNode.Connections = append(newNode.Connections, newbackConn)
			}

			path.Weight += 1
			path.Nodes = append(path.Nodes, newNode)

		}

		return path
	}

	return nil
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
func ConvertIDsToPath(ids []int, topo *topology.Topology) *Path {
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
			ID:          realNode.ID,
			Shape:       realNode.Shape,         // 如果在 realNode 中有 shape
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

// func (path_set *Path_set)checkListenerAndTalker(source int, destination int)bool{

// 	for _, p := range path_set.CAN2TSNPath {

// 		s := p.GetNodeByID(source)
// 		if s == nil {          // ←★★ 防護：source 不在此 path
// 			continue
// 		}
// 		d := p.GetNodeByID(destination)
// 		if d == nil {          // ←★★ 防護：dest 不在此 path
// 			continue
// 		}
// 		// 兩個都存在才算「已經有這條 path」
// 		return true
// 	}
// 	return false
// }

// func (ps *Path_set) GetPathByMethod(m string) []*Path {
//     var ret []*Path
//     for _, p := range ps.CAN2TSNPath { // 依你的欄位名稱
//         if p == nil || p.Method == m {
//             ret = append(ret, p)  // nil 也放
//         }
//     }
//     return ret
// }

