package schedule

import (
	"fmt"
	"sort"
	"src/network/flow"
)

func EncapsulateCAN2TSN(f *flow.Flows, hyperperiod int) (*CAN2TSN_Flow_Set, int , int) {
	can2tsnFlowSet := &CAN2TSN_Flow_Set{}
	// create flow set
	for _, impf := range f.ImportantCANFlows {
		can2tsnFlowSet.searchCAN2TSNFlow(impf)
	}
	for _, unimpf := range f.UnimportantCANFlows {
		can2tsnFlowSet.searchCAN2TSNFlow(unimpf)
	}

	// encapsulate
	o1_can_drop_1	:= 0
	o1_can_drop_2	:= 0
	method			:= "deadline"
	datasize_max 	:= 100.
	deadline 		:= 5000
	step 			:= 5000

	for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
		queue := &Queue{}
		
		for current_time := 0; current_time < hyperperiod; current_time += step {
			
			datasize_count := 0.
			queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
			queue.sortQueue(method)
			
			// M1
			// example:
			// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
			// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
			for _, queue_stream := range queue.Streams {
				
				if current_time > queue_stream.FinishTime {
					o1_can_drop_1 += 1
					queue.popQueue()
					continue
				} 
				datasize_count += queue_stream.DataSize

				if datasize_count >= datasize_max {
					can2tsnStream := createCAN2TSNStream(current_time, deadline, datasize_max)
					can2tsnFlow.CAN2TSN_Flow.Streams = append(can2tsnFlow.CAN2TSN_Flow.Streams, can2tsnStream)
					datasize_count = 0
					break
				}

				// if queue_stream.FinishTime - current_time < 5000{
				// 	can2tsnStream := createCAN2TSNStream(current_time, deadline, datasize_max)
				// 	can2tsnFlow.CAN2TSN_Flow.Streams = append(can2tsnFlow.CAN2TSN_Flow.Streams, can2tsnStream)
				// 	datasize_count = 0
				// 	break
				// }
				queue.popQueue()
					
			}	
		}

		queue_1 := &Queue{}
		for current_time := 0; current_time < hyperperiod; current_time += step {
			datasize_count := 0.
			
			queue_1.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
			queue_1.sortQueue(method)
			// M2
			// ❗ 不用 range；自己控制 pop，才不會亂位
			for len(queue_1.Streams) > 0 {
				qs :=queue_1.Streams[0]      // 先抓頭
				queue_1.popQueue()            // 立刻移除，避免重複

				if current_time >= qs.FinishTime {   // 已逾期
					o1_can_drop_2++
					continue
				}

				datasize_count += qs.DataSize
				timeLeft := qs.FinishTime - current_time

				// 決定是否提前封裝
				shouldSend := datasize_count >= datasize_max || // 滿載
				timeLeft < step  ||                       // 快逾期
				len(queue_1.Streams) == 0                       // 佇列已空

				
				if shouldSend {
					size := datasize_count    // 用實際累積量，而不是 datasizeMax
					can2tsnStream := createCAN2TSNStream(current_time, deadline, size)
					can2tsnFlow.CAN2TSN_Flow.Streams =
						append(can2tsnFlow.CAN2TSN_Flow.Streams, can2tsnStream)

					datasize_count = 0        // 重置計數
					break                    // 本時槽已送一包，下一包留到 5 ms 後
				}
			}
		}
		

	}
	// mq.Show_MQ()
	
	return can2tsnFlowSet, o1_can_drop_1 , o1_can_drop_2

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


