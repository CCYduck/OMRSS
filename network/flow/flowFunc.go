package flow

func (flows *TSNFlows) Input_flow_set() *TSNFlows {
	Input_flow_set := new_TSNFlows()

	Input_flow_set.TSNFlows = append(Input_flow_set.TSNFlows, flows.TSNFlows[bg_tsnflows_end:]...)
	Input_flow_set.AVBFlows = append(Input_flow_set.AVBFlows, flows.AVBFlows[bg_avbflow_end:]...)

	return Input_flow_set
}
func (flows *CANFlows) Input_flow_set() *CANFlows {
	Input_flow_set := new_CANFlows()

	Input_flow_set.importantCANFlows = append(Input_flow_set.importantCANFlows, flows.importantCANFlows...)
	Input_flow_set.unimportantCANFlows = append(Input_flow_set.unimportantCANFlows, flows.unimportantCANFlows...)

	return Input_flow_set
}

func (flows *TSNFlows) BG_flow_set() *TSNFlows {
	BG_flow_set := new_TSNFlows()

	BG_flow_set.TSNFlows = append(BG_flow_set.TSNFlows, flows.TSNFlows[:bg_tsnflows_end]...)
	BG_flow_set.AVBFlows = append(BG_flow_set.AVBFlows, flows.AVBFlows[:bg_avbflow_end]...)

	return BG_flow_set
}
