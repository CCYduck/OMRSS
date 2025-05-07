package flow

// import "fmt"

func (flows *Flows) Input_TSNflow_set() *Flows {
	Input_flow_set := new_Flows()

	Input_flow_set.TSNFlows = append(Input_flow_set.TSNFlows, flows.TSNFlows[bg_tsnflows_end:]...)
	Input_flow_set.AVBFlows = append(Input_flow_set.AVBFlows, flows.AVBFlows[bg_avbflows_end:]...)
	Input_flow_set.Encapsulate =append(Input_flow_set.Encapsulate,flows.Encapsulate... )

	return Input_flow_set
}


func (flows *Flows) BG_flow_set() *Flows {
	BG_flow_set := new_Flows()

	BG_flow_set.TSNFlows = append(BG_flow_set.TSNFlows, flows.TSNFlows[:bg_tsnflows_end]...)
	BG_flow_set.AVBFlows = append(BG_flow_set.AVBFlows, flows.AVBFlows[:bg_avbflows_end]...)

	return BG_flow_set
}

func (flows *Flows) FindMethod(methodname string) []*Flow{
	for _,method :=range flows.Encapsulate{
		if method.Method_Name ==methodname {
			// fmt.Println(len(method.CAN2TSNFlows))
			return method.CAN2TSNFlows
		}
	}
	return nil
}