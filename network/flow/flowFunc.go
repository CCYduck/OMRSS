package flow

func (flows *Flows) Input_flow_set() *Flows {
	Input_flow_set := new_Flows()
	Can_flow_set :=new_CANFlows()

	Input_flow_set.TSNFlows = append(Input_flow_set.TSNFlows, flows.TSNFlows[bg_tsnflows_end:]...)
	Input_flow_set.AVBFlows = append(Input_flow_set.AVBFlows, flows.AVBFlows[bg_avbflow_end:]...)

	Can_flow_set.importantCANFlows=append(Can_flow_set.importantCANFlows, flows.importantCANFlows...)
	Can_flow_set.unimportantCANFlows=append(Can_flow_set.unimportantCANFlows, flows.unimportantCANFlows...)

	return Input_flow_set
}

func (flows *Flows) BG_flow_set() *Flows {
	BG_flow_set := new_Flows()

	BG_flow_set.TSNFlows = append(BG_flow_set.TSNFlows, flows.TSNFlows[:bg_tsnflows_end]...)
	BG_flow_set.AVBFlows = append(BG_flow_set.AVBFlows, flows.AVBFlows[:bg_avbflow_end]...)

	return BG_flow_set
}
