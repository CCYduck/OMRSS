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
	fmt.Println("Complete Generating Flows.")
	fmt.Println()

	// 4. Simulating graphs using flows in topology
	fmt.Println("Simulating Graphs")
	fmt.Println("----------------------------------------")
	network.Graph_Set = graph.Generate_Graphs(network.Topology, network.TSNFlow_Set, network.BytesRate)
	fmt.Println("Complete Simulating Graphs.")
	fmt.Println()
}

func new_Network(topology_name string, bg_tsn int, bg_avb int, input_tsn int, input_avb int, importantCAN int, unimportantCAN int, hyperperiod int, bandwidth float64) *Network {
	bw := (bandwidth / 8) * 1e-6 // bytes/us ==> 125 bytes
	bytes_rate := 1. / bw        // The number of bytes that can be transmitted in 1us ==> 1/125
	bw *= float64(hyperperiod)   // The bytes that can be transmitted in 6000us (bytes/us * hyperperiod) ==> 750000 bytes

	Network := &Network{
		HyperPeriod:    hyperperiod,
		BytesRate:      bytes_rate,
		Bandwidth:      bw,
		TopologyName:   topology_name,
		BG_TSN:         bg_tsn,
		BG_AVB:         bg_avb,
		Input_TSN:      input_tsn,
		Input_AVB:      input_avb,
		ImportantCAN:   importantCAN,
		UnimportantCAN: unimportantCAN,
	}

	fmt.Println("Define network parameters")
	fmt.Println("----------------------------------------")
	fmt.Printf("HyperPeriod: %d us \n", Network.HyperPeriod)
	fmt.Printf("Bandwidth:  %f bytes/6000us \n", Network.Bandwidth)
	fmt.Printf("BytesRate:  %f BPS \n", Network.BytesRate)
	fmt.Printf("Topology file name: %s \n", Network.TopologyName)
	fmt.Printf("TSN flow: %d, AVB flow: %d \n", Network.Input_TSN+Network.BG_TSN, Network.Input_AVB+Network.BG_AVB)
	fmt.Printf("ImportantCAN flow: %d, UnimportantCAN flow: %d \n", Network.ImportantCAN, Network.UnimportantCAN)
	fmt.Println()

	return Network
}
