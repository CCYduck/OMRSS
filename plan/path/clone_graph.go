package path

// ------------------------------------------------------------
// Graph deep-clone  -- 只複製 Dijkstra/ Yen 會用到的欄位
// ------------------------------------------------------------
func (g *Graph) Clone() *Graph {
	out := &Graph{ToVertex: g.ToVertex}
	out.Vertexs = make([]*Vertex, len(g.Vertexs))
	for i, v := range g.Vertexs {
		nv := *v
		nv.Edges = make([]*Edge, len(v.Edges))
		for j, e := range v.Edges {
			ne := *e
			nv.Edges[j] = &ne
		}
		out.Vertexs[i] = &nv
	}
	return out
}

// 刪 (u → v) 單向邊
func (g *Graph) RemoveEdge(u, v int) {
	if vert := g.FindVertex(u); vert != nil {
		keep := vert.Edges[:0]
		for _, e := range vert.Edges {
			if !(e.Strat == u && e.End == v) {
				keep = append(keep, e)
			}
		}
		vert.Edges = keep
	}
}

// 刪頂點 + incident edges
func (g *Graph) RemoveVertex(id int) {
	// 1. remove vertex
	v2 := g.Vertexs[:0]
	for _, v := range g.Vertexs {
		if v.ID != id {
			v2 = append(v2, v)
		}
	}
	g.Vertexs = v2
	// 2. remove edges to that vertex
	for _, v := range g.Vertexs {
		e2 := v.Edges[:0]
		for _, e := range v.Edges {
			if e.End != id {
				e2 = append(e2, e)
			}
		}
		v.Edges = e2
	}
}
