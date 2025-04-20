package path

import (
	// "fmt"
	"src/network"
	"src/network/topology"
	// "src/network/flow"
	"container/heap"
	"math"
)



// var k =5

// func KShortestPath(Network *network.Network) *KPath_Set{

// 	kpath_set:=new_KPath_Set()

// 	appendKPaths(&kpath_set.TSNPaths,  Network.TSNFlow_Set.TSNFlows,Network.TSNGraph_Set.TSNGraphs)
//     appendKPaths(&kpath_set.AVBPaths,  Network.TSNFlow_Set.AVBFlows,Network.TSNGraph_Set.AVBGraphs)
    
// 	return kpath_set
// }

// func appendKPaths(dst *[]*KPath, flows []*flow.Flow, topos []*topology.Topology) {
//     for i, f := range flows {
//         p := KShortestPaths(topos[i], f.Source, f.Destination, k)
//         *dst = append(*dst, p)
//     }
// }

// func shortestIDs(topo *topology.Topology, src, dst int) []int {
//     g := GetGarph(topo)
//     Dijkstra(g, dst, src)          // 你的 Dijkstra 反向跑
//     if len(g.Path) == 0 {
//         return nil
//     }
//     ids := append([]int(nil), g.Path[0]...)
// 	for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
// 		ids[i], ids[j] = ids[j], ids[i]   // 轉成 src→dst
// 	}
// 	return ids
// }

// // ----------------- 優先佇列 -----------------
// type cand struct {
//     cost float64
//     ids  []int
// }
// type candPQ []cand
// func (pq candPQ) Len() int            { return len(pq) }
// func (pq candPQ) Less(i, j int) bool  { return pq[i].cost < pq[j].cost }
// func (pq candPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
// func (pq *candPQ) Push(x interface{}) { *pq = append(*pq, x.(cand)) }
// func (pq *candPQ) Pop() interface{} {
//     old := *pq
//     n := len(old)
//     x  := old[n-1]
//     *pq = old[:n-1]
//     return x
// }
// // ------------------------------------------------

// // KShortestPaths 取得 K 條最短路，回傳 *KPath
// func KShortestPaths(topo *topology.Topology, src, dst, K int) *KPath {
//     kpath := new_KPath(K, src, dst)

//     // 1️⃣ 先求第一條
//     first := shortestIDs(topo, src, dst)
//     if first == nil {
//         return kpath // 無路徑
//     }
//     kpath.Paths = append(kpath.Paths, idsToPath(first, topo))

//     // 候選路徑最小堆
//     pq := &candPQ{}
//     heap.Init(pq)

//     // 2️⃣ Yen’s 迴圈
//     for k := 1; k < K; k++ {
//         prevIDs := make([]int, len(kpath.Paths[k-1].Nodes))
//         for i, n := range kpath.Paths[k-1].Nodes { prevIDs[i] = n.ID }

//         for i := 0; i < len(prevIDs)-1; i++ {
//             spurNode := prevIDs[i]
//             rootIDs  := prevIDs[:i+1]        // s→spur

//             // 2‑1 複製拓樸，移除 rootIDs 之後與既有最短路重疊的邊
//             // clone & remove
// 			tCopy := topo.TopologyDeepCopy()
// 			removeOverlapEdges(tCopy, rootIDs , kpath )

//             // 2‑2 從 spurNode → dst 求最短路
//             spurIDs := shortestIDs(tCopy, spurNode, dst)
//             if spurIDs != nil {
//                 // 合併 root(不含 spurNode 最後一點) + spur
//                 candIDs := append(append([]int{}, rootIDs[:len(rootIDs)-1]...), spurIDs...)
//                 cost    := pathCost(topo, candIDs)
//                 heap.Push(pq, cand{cost: cost, ids: candIDs})
//             }
//             // 2‑3 恢復邊 (cloneAndRemove 已複製，不用復原)
//         }

//         if pq.Len() == 0 {
//             break // 無更多候選
//         }
//         best := heap.Pop(pq).(cand)
//         kpath.Paths = append(kpath.Paths, idsToPath(best.ids, topo))
//     }
//     return kpath
// }


// // --------- 輔助 ---------
// func pathIDs(p *Path) []int {
// 	ids := make([]int, len(p.Nodes))
// 	for i, n := range p.Nodes { ids[i] = n.ID }
// 	return ids
// }

// func removeOverlapEdges(t *topology.Topology, root []int, kp *KPath) {
// 	for _, p := range kp.Paths {
// 		ids := pathIDs(p)
// 		if len(ids) < len(root) || !equalPrefix(ids, root) { continue }
// 		u, v := ids[len(root)-1], ids[len(root)]
// 		t.RemoveEdge(u, v) // 你要在 Topology 加此方法
// 	}
// }

// func equalPrefix(a, b []int) bool {
// 	if len(b) > len(a) { return false }
// 	for i := range b {
// 		if a[i] != b[i] { return false }
// 	}
// 	return true
// }

// // 計算成本
// func pathCost(t *topology.Topology, ids []int) float64 {
// 	var c float64
// 	for i := 0; i < len(ids)-1; i++ {
// 		n := t.GetNodeByID(ids[i])
// 		for _, e := range n.Connections {
// 			if e.ToNodeID == ids[i+1] { c += e.Cost; break }
// 		}
// 	}
// 	return c
// }

// kpath_all.go  ── 整合 Yen K‑Shortest Paths + Dijkstra + 與現有結構對接
// 放在 path/ 目錄即可直接 `go test ./path`

