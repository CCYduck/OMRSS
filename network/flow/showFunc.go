package flow

import (
	"fmt"
)

func (flows *Flows) Show_TSNStream() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNflow", number)
		fmt.Println(name)
		for _, stream := range flow.Streams {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
		}
		number += 1
		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBflow", number)
		fmt.Println(name)
		for _, stream := range flow.Streams {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
		}
		number += 1
		break
	}
}

func (flows *Flows) Show_CANStream() {
	Encapsulate := flows.Encapsulate
	number := 1
	for _, flow := range Encapsulate {
		for _, method := range flow.CAN2TSNFlows{
			fmt.Printf("%v ImportantCANflow %v",method, number)
			for _, stream := range method.Streams {
				fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
			}
		}
		number += 1
		break
	}
}

func (flows *Flows) Show_TSNFlow() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows

	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNFlow", number)
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destination)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBFlow", number)
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destination)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}
}

func (flows *Flows) Show_CANFlow() {
	Encapsulate := flows.Encapsulate

	number := 1
	for _, flow := range Encapsulate {
		for _, method := range flow.CAN2TSNFlows{
			name := fmt.Sprint("ImportantCANFlow", number)
			fmt.Printf("Source: %d\n", method.Source)
			fmt.Printf("Destinations: %v\n", method.Destination)
			fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, method.Period, method.Deadline, method.DataSize)
			number += 1

		break
		}
		
	}
}

func (flows *Flows) Show_TSNFlows() {
	// Display all flows.
	fmt.Printf("Total Flows:%d ( TSN Flows:%d  AVB Flows:%d )\n",
		len(flows.TSNFlows)+len(flows.AVBFlows), len(flows.TSNFlows), len(flows.AVBFlows))
}

func (flows *Flows) Show_CANFlows() {
	// Display all flows.
	for _, flow := range flows.Encapsulate {
		fmt.Printf(" %s: Total CAN2TSN Flows: %d StreamSize: %v StreamCount: %v CAN2TSN_O1_Drop: %d CAN_Area_O1_Drop: %d Delay: %v \n", 
		flow.Method_Name, len(flow.CAN2TSNFlows) , flow.BytesSent, flow.TSNFrameCount , flow.CAN2TSN_O1_Drop, flow.CAN_Area_O1_Drop, flow.CAN2TSN_Delay)
	}
}
