package flow

import (
	"fmt"
)

var (
	bg_tsnflows_end int
	bg_avbflows_end  int
)

func Generate_TSNFlows(Nnode int, bg_tsn int, bg_avb int, input_tsn int, input_avb int, HyperPeriod int) *TSNFlows {
	// Constructing Flows structures
	flow_set := new_TSNFlows()
	bg_tsnflows_end = bg_tsn
	bg_avbflows_end = bg_avb

	// round 1
	Generate_TSNFlow(flow_set, Nnode, bg_tsn, HyperPeriod)
	Generate_AVBFlow(flow_set, Nnode, bg_avb, HyperPeriod)
	fmt.Printf("Complete generating round%d bgstreams.\n", 1)

	// round 2
	Generate_TSNFlow(flow_set, Nnode, input_tsn, HyperPeriod)
	Generate_AVBFlow(flow_set, Nnode, input_avb, HyperPeriod)
	fmt.Printf("Complete generating round%d tsnstreams.\n", 2)

	return flow_set
}

func Generate_TSNFlow(flows *TSNFlows, Nnode int, TS int, HyperPeriod int) {
	for flow := 0; flow < TS; flow++ {
		tsn := TSN_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := Random_Devices(Nnode)

		Flow := Generate_TSNstream(tsn.Period, tsn.Deadline, tsn.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.TSNFlows = append(flows.TSNFlows, Flow)
	}
}

func Generate_AVBFlow(flows *TSNFlows, Nnode int, AS int, HyperPeriod int) {
	for flow := 0; flow < AS; flow++ {
		avb := AVB_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := Random_Devices(Nnode)

		Flow := Generate_TSNstream(avb.Period, avb.Deadline, avb.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.AVBFlows = append(flows.AVBFlows, Flow)
	}
}

func Generate_TSNstream(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var (
		ArrivalTime int = 0
		FinishTime  int = 0
		Deadline    int = 0
		number      int = 0
	)

	flow := new_TSNFlow(period, deadline, datasize, HyperPeriod)

	for FinishTime < HyperPeriod {
		Deadline += deadline
		FinishTime += period
		name := fmt.Sprint("tsnstream", number)

		stream := new_TSNStream(name, ArrivalTime, datasize, deadline, FinishTime)
		flow.Streams = append(flow.Streams, stream)
		ArrivalTime += period
		number += 1
	}

	return flow
}