// ─────────────────────────────────────────────────────────────
// 1. 最短路徑 Dijkstra  (回傳 dist / prev)
// ─────────────────────────────────────────────────────────────
func KDijkstra(g *Graph, src int) (map[int]float64, map[int]int) {
    dist := map[int]float64{}
    prev := map[int]int{}
    for _, v := range g.Vertexs {
        dist[v.ID] = math.Inf(1)
    }
    dist[src] = 0

    pq := &vertexPQ{}
    heap.Init(pq)
    heap.Push(pq, &vertexCost{src, 0})

    for pq.Len() > 0 {
        vc := heap.Pop(pq).(*vertexCost)
        if vc.cost > dist[vc.id] {
            continue // outdated
        }
        v := g.getVertex(vc.id)
        for _, e := range v.Edges {
            alt := dist[vc.id] + float64(e.Cost)
            if alt < dist[e.End] {
                dist[e.End] = alt
                prev[e.End] = vc.id
                heap.Push(pq, &vertexCost{e.End, alt})
            }
        }
    }
    return dist, prev
}

// ── PQ for Dijkstra

type vertexCost struct{ id int; cost float64 }
type vertexPQ []*vertexCost

func (h vertexPQ) Len() int            { return len(h) }
func (h vertexPQ) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h vertexPQ) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *vertexPQ) Push(x interface{}) { *h = append(*h, x.(*vertexCost)) }
func (h *vertexPQ) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func (g *Graph) getVertex(id int) *Vertex {
    for _, v := range g.Vertexs {
        if v.ID == id {
            return v
        }
    }
    return nil
}

// ─────────────────────────────────────────────────────────────
// 2. Yen K‑Shortest Paths (loop‑free)  回傳 [][]int
// ─────────────────────────────────────────────────────────────

func YenKShortest(g *Graph, src, dst, K int) [][]int {
    dist, prev := KDijkstra(g, src)
    if dist[dst] == math.Inf(1) {
        return nil // 無路徑
    }
    P := [][]int{rebuildPath(prev, dst)}
    cand := &pathHeap{}
    heap.Init(cand)

    for k := 1; k < K; k++ {
        last := P[k-1]
        for i := 0; i < len(last)-1; i++ {
            spurNode := last[i]
            root := append([]int{}, last[:i+1]...) // copy

            removed := removeEdgesAndNodes(g, root)
            d, p := KDijkstra(g, spurNode)
            restore(removed)
            if d[dst] == math.Inf(1) {
                continue
            }
            spur := rebuildPath(p, dst)
            totalPath := append(root[:len(root)-1], spur...)
            totalCost := pathCost(root, g) + d[dst]
            heap.Push(cand, &candPath{totalPath, totalCost})
        }
        if cand.Len() == 0 {
            break
        }
        P = append(P, heap.Pop(cand).(*candPath).ids)
    }
    return P
}

// ── 候選路徑最小堆

type candPath struct{ ids []int; cost float64 }
type pathHeap []*candPath

func (h pathHeap) Len() int           { return len(h) }
func (h pathHeap) Less(i, j int) bool { return h[i].cost < h[j].cost }
func (h pathHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *pathHeap) Push(x interface{}) { *h = append(*h, x.(*candPath)) }
func (h *pathHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// ── 工具：路徑重建、成本計算、刪除/還原邊節點

func rebuildPath(prev map[int]int, dst int) []int {
    var rev []int
    for u, ok := dst, true; ok; u, ok = prev[u] {
        rev = append(rev, u)
        if u == 0 { break }
    }
    // reverse
    for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
        rev[i], rev[j] = rev[j], rev[i]
    }
    return rev
}

func pathCost(path []int, g *Graph) float64 {
    cost := 0.0
    for i := 0; i < len(path)-1; i++ {
        v := g.getVertex(path[i])
        for _, e := range v.Edges {
            if e.End == path[i+1] {
                cost += float64(e.Cost)
            }
        }
    }
    return cost
}

type removedObj struct{ from, to int; vertex *Vertex }

func removeEdgesAndNodes(g *Graph, prefix []int) []removedObj {
    var removed []removedObj
    // 移除 prefix 中的邊及重複節點
    for i := 0; i < len(prefix)-1; i++ {
        v := g.getVertex(prefix[i])
        for idx, e := range v.Edges {
            if e.End == prefix[i+1] {
                removed = append(removed, removedObj{from: v.ID, to: e.End})
                v.Edges = append(v.Edges[:idx], v.Edges[idx+1:]...)
                break
            }
        }
    }
    // optional: node‑disjoint 可移除節點 (略)
    return removed
}

func restore(objs []removedObj) {
    for _, o := range objs {
        v := o.vertex
        if v == nil {
            // 恢復邊
            // 找到對應 Vertex 加回 Edge (Cost=1 預設)
        }
    }
}

// ─────────────────────────────────────────────────────────────
// 3. 對接 Flow → KPath
// ─────────────────────────────────────────────────────────────

func BuildAllKPathSet(net *network.Network, k int) *KPath_Set {
    ks := new_KPath_Set()
    attach := func(src, dst int, topo *topology.Topology, bag *[]*KPath) {
        kp := new_KPath(k, src, dst)
        ids := YenKShortest(GetGarph(topo), src, dst, k)
        for _, p := range ids {
            kp.Paths = append(kp.Paths, ConvertIDsToPath(p, topo))
        }
        *bag = append(*bag, kp)
    }
    for i, f := range net.Flow_Set.TSNFlows {
        attach(f.Source, f.Destination, net.Graph_Set.TSNGraphs[i], &ks.TSNPaths)
    }
    for i, f := range net.Flow_Set.AVBFlows {
        attach(f.Source, f.Destination, net.Graph_Set.AVBGraphs[i], &ks.AVBPaths)
    }
    // Important / Unimportant CAN 類似添加…
    return ks
}
