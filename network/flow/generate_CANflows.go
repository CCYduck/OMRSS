package flow

import (
	"fmt"
)

func (flow_set *Flows) Generate_CAN_Flows(CANnode []int, importantCAN int, unimportantCAN int, HyperPeriod int) {
	// Generate CAN Flows
	Generate_Important_CANFlow(flow_set, CANnode, importantCAN, HyperPeriod)
	Generate_Unimportant_CANFlow(flow_set, CANnode, unimportantCAN, HyperPeriod)
	fmt.Println("Complete generating can streams.")
}

func Generate_Important_CANFlow(flows *Flows, CANnode []int, impcan int, HyperPeriod int) {
	for flow := 0; flow < impcan; flow++ {
		importantCAN := config_ImportantCAN_Stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := random_CAN_Devices_For_Path(CANnode)

		Flow := Generate_CAN_Streams(importantCAN.Period, importantCAN.Deadline, importantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.ImportantCANFlows = append(flows.ImportantCANFlows, Flow)
	}
}

func Generate_Unimportant_CANFlow(flows *Flows, CANnode []int, umimpcan int, HyperPeriod int) {
	for flow := 0; flow < umimpcan; flow++ {
		unimportantCAN := config_UnimportantCAN_Stream()

		// Random End Devices 1. source(Talker) 2. destinations(listener)
		source, destination := random_CAN_Devices_For_Path(CANnode)

		Flow := Generate_CAN_Streams(unimportantCAN.Period, unimportantCAN.Deadline, unimportantCAN.DataSize, HyperPeriod)
		Flow.Source = source
		Flow.Destination = destination

		flows.UnimportantCANFlows = append(flows.UnimportantCANFlows, Flow)
	}
}

func Generate_CAN_Streams(period int, deadline int, datasize float64, HyperPeriod int) *Flow {
	var number int = 0

	flow := new_CANFlow(period, deadline, datasize, HyperPeriod)
	for ArrivalTime := 0; ArrivalTime < HyperPeriod; ArrivalTime += period {
		FinishTime := ArrivalTime + deadline
		name := fmt.Sprint("canstream", number)
		stream := new_CANStream(name, ArrivalTime, datasize, deadline, FinishTime)
		flow.Streams = append(flow.Streams, stream)
		number += 1
	}

	return flow
}
