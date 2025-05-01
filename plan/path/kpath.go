package path

import (
	"sort"
	"fmt"
	"src/network"
	"math"
	"src/network/topology"
)

func BuildKPath(k, src, dst int, topo *topology.Topology)*KPath{
	
	g  := BuildGraphFromTopology(topo)
	ids := YenKPaths(g,src,dst,k)

	kp := new_KPath(k,src,dst)
	for _, idPath := range ids {
		kp.Paths = append(kp.Paths, ConvertIDsToPath_k(idPath.IDs, topo))
	}
	return kp
}

func KShortestPath(Network *network.Network) *KPath_Set{
	const k = 3    		// 你要幾條路
	kpath_set := new_KPath_Set()                   

	// -------- lookup table，用來判斷 (src,dst) 是否已經建立過 ----------
    type sd struct{ s, d int }
	// done := make(map[sd]bool)	// 全域收集
	tsnDone := make(map[sd]bool)	// 全域收集
	avbDone := make(map[sd]bool)	// 全域收集
	c2tDone	:= make(map[sd]bool)	// 全域收集


	for idx, flow := range Network.Flow_Set.TSNFlows {
		key := sd{flow.Source, flow.Destination}
		if tsnDone[key] {continue}
		kp := BuildKPath(k, flow.Source, flow.Destination, Network.Graph_Set.TSNGraphs[idx])
		kpath_set.TSNPaths = append(kpath_set.TSNPaths, kp)
		tsnDone[key]=true
	}
	
	// -------- AVB ----------
	for idx, flow := range Network.Flow_Set.AVBFlows {
		key := sd{flow.Source, flow.Destination}
        if avbDone[key]{continue}
		kp := BuildKPath(k, flow.Source,flow.Destination, Network.Graph_Set.AVBGraphs[idx])
		kpath_set.AVBPaths = append(kpath_set.AVBPaths,kp)
		avbDone[key]=true
	}

	// -------- CAN→TSN (封裝流) ----------
	for _, m := range Network.Flow_Set.Encapsulate {   // 每種封裝方法
		for _, f := range m.CAN2TSNFlows {             // 每條 CAN→TSN flow
			key := sd{f.Source,f.Destination}
			if c2tDone[key]{continue}
			topo := Network.Graph_Set.GetGarphBySD(f.Source,f.Destination)
			kp   := BuildKPath(k,f.Source,f.Destination, topo)
			kpath_set.CAN2TSNPaths = append(kpath_set.CAN2TSNPaths,kp)
			c2tDone[key]=true
		}
	}
	return kpath_set
}


// YenKPaths 回傳 K 條最短簡單路，已依 Weight 由小到大
func YenKPaths(g *Graph, src, dst, K int) []*Path {
	A := []*Path{}
	first := oneShortest(g, src, dst)
	if first==nil { return nil }
	A = append(A, first)
	B := []*Path{}

	for ki:=1; ki<k; ki++ {
		last := A[ki-1]
		for i:=0;i<len(last.IDs)-1;i++ {   // 每個 spur node
			spurNode := last.IDs[i]
			rootPath := last.IDs[:i+1]

			// 殘圖
			g2 := g.Clone()
			for _, p := range A {
				if len(p.IDs)>i && equal(rootPath, p.IDs[:i+1]) {
					g2.RemoveEdge(p.IDs[i], p.IDs[i+1])
				}
			}
			for _, rid := range rootPath[:len(rootPath)-1] {
				g2.RemoveVertex(rid)
			}
			spur := oneShortest(g2, spurNode, dst)
			if spur==nil { continue }

			full := append(append([]int{}, rootPath[:len(rootPath)-1]...), spur.IDs...)
			B = append(B, &Path{IDs: full, Weight: float64(len(full)-1)})
		}
		if len(B)==0 { break }
		sort.Slice(B, func(i,j int)bool{ return B[i].Weight<B[j].Weight })
		A = append(A, B[0])
		B = B[1:]
	}
	return A
}

func equal(a,b []int)bool{
	if len(a)!=len(b) { return false }
	for i:=range a { if a[i]!=b[i] { return false } }
	return true
}
// ---------- internal helpers ----------


func oneShortest(g *Graph, s, t int) *Path {
	const inf = math.MaxInt32
	dist := map[int]int{}
	prev := map[int]int{}
	for id := range g.V { dist[id] = inf }
	dist[s] = 0
	visited := map[int]bool{}

	for len(visited) < len(g.V) {
		// 找未訪問且 dist 最小的頂點
		minID, minV := -1, inf
		for id, d := range dist {
			if !visited[id] && d < minV {
				minV, minID = d, id
			}
		}
		if minID == -1 || minID == t { break }
		visited[minID] = true
		for _, e := range g.V[minID].Edges {
			if visited[e.End] { continue }
			if alt := dist[minID] + e.Cost; alt < dist[e.End] {
				dist[e.End] = alt
				prev[e.End] = minID
			}
		}
	}

	if _, ok := dist[t]; !ok || dist[t]==inf { return nil }

	// 回溯路徑
	ids := []int{t}
	for cur := t; cur != s; {
		cur = prev[cur]
		ids = append(ids, cur)
	}
	// 反轉
	for i,j := 0,len(ids)-1; i<j; i,j = i+1,j-1 {
		ids[i],ids[j] = ids[j],ids[i]
	}
	return &Path{IDs: ids, Weight: float64(dist[t])}
}

func ConvertIDsToPath_k(ids []int, topo *topology.Topology) *Path {
	p := &Path{
        IDs:    ids,                         // ← (1) 補上
        Nodes:  make([]*Node, 0, len(ids)),
        Weight: 0,
    }
    for i, id := range ids {
        real := topo.GetNodeByID(id)
        if real == nil {                    // ← (2) 直接認定路徑無效
            fmt.Printf("Node %d not found, abort path\n", id)
            return nil
        }
        node := &Node{
            ID:    real.ID,
            Shape: real.Shape,
        }
        p.Nodes = append(p.Nodes, node)

        if i < len(ids)-1 {
            if c := findConnectionInNode(real, ids[i+1]); c != nil {
                node.Connections = append(node.Connections, &Connection{
                    FromNodeID: c.FromNodeID,
                    ToNodeID:   c.ToNodeID,
                    Cost:       c.Cost,
                })
                p.Weight += c.Cost          // ← (3) 累加 hop 成本
            }
        }
    }
    return p
}
