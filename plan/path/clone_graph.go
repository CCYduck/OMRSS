package path


import "src/network/topology"
// ------------------------------------------------------------
// Graph deep-clone  -- 只複製 Dijkstra/ Yen 會用到的欄位
// ------------------------------------------------------------
func (g *Graph) Clone() *Graph {
	out := &Graph{V: make(map[int]*Vertex, len(g.V))}
	for id, v := range g.V {
		nv := *v
		nv.Edges = make([]*Edge, len(v.Edges))
		for i, e := range v.Edges {
			ne := *e
			nv.Edges[i] = &ne
		}
		out.V[id] = &nv
	}
	return out
}

// 刪 (u → v) 單向邊
func (g *Graph) RemoveEdge(u, v int) {
	if vert := g.V[u]; vert != nil {
		e2 := vert.Edges[:0]
		for _, e := range vert.Edges {
			if !(e.Strat == u && e.End== v) {
				e2 = append(e2, e)
			}
		}
		vert.Edges = e2
	}
}

// 刪頂點 + incident edges
func (g *Graph) RemoveVertex(id int) { delete(g.V, id); for _, v := range g.V {
	e2 := v.Edges[:0]
	for _, e := range v.Edges {
		if e.Strat != id { e2 = append(e2, e) }
	}
	v.Edges = e2
} }

// 把 topology 轉成簡單圖（單向邊）
func BuildGraphFromTopology(t *topology.Topology) *Graph {
	g := &Graph{V: map[int]*Vertex{}}

	add := func(u,v int, cost float64) {
		if g.V[u] == nil { g.V[u] = &Vertex{ID:u} }
		g.V[u].Edges = append(g.V[u].Edges, &Edge{Strat:u, End:v, Cost:int(cost)})
	}

	iterNode := func(n *topology.Node){
		for _, c := range n.Connections { add(c.FromNodeID, c.ToNodeID, c.Cost) }
	}
	for _, n := range append(append(t.Talker, t.Switch...), t.Listener...) {
		iterNode(n)
	}

	return g
}