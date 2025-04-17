package flow

import (
	"fmt"
)

func (flow_set *Flows) Generate_TT_Flows(Nnode_length int, bg_tsn int, bg_avb int, input_tsn int, input_avb int, HyperPeriod int) {
	// Generate TT BG Flows, round 1
	Generate_TT_TSNFlow(flow_set, Nnode_length, bg_tsn, HyperPeriod)
	Generate_TT_AVBFlow(flow_set, Nnode_length, bg_avb, HyperPeriod)
	fmt.Printf("Complete generating round%d bgstreams.\n", 1)

	// Generate TT Input Flows,round 2
	Generate_TT_TSNFlow(flow_set, Nnode_length, input_tsn, HyperPeriod)
	Generate_TT_AVBFlow(flow_set, Nnode_length, input_avb, HyperPeriod)
	fmt.Printf("Complete generating round%d tsnstreams.\n", 2)
}

func Generate_TT_TSNFlow(flows *Flows, Nnode_length int, TS int, HyperPeriod int) {
	for flow := 0; flow < TS; flow++ {
		tsn := config_TSN_Stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := random_TT_Devices_For_Path(Nnode_length)

		Flow := Generate_TT_Streams(tsn.Period, tsn.Deadline, tsn.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.TSNFlows = append(flows.TSNFlows, Flow)
	}
}

func Generate_TT_AVBFlow(flows *Flows, Nnode_length int, AS int, HyperPeriod int) {
	for flow := 0; flow < AS; flow++ {
		avb := config_AVB_Stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := random_TT_Devices_For_Path(Nnode_length)

		Flow := Generate_TT_Streams(avb.Period, avb.Deadline, avb.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.AVBFlows = append(flows.AVBFlows, Flow)
	}
}

func Generate_TT_Streams(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var number int = 0

	flow := new_TTFlow(period, deadline, datasize, HyperPeriod)
	for ArrivalTime := 0; ArrivalTime < HyperPeriod; ArrivalTime += period {
		FinishTime := ArrivalTime + deadline
		name := fmt.Sprint("ttstream", number)
		stream := new_TTStream(name, ArrivalTime, datasize, deadline, FinishTime)
		flow.Streams = append(flow.Streams, stream)
		number += 1
	}

	return flow
}
