package network

func (network *Network) Show_Network() {
	network.Topology.Show_Topology()

	network.Flow_Set.Show_TSNFlows()
	network.Flow_Set.Show_TSNFlow()
	network.Flow_Set.Show_TSNStream()

	network.Flow_Set.Show_CANFlows()
	network.Flow_Set.Show_CANFlow()
	network.Flow_Set.Show_CANStream()

	network.Graph_Set.Show_Graphs()
	network.Graph_Set.Show_Graphs()

}
