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
	fmt.Printf("Complete generating round%d canstreams.\n", 2)

	return canflow_set
}



func Generate_importantCANFlow(flows *CANFlows, Nnode int, impcan int, HyperPeriod int) {
	for flow := 0; flow < impcan; flow++ {
		importantCAN := importantCAN_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener) 
		// 這裡要修改Random方法
		source, destination := Random_Devices(Nnode)
		
		// 這裡要加上封裝的延遲
		Flow := Generate_CANstream(importantCAN.Period, importantCAN.Deadline, importantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.ImportantCANFlows = append(flows.ImportantCANFlows, Flow)
	}
}

func Generate_unimportantCANFlow(flows *CANFlows, Nnode int, umimpcan int, HyperPeriod int) {
	for flow := 0; flow < umimpcan; flow++ {
		unimportantCAN := unimportantCAN_stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := Random_Devices(Nnode)

		Flow := Generate_CANstream(unimportantCAN.Period, unimportantCAN.Deadline, unimportantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.UnimportantCANFlows = append(flows.UnimportantCANFlows, Flow)
	}
}

func Generate_CANstream(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var (
		ArrivalTime int = 0
		FinishTime  int = 0
		Deadline    int = 0
		number      int = 0
	)
	//這裡要修改
	flow := new_CANFlow(period, deadline, datasize, HyperPeriod)

	for FinishTime < HyperPeriod {
		Deadline += deadline
		FinishTime += period
		name := fmt.Sprint("canstream", number)

		stream := new_CANStream(name, ArrivalTime, datasize, deadline, FinishTime)
		flow.Streams = append(flow.Streams, stream)
		ArrivalTime += period
		number += 1
	}

	return flow
}