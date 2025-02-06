package network

import (
	"fmt"
	"src/network/flow"
	"src/network/graph"
	"src/network/topology"
)

func (network *Network) Generate_Network() {
	// 2. Generate topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	network.Topology = topology.Generate_Topology(network.TopologyName, network.BytesRate)
	fmt.Println("Complete Generating Topology.")
	fmt.Println()

	// 3. Generate flows
	fmt.Println("Generate Flows")
	fmt.Println("----------------------------------------")
	network.TSNFlow_Set = flow.Generate_TSNFlows(len(network.Topology.Nodes), network.BG_TSN, network.BG_AVB, network.Input_TSN, network.Input_AVB, network.HyperPeriod)
	network.CANFlow_Set = flow.Generate_CANFlows(len(network.Topology.Nodes), network.Important_CAN, network.Unimportant_CAN, network.HyperPeriod)
	fmt.Println("Complete Generating Flows.")
	fmt.Println()

	// 4. Simulating graphs using flows in topology
	fmt.Println("Simulating Graphs")
	fmt.Println("----------------------------------------")
	network.Graph_Set = graph.Generate_TSNGraphs(network.Topology, network.TSNFlow_Set, network.BytesRate)
	network.Graph_Set = graph.Generate_CANGraphs(network.Topology, network.CANFlow_Set, network.BytesRate)
	fmt.Println("Complete Simulating Graphs.")
	fmt.Println()
}

func (network *OSRO_Network) Generate_Network() {
	// 2. Generate topology
	fmt.Println("Generate Topology")
	fmt.Println("----------------------------------------")
	network.Topology = topology.Generate_Topology(network.TopologyName, network.BytesRate)
	fmt.Println("Complete Generating Topology.")
	fmt.Println()

	// select CAN node
	CAN_Node_Set := network.Topology.Select_CAN_Node_Set()
	fmt.Printf("CAN nodes: %v", CAN_Node_Set)
	fmt.Println()

}


//func (network *Network) Generate_Network() *Network {
//
//}
