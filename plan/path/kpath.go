package path

import (
	"sort"
	"fmt"
	"src/network"
	// "src/network/topology"
)

func KShortestPath(Network *network.Network) *KPath_Set{
	const k = 3    		// 你要幾條路
	kpath_set := new_KPath_Set()                   

	// -------- lookup table，用來判斷 (src,dst) 是否已經建立過 ----------
    type sd struct{ s, d int }
	// done := make(map[sd]bool)	// 全域收集
	tsndone := make(map[sd]bool)	// 全域收集
	avbdone := make(map[sd]bool)	// 全域收集
	can2tsndone := make(map[sd]bool)	// 全域收集


	for idx, flow := range Network.Flow_Set.TSNFlows {
		pair := sd{flow.Source, flow.Destination}
        if tsndone[pair] {                 // 已經做過 → skip
            continue
        }

		topo := Network.Graph_Set.TSNGraphs[idx]
		g    := GetGarph(topo)
		kp   := BuildKPath(k, flow.Source, flow.Destination, g)
		kpath_set.TSNPaths = append(kpath_set.TSNPaths, kp)
		// done[pair] = true 
		tsndone[pair] = true 
	}
	
	// -------- AVB ----------
	for idx, flow := range Network.Flow_Set.AVBFlows {
		pair := sd{flow.Source, flow.Destination}
        if avbdone[pair] {                 // 已經做過 → skip
            continue
        }
		topo := Network.Graph_Set.AVBGraphs[idx]
		g    := GetGarph(topo)
		kp   := BuildKPath(k, flow.Source, flow.Destination, g)
		kpath_set.AVBPaths = append(kpath_set.AVBPaths, kp)
		// done[pair] = true 
		avbdone[pair] = true
	}
	
	// -------- CAN→TSN (封裝流) ----------
	for _, m := range Network.Flow_Set.Encapsulate {   // 每種封裝方法
		for _, f := range m.CAN2TSNFlows {             // 每條 CAN→TSN flow
			pair := sd{f.Source, f.Destination}
            if can2tsndone[pair] { continue }

			topo := Network.Graph_Set.GetGarphBySD(f.Source, f.Destination)
			g    := GetGarph(topo)
			kp   := BuildKPath(k, f.Source, f.Destination, g)
			kpath_set.CAN2TSNPaths = append(kpath_set.CAN2TSNPaths, kp)
			can2tsndone[pair] = true 
		}
	}
	return kpath_set
}

func BuildKPath(k int, src, dst int, g *Graph) *KPath {
	kp := new_KPath(k, src, dst)
	kp.Paths = YenKPaths(g, src, dst, k)
	return kp
}

// YenKPaths 回傳 K 條最短簡單路，已依 Weight 由小到大
func YenKPaths(base *Graph, src, dst, K int) []*Path {
	A := make([]*Path, 0, K) // 已確定
	B := make([]*Path, 0)    // 候選

	first := oneShortest(base, src, dst)
	if first == nil { return nil }
	A = append(A, first)

	for k := 1; k < K; k++ {
		last := A[k-1]
		for i := 0; i < len(last.Nodes)-1; i++ { // 每個 spur node
			spur := last.Nodes[i].ID
			rootIDs := last.PrefixIDs(i) 	//把node結構變成list結構

			// 1. 建殘圖
			g := base.Clone()
			for _, rid := range rootIDs[:len(rootIDs)-1] {
				g.RemoveVertex(rid)
			}
			for _, p := range A {
				if len(p.Nodes) <= i { continue }
				if equalSlice(rootIDs, p.PrefixIDs(i)) {
					u, v := p.Nodes[i].ID, p.Nodes[i+1].ID
					g.RemoveEdge(u, v)
				}
			}

			// 2. spur → dst 最短路
			spurPath := oneShortest(g, spur, dst)
			if spurPath == nil { continue }

			// 3. 組 full path
			full := join(rootIDs, spurPath)
			fmt.Println(full.Nodes)
			B = append(B, full)
		}
		if len(B) == 0 { break }
		sort.Slice(B, func(i, j int) bool { return B[i].Weight < B[j].Weight })
		A = append(A, B[0])
		B = B[1:]
	}
	return A
}

// ---------- internal helpers ----------

func oneShortest(g *Graph, s, t int) *Path {
	ret := Dijkstra(g, s, t)
	if len(ret.Path) == 0 { return nil }

	ids := ret.Path[0]
	p := &Path{Weight: float64(len(ids) - 1)}

	for i, id := range ids {
        n := &Node{ID: id}
        if i < len(ids)-1 {
            n.Connections = append(n.Connections, &Connection{
                FromNodeID: id,
                ToNodeID:   ids[i+1],
                Cost:       1,
            })
        }
        // 也可補前向邊、反向邊都行
        p.Nodes = append(p.Nodes, n)
    }
	return p
}

func join(rootIDs []int, spur *Path) *Path {
	out := &Path{}
	for _, id := range rootIDs[:len(rootIDs)-1] {
		out.Nodes = append(out.Nodes, &Node{ID: id})
	}
	out.Nodes = append(out.Nodes, spur.Nodes...)
	out.Weight = float64(len(out.Nodes) - 1)
	return out
}
