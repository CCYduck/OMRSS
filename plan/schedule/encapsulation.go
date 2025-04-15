package schedule

import (
	"fmt"
	"sort"
	"src/network/flow"
)

func EncapsulateCAN2TSN(f *flow.CANFlows, hyperperiod int) (*CAN2TSN_Flow_Set, int) {
	can2tsnFlowSet := &CAN2TSN_Flow_Set{}
	// create flow set
	for _, impf := range f.ImportantCANFlows {
		can2tsnFlowSet.searchCAN2TSNFlow(impf)
	}
	for _, unimpf := range f.UnimportantCANFlows {
		can2tsnFlowSet.searchCAN2TSNFlow(unimpf)
	}

	// encapsulate
	o1_can_drop := 0
	method := "deadline"
	datasize_max := 100.
	deadline := 5000
	for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
		queue := &Queue{}
		datasize_count := 0.
		for current_time := 0; current_time < hyperperiod; current_time += 5000 {
			queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
			queue.sortQueue(method)

			// example:
			// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
			// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
			for _, queue_stream := range queue.Streams {
				if current_time > queue_stream.FinishTime {
					o1_can_drop += 1

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
	}
	// mq.Show_MQ()

	return can2tsnFlowSet, o1_can_drop

}

type CAN2TSN_Flow_Set struct {
	CAN2TSN_Flows []*CAN2TSN_Flow
}

type CAN2TSN_Flow struct {
	Source       int
	Destination  int
	CAN_Streams  []*flow.Stream
	CAN2TSN_Flow *flow.Flow
}

func (can2tsnFlow *CAN2TSN_Flow) appendCANStreamsToFlow(f *flow.Flow) {
	can2tsnFlow.CAN_Streams = append(can2tsnFlow.CAN_Streams, f.Streams...)
}

func (can2tsnFlow *CAN2TSN_Flow) getStreamsByCurrentTime(current_time int) []*flow.Stream {
	streams := []*flow.Stream{}
	for _, stream := range can2tsnFlow.CAN_Streams {
		if stream.ArrivalTime == current_time {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) addNewCAN2TSNFlowToSet(f *flow.Flow) {
	can2tsnFlow := &CAN2TSN_Flow{}
	can2tsnFlow.Source = f.Source
	can2tsnFlow.Destination = f.Destination
	can2tsnFlow.CAN2TSN_Flow = createCAN2TSNFlow(f.Source, f.Destination)
	can2tsnFlow.appendCANStreamsToFlow(f)
	can2tsnFlowSet.CAN2TSN_Flows = append(can2tsnFlowSet.CAN2TSN_Flows, can2tsnFlow)
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) searchCAN2TSNFlow(f *flow.Flow) {
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

func createCAN2TSNFlow(source int, destination int) *flow.Flow {
	newFlow := &flow.Flow{
		Period:      5000,
		Deadline:    5000,
		DataSize:    100,
		Source:      source,
		Destination: destination,
	}

	return newFlow
}

func createCAN2TSNStream(arrival_time int, deadline int, datasize float64) *flow.Stream {
	newStream := &flow.Stream{
		ArrivalTime: arrival_time,
		Deadline:    deadline,
		DataSize:    datasize,
		FinishTime:  arrival_time + deadline,
	}

	return newStream
}

type Queue struct {
	Streams []*flow.Stream
}

func (q *Queue) appendQueue(streams []*flow.Stream) {
	q.Streams = append(q.Streams, streams...)
}

func (q *Queue) popQueue() {
	q.Streams = q.Streams[1:]
}

func (q *Queue) sortQueue(method string) {
	switch method {
	case "fifo":
		// 到達時間小 → 大
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].ArrivalTime < q.Streams[j].ArrivalTime
		})
	case "deadline":
		// Deadline 小 → 大（最急先送）
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].Deadline < q.Streams[j].Deadline
		})
	default:
		// 預設 FIFO
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].ArrivalTime < q.Streams[j].ArrivalTime
		})
	}
}
