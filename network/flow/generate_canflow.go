package flow

import (
	"fmt"
)

func Generate_CANFlows(Nnode int,importantCAN int, unimportantCAN int, HyperPeriod int) *CANFlows {
	// Constructing Flows structures

	canflow_set := new_CANFlows()
	//importantCAN
	Generate_importantCANFlow(canflow_set, Nnode, importantCAN, HyperPeriod)
	//unimportantCAN
	Generate_unimportantCANFlow(canflow_set, Nnode, unimportantCAN, HyperPeriod)
	fmt.Printf("Complete generating round%d streams.\n", 2)

	return canflow_set
}



func Generate_importantCANFlow(flows *CANFlows, Nnode int, impcan int, HyperPeriod int) {
	for flow := 0; flow < impcan; flow++ {
		importantCAN := importantCAN_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destinations := Random_Devices(Nnode)

		Flow := Generate_stream(importantCAN.Period, importantCAN.Deadline, importantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows.importantCANFlows = append(flows.importantCANFlows, Flow)
	}
}

func Generate_unimportantCANFlow(flows *CANFlows, Nnode int, umimpcan int, HyperPeriod int) {
	for flow := 0; flow < umimpcan; flow++ {
		unimportantCAN := unimportantCAN_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destinations := Random_Devices(Nnode)

		Flow := Generate_stream(unimportantCAN.Period, unimportantCAN.Deadline, unimportantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destinations = destinations

		flows.unimportantCANFlows = append(flows.unimportantCANFlows, Flow)
	}
}

func Generate_CANstream(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var (
		ArrivalTime int = 0
		FinishTime  int = 0
		Deadline    int = 0
		number      int = 0
	)

	flow := new_Flow(period, deadline, datasize, HyperPeriod)

	for FinishTime < HyperPeriod {
		Deadline += deadline
		FinishTime += period
		name := fmt.Sprint("stream", number)

		stream := new_Stream(name, ArrivalTime, datasize, deadline, FinishTime)

		flow.Streams = append(flow.Streams, stream)
		ArrivalTime += period
		number += 1
	}

	return flow
}
