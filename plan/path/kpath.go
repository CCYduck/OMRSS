package path

// func KShortestPaths(g0 *Graph, s, t, K int) *KPath {
// 	kp := &KPath{K: K, Source: s, Target: t}

// 	// 1️⃣ 先求第一條最短路
// 	p0, ok := shortestPath(g0, s, t)
// 	if !ok {
// 		return kp // 無路徑
// 	}
// 	kp.Paths = append(kp.Paths, p0)

// 	// 2️⃣ 儲存候選路徑的最小堆
// 	type cand struct{ cost float64; path []int }
// 	candPQ := &priorityQueue{} // 實作 container/heap
// 	heap.Init(candPQ)

// 	// 3️⃣ Yen’s 主迴圈
// 	for k := 1; k < K; k++ {

// 		prevPath := kp.Paths[k-1] // P_{k-1}

// 		for i := 0; i < len(prevPath)-1; i++ {
// 			spurNode := prevPath[i]          // 偏差點
// 			rootPath := prevPath[:i+1]       // s → spurNode

// 			// 3‑1 暫時移除：rootPath 中與既有最短路重複的邊 / 節點
// 			removed := removeRootEdges(g0, rootPath, kp.Paths)

// 			// 3‑2 從 spurNode 到 t 再跑一次最短路
// 			spurPath, ok := shortestPath(g0, spurNode, t)
// 			if ok {
// 				candidate := append([]int{}, rootPath[:len(rootPath)-1]...)
// 				candidate = append(candidate, spurPath...)
// 				cost := pathCost(g0, candidate)
// 				heap.Push(candPQ, cand{cost: cost, path: candidate})
// 			}

// 			// 3‑3 恢復被移除的邊 / 節點
// 			restoreEdges(g0, removed)
// 		}

// 		if candPQ.Len() == 0 { // 再也沒有候選路徑
// 			break
// 		}
// 		// 4️⃣ 取 cost 最小者作為第 k 條最短路
// 		best := heap.Pop(candPQ).(cand)
// 		kp.Paths = append(kp.Paths, best.path)
// 	}

// 	return kp
// }