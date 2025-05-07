package path

type KPath_Set struct {
	TSNPaths 		[]*KPath
	AVBPaths 		[]*KPath
	CAN2TSNPaths	[]*KPath
}

func new_KPath_Set() *KPath_Set {
	return &KPath_Set{}
}

type KPath struct {
	Method string
	K      int
	Source int
	Target int
	Paths  []*Path
}

func new_KPath(k int, source, target int) *KPath {
	return &KPath{
		K:      k,
		Source: source,
		Target: target,
		Paths:  []*Path{},
	}
}

type Path_set struct {
	TSNPath     []*Path
	AVBPath     []*Path
	CAN2TSNPath []*Path
}

func new_Path_Set() *Path_set {
	return &Path_set{}
}

type Path struct {
	Method 	string
	IDs 	[]int
	Nodes  	[]*Node
	Weight 	float64
}

func new_Path() *Path {
	return &Path{}
}

type Node struct {
	ID          int
	Shape       string
	Connections []*Connection
}

type Connection struct {
	FromNodeID int     // start
	ToNodeID   int     // next
	Cost       float64 // 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link ==> 125 bytes/us
}

func new_Connection(fromNodeID int, toNodeID int, cost float64) *Connection {
	return &Connection{
		FromNodeID: fromNodeID,
		ToNodeID:   toNodeID,
		Cost:       cost,
	}
}

type V2V struct {
	V2VEdges []*V2VEdge
}

type V2VEdge struct {
	FromVertex int
	Graphs     []*Graph
}

type Graph struct {
	Vertexs  	[]*Vertex
	V 			map[int]*Vertex
	ToVertex 	int
	Path     	[][]int
	ConnPath 	[][]Connection
}

type Vertex struct {
	ID      int
	Visited bool
	Cost    int
	Path    int
	Edges   []*Edge
}

type Edge struct {
	Strat int
	End   int
	Cost  int
}
