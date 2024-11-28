package topology

type Node struct {
	ID          int
	Connections []*Connection
}

type Connection struct {
	FromNodeID int     // strat
	ToNodeID   int     // next
	Cost       float64 // 1Gbps => (750,000 bytes/6ms) 750,000 bytes under 6ms for each link ==> 125 bytes/ns
}

func new_Connection(fromNodeID int, toNodeID int, cost float64) *Connection {
	return &Connection{
		FromNodeID: fromNodeID,
		ToNodeID:   toNodeID,
		Cost:       cost,
	}
}

type Topology struct {
	Talker   []*Node //SOURCE
	Switch   []*Node //TSN switch
	Listener []*Node //destination
	Nodes    []*Node //Not Use node
}

func new_Topology() *Topology {
	return &Topology{}
}

type Edge struct {
	Ends []int `yaml:"ends"`
}

type Data struct {
	Scale struct {
		EndStations int `yaml:"end_stations"`
		Bridges     int `yaml:"bridges"`
	} `yaml:"scale"`
	EndStationEdges []Edge `yaml:"end_station_edges"`
	BridgeEdges     []Edge `yaml:"bridge_edges"`
}
