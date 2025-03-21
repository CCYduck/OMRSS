package network

func (network *Network) Show_Network() {
	network.Topology.Show_Topology()

	network.TSNFlow_Set.Show_TSNFlows()
	network.TSNFlow_Set.Show_TSNFlow()
	network.TSNFlow_Set.Show_TSNStream()

	network.CANFlow_Set.Show_CANFlows()
	network.CANFlow_Set.Show_CANFlow()
	network.CANFlow_Set.Show_CANStream()

	network.TSNGraph_Set.Show_Graphs()
	network.CANGraph_Set.Show_Graphs()
	
}
