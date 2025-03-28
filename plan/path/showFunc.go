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
		graph.Show_Path()
	}
}

func (graph *Graph) Show_Path() {
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
	tsn := 1
	for _, KPath := range Paths.TSNPaths {
		fmt.Printf("\nTSN Path %d \n", tsn)
		KPath.Show_KPath()
		tsn++
		break
	}
	avb := 1
	for _, KPath := range Paths.AVBPaths {
		fmt.Printf("\nAVB Path %d \n", avb)
		KPath.Show_KPath()
		avb++
		break
	}
}

func (Paths *KPath_Set) Show_Path_Set() {
	tsn := 1
	for _, path := range Paths.TSNPaths {
		fmt.Printf("\nTSN Path %d \n", tsn)
		path.Show_KPath()
		tsn++

		break
	}
	avb := 1
	for _, path := range Paths.AVBPaths {
		fmt.Printf("\nAVB Path %d \n", avb)
		path.Show_KPath()
		avb++
		break
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
		fmt.Println(node.ID)
		for _, c := range node.Connections {
			fmt.Printf("%d --> %d \n", c.FromNodeID, c.ToNodeID)
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
