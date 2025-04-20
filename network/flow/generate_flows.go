package flow

var (
	bg_tsnflows_end int
	bg_avbflows_end int
)

func Generate_OSRO_Flows(CANnode []int, importantCAN int, unimportantCAN int, Nnode_length int, bg_tsn int, bg_avb int, input_tsn int, input_avb int, HyperPeriod int) *Flows {
	flow_set := new_Flows()
	bg_tsnflows_end = bg_tsn
	bg_avbflows_end = bg_avb

	flow_set.Generate_TT_Flows(Nnode_length, bg_tsn, bg_avb, input_tsn, input_avb, HyperPeriod)
	flow_set.Generate_CAN2TSN_Flows(CANnode, importantCAN, unimportantCAN, HyperPeriod, "FIFO")
	flow_set.Show_CANFlows()

	return flow_set
}
