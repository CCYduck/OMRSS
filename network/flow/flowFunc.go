package flow

func (flows *TSNFlows) Input_TSNflow_set() *TSNFlows {
	Input_flow_set := new_Flows()

	Input_flow_set.TSNFlows = append(Input_flow_set.TSNFlows, flows.TSNFlows[bg_tsnflows_end:]...)
	Input_flow_set.AVBFlows = append(Input_flow_set.AVBFlows, flows.AVBFlows[bg_avbflow_end:]...)

	return Input_flow_set
}

func (flows *TSNFlows) BG_flow_set() *TSNFlows {
	BG_flow_set := new_Flows()

	BG_flow_set.TSNFlows = append(BG_flow_set.TSNFlows, flows.TSNFlows[:bg_tsnflows_end]...)
	BG_flow_set.AVBFlows = append(BG_flow_set.AVBFlows, flows.AVBFlows[:bg_avbflow_end]...)

	return BG_flow_set
}

func (flows *CANFlows) Input_CANflow_set() *CANFlows {
	Can_flow_set :=new_CANFlows()

	Can_flow_set.importantCANFlows=append(Can_flow_set.importantCANFlows, flows.importantCANFlows...)
	Can_flow_set.unimportantCANFlows=append(Can_flow_set.unimportantCANFlows, flows.unimportantCANFlows...)

	return Can_flow_set
}
