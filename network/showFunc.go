package network

func (network *Network) Show_Network() {
	network.Topology.Show_Topology()

	network.TSNFlow_Set.Show_TSNFlows()
	network.TSNFlow_Set.Show_TSNFlow()
	network.TSNFlow_Set.Show_TSNStream()

	network.Graph_Set.Show_Graphs()
}

func (network *OSRO_Network) Show_Network() {
	network.Topology.Show_Topology()

	network.TSNFlow_Set.Show_TSNFlows()
	network.TSNFlow_Set.Show_TSNFlow()
	network.TSNFlow_Set.Show_TSNStream()

	network.CANFlow_Set.Show_CANFlows()
	network.CANFlow_Set.Show_CANFlow()
	network.CANFlow_Set.Show_CANStream()

	network.Graph_Set.Show_Graphs()
	
}
