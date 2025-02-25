package path

type KPath_set struct {
	TSNPaths []*KPath
	AVBPaths []*KPath
}

func new_KPath_Set() *KPath_set {
	return &KPath_set{}
}

type KPath struct {
	K      int      
	Source int      
	Target int      
	Paths  []*Path  
}

func new_KPath(k int, source, target int) *KPath{
	return &KPath{
		K:      k,
		Source: source,
		Target: target,
		Paths:  []*Path{},
	}
}

type Path_set struct {
	TSNPaths []*Path
	AVBPaths []*Path
}

func new_Path_Set() *Path_set {
	return &Path_set{}
}

type Path struct {
	Nodes  []*Node
	Weight float64
}

func new_Path() *Path {
	return &Path{}
}

type Node struct {
	ID          int
	Connections []*Connection
}

type Connection struct {
	FromNodeID int     // strat
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
	Vertices map[int]*Vertex 
	Edges    []*Edge         
}

type Vertex struct {
	ID      int      
	Visited bool     
	Cost    float64     
	Path    int      
	Edges   []*Edge 
}

type Edge struct {
	Start int  
	End   int  
	Cost  float64
}