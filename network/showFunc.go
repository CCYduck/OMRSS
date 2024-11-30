package network

func (network *Network) Show_Network() {
	network.Topology.Show_Topology()

	network.TSNFlow_Set.Show_Flows()
	network.TSNFlow_Set.Show_Flow()
	network.TSNFlow_Set.Show_Stream()

	network.CANFlow_Set.Show_Flows()
	network.CANFlow_Set.Show_Flow()
	network.CANFlow_Set.Show_Stream()
	
	
	network.Graph_Set.Show_Graphs()
}
