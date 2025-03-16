package flow

func (flows *TSNFlows) Input_TSNflow_set() *TSNFlows {
	Input_flow_set := new_TSNFlows()

	Input_flow_set.TSNFlows = append(Input_flow_set.TSNFlows, flows.TSNFlows[bg_tsnflows_end:]...)
	Input_flow_set.AVBFlows = append(Input_flow_set.AVBFlows, flows.AVBFlows[bg_avbflows_end:]...)

	return Input_flow_set
}
func (flows *CANFlows) Input_CANflow_set() *CANFlows {
	Input_flow_set := new_CANFlows()

	Input_flow_set.ImportantCANFlows = append(Input_flow_set.ImportantCANFlows, flows.ImportantCANFlows...)
	Input_flow_set.UnimportantCANFlows = append(Input_flow_set.UnimportantCANFlows, flows.UnimportantCANFlows...)

	return Input_flow_set
}

func (flows *TSNFlows) BG_flow_set() *TSNFlows {
	BG_flow_set := new_TSNFlows()

	BG_flow_set.TSNFlows = append(BG_flow_set.TSNFlows, flows.TSNFlows[:bg_tsnflows_end]...)
	BG_flow_set.AVBFlows = append(BG_flow_set.AVBFlows, flows.AVBFlows[:bg_avbflows_end]...)

	return BG_flow_set
}
