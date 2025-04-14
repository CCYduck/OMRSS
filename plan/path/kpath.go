package path

// import (
// 	// "fmt"
// 	"src/network"
// 	"src/network/topology"
// 	"src/network/flow"
// 	"container/heap"
// )

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