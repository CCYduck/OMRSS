package flow

import (
	"fmt"
)

func (flows *TSNFlows) Show_TSNStream() {
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

func (flows *CANFlows) Show_CANStream() {
	importantCANFlows := flows.importantCANFlows
	unimportantCANFlows := flows.unimportantCANFlows
	number := 1
	for _, flow := range importantCANFlows {
		name := fmt.Sprint("ImportantCANflow", number)
		fmt.Println(name)
		for _, stream := range flow.Streams {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
		}
		number += 1
		break
	}

	number = 1
	for _, flow := range unimportantCANFlows {
		name := fmt.Sprint("UnimportantCANflow", number)
		fmt.Println(name)
		for _, stream := range flow.Streams {
			fmt.Printf("%s ArrivalTime:%d DataSize:%f Deadline:%d FinishTime:%d\n",
				stream.Name, stream.ArrivalTime, stream.DataSize, stream.Deadline, stream.FinishTime)
		}
		number += 1
		break
	}
}

func (flows *TSNFlows) Show_TSNFlow() {
	TSNFlows := flows.TSNFlows
	AVBFlows := flows.AVBFlows
	
	number := 1
	for _, flow := range TSNFlows {
		name := fmt.Sprint("TSNFlow", number)
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}

	number = 1
	for _, flow := range AVBFlows {
		name := fmt.Sprint("AVBFlow", number)
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}
}

func (flows *CANFlows) Show_CANFlow() {
	importantCANFlows := flows.importantCANFlows
	unimportantCANFlows := flows.unimportantCANFlows
	
	number := 1
	for _, flow := range importantCANFlows {
		name := fmt.Sprint("ImportantCANFlow", number)
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}

	number = 1
	for _, flow := range unimportantCANFlows {
		name := fmt.Sprint("UnimportantCANFlow", number)
		fmt.Printf("Source: %d\n", flow.Source)
		fmt.Printf("Destinations: %v\n", flow.Destinations)
		fmt.Printf("%s : period:%d us, deadline:%d us, datasize:%f bytes\n",
			name, flow.Period, flow.Deadline, flow.DataSize)
		number += 1

		break
	}
}

func (flows *TSNFlows) Show_TSNFlows() {
	// Display all flows.
	fmt.Printf("Total Flows:%d ( TSN Flows:%d  AVB Flows:%d )\n",
		len(flows.TSNFlows)+len(flows.AVBFlows), len(flows.TSNFlows), len(flows.AVBFlows))
}

func (flows *CANFlows) Show_CANFlows() {
	// Display all flows.
	fmt.Printf("Total Flows:%d (  ImportantCAN Flows:%d   UnimportantCAN Flows:%d )\n",
		len(flows. importantCANFlows)+len(flows. unimportantCANFlows), len(flows. importantCANFlows), len(flows. unimportantCANFlows))
}
