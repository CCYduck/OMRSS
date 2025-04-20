package flow

import (
	"fmt"
	"sort"
	"time"
)

func (flow_set *Flows) Generate_CAN2TSN_Flows(CANnode []int, importantCAN int, unimportantCAN int, hyperperiod int, method string) {
	ImportantCANFlows, UnimportantCANFlows := Generate_CAN_Flows(CANnode, importantCAN, unimportantCAN, hyperperiod)
	//fmt.Println(ImportantCANFlows)
	// create flow set
	can2tsnFlowSet := &CAN2TSN_Flow_Set{}
	for _, impf := range ImportantCANFlows {
		can2tsnFlowSet.searchCAN2TSNFlow(impf)
	}
	for _, unimpf := range UnimportantCANFlows {
		can2tsnFlowSet.searchCAN2TSNFlow(unimpf)
	}

	// encapsulate
	start_time := time.Now()
	can2tsnFlowSet.EncapsulateCAN2TSN(hyperperiod, method)
	can2tsn_delay := time.Since(start_time)

	flow_set.CAN2TSN_O1_Drop = can2tsnFlowSet.O1_Drop
	flow_set.CAN2TSN_Delay = can2tsn_delay
	for _, can2tsn_flow := range can2tsnFlowSet.CAN2TSN_Flows {
		flow_set.CAN2TSNFlows = append(flow_set.CAN2TSNFlows, can2tsn_flow.CAN2TSN_Flow)
	}
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) EncapsulateCAN2TSN(hyperperiod int, method string) {
	datasize_max := 100.
	deadline := 5000
	for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
		queue := &Queue{}
		datasize_count := 0.
		for current_time := 0; current_time < hyperperiod; current_time += 5000 {
			queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
			queue.sortQueue(method, current_time)

			// example:
			// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
			// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
			for _, queue_stream := range queue.Streams {
				if current_time > queue_stream.FinishTime {
					can2tsnFlowSet.O1_Drop += 1
					queue.popQueue()

				} else {
					datasize_count += queue_stream.DataSize
					if datasize_count >= datasize_max {
						can2tsnStream := createCAN2TSNStream(current_time, deadline, datasize_max)
						can2tsnFlow.CAN2TSN_Flow.Streams = append(can2tsnFlow.CAN2TSN_Flow.Streams, can2tsnStream)
						datasize_count = 0
						break
					}
					queue.popQueue()
				}
			}
		}
		last_index := len(can2tsnFlow.CAN2TSN_Flow.Streams) - 1
		if last_index >= 0 {
			if can2tsnFlow.CAN2TSN_Flow.Streams[last_index].FinishTime < hyperperiod-5000 {
				if datasize_count > 0 {
					can2tsnStream := createCAN2TSNStream(can2tsnFlow.CAN2TSN_Flow.Streams[last_index].FinishTime, deadline, datasize_max)
					can2tsnFlow.CAN2TSN_Flow.Streams = append(can2tsnFlow.CAN2TSN_Flow.Streams, can2tsnStream)
					datasize_count = 0
				}

			}
		}
	}

}

type CAN2TSN_Flow_Set struct {
	CAN2TSN_Flows []*CAN2TSN_Flow
	O1_Drop       int
}

type CAN2TSN_Flow struct {
	Source       int
	Destination  int
	CAN_Streams  []*Stream
	CAN2TSN_Flow *Flow
}

func (can2tsnFlow *CAN2TSN_Flow) appendCANStreamsToFlow(f *Flow) {
	can2tsnFlow.CAN_Streams = append(can2tsnFlow.CAN_Streams, f.Streams...)
}

func (can2tsnFlow *CAN2TSN_Flow) getStreamsByCurrentTime(current_time int) []*Stream {
	streams := []*Stream{}
	for _, stream := range can2tsnFlow.CAN_Streams {
		if stream.ArrivalTime == current_time {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) addNewCAN2TSNFlowToSet(f *Flow) {
	can2tsnFlow := &CAN2TSN_Flow{}
	can2tsnFlow.Source = f.Source
	can2tsnFlow.Destination = f.Destination
	can2tsnFlow.CAN2TSN_Flow = createCAN2TSNFlow(f.Source, f.Destination)
	can2tsnFlow.appendCANStreamsToFlow(f)
	can2tsnFlowSet.CAN2TSN_Flows = append(can2tsnFlowSet.CAN2TSN_Flows, can2tsnFlow)
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) searchCAN2TSNFlow(f *Flow) {
	s := f.Source
	d := f.Destination

	for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
		if can2tsnFlow.Source == s && can2tsnFlow.Destination == d {
			can2tsnFlow.appendCANStreamsToFlow(f)
			return
		}
	}
	can2tsnFlowSet.addNewCAN2TSNFlowToSet(f)
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) Show_CAN2TSNFlowSet() {
	for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
		fmt.Printf("Queue (%d→%d) streams=%d\n", can2tsnFlow.Source, can2tsnFlow.Destination, len(can2tsnFlow.CAN_Streams))
		for ind, stream := range can2tsnFlow.CAN2TSN_Flow.Streams {
			fmt.Printf("Stream %v ,ArrivalTime: %v  ,Deadline: %v ,Datasize: %v\n", ind, stream.ArrivalTime, stream.Deadline, stream.DataSize)
		}

	}
}

func createCAN2TSNFlow(source int, destination int) *Flow {
	newFlow := &Flow{
		Period:      5000,
		Deadline:    5000,
		DataSize:    100,
		Source:      source,
		Destination: destination,
	}

	return newFlow
}

func createCAN2TSNStream(arrival_time int, deadline int, datasize float64) *Stream {
	newStream := &Stream{
		ArrivalTime: arrival_time,
		Deadline:    deadline,
		DataSize:    datasize,
		FinishTime:  arrival_time + deadline,
	}

	return newStream
}

type Queue struct {
	Streams []*Stream
}

func (q *Queue) appendQueue(streams []*Stream) {
	q.Streams = append(q.Streams, streams...)
}

func (q *Queue) popQueue() {
	q.Streams = q.Streams[1:]
}

func (q *Queue) sortQueue(method string, current_time int) {
	switch method {
	case "fifo":
		// 到達時間小 → 大
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].ArrivalTime < q.Streams[j].ArrivalTime
		})

	case "priority":
		// Deadline 小 → 大（最急先送）
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].Deadline < q.Streams[j].Deadline
		})

	case "mat":
		// MAT 根據剩餘時間
		sort.Slice(q.Streams, func(i, j int) bool {
			ti := q.Streams[i].FinishTime - current_time
			tj := q.Streams[j].FinishTime - current_time
			return ti < tj
		})

	default:
		// 預設 FIFO
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].ArrivalTime < q.Streams[j].ArrivalTime
		})
	}
}
