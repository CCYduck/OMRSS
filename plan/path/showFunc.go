package path

import (
	"fmt"
)

func (v2v *V2V) Show_V2Vs() {
	for _, v2vedge := range v2v.V2VEdges {
		v2vedge.Show_VertexToVertex()
	}
}

func (v2vedge *V2VEdge) Show_VertexToVertex() {
	for _, graph := range v2vedge.Graphs {
		fmt.Printf("From Vertex: %d\n", v2vedge.FromVertex)
		graph.Show_Graph_Path()
	}
}

func (graph *Graph) Show_Graph_Path() {
	fmt.Println("Graph vertices:")
	for id, vertex := range graph.Vertexs {
		fmt.Printf("Vertex %d: %+v\n", id, vertex)
		fmt.Println("Graph edges:")
		for _, edge := range vertex.Edges {
			fmt.Printf("Edge from %d to %d with cost %d\n", edge.Strat, edge.End, edge.Cost)
		}
	}
}

func (Paths *KPath_Set) Show_KPath_Set() {

	for ind, KPath := range Paths.TSNPaths {
		fmt.Printf("\nTSN #%d: %d→%d  (K=%d)\n", ind, KPath.Source, KPath.Target, KPath.K)
		KPath.Show_KPath()

		// break
	}

	for ind, KPath := range Paths.AVBPaths {
		fmt.Printf("\nAVB #%d: %d→%d  (K=%d)\n", ind, KPath.Source, KPath.Target, KPath.K)
		KPath.Show_KPath()

		// break
	}

	for ind, KPath := range Paths.CAN2TSNPaths {
		fmt.Printf("\nCAN2TSN #%d: %d→%d  (K=%d)\n", ind, KPath.Source, KPath.Target, KPath.K)
		KPath.Show_KPath()

		// break
	}

}

func (Paths *Path_set) Show_Path_Set() {

	for index, path := range Paths.TSNPath {
		fmt.Printf("TSN Path %d\n", index)
		path.Show_Path()

		// break
	}

	for index, path := range Paths.AVBPath {
		fmt.Printf("\nAVB Path %d \n", index)
		path.Show_Path()

		// break
	}

	for index, path := range Paths.CAN2TSNPath {
		fmt.Printf("\nImportCan Path %d \n", index)
		path.Show_Path()

		// break
	}
}

func (Kpath *KPath) Show_KPath() {
	for index, Path := range Kpath.Paths {
		fmt.Printf("Path %d \n", index)
		fmt.Printf("Path weight: %e \n", Path.Weight)
		Path.Show_Path()
	}
}

func (Path *Path) Show_Path() {
	for _, node := range Path.Nodes {
		// fmt.Println(node.ID)
		for _, c := range node.Connections {
			fmt.Printf("%d --> %d\n", c.FromNodeID, c.ToNodeID)
		}
	}
}

func (path *Path) Show_Cycle() {
	b, cyclelist := path.FindCyCle()
	if b {
		fmt.Println("The MST has cycle")
		fmt.Println(cyclelist)

	} else {
		fmt.Println("The MST has no cycle")
	}
}
