package flow

import (
	"fmt"
	"sort"
	"time"
)

func (flow_set *Flows) Generate_CAN2TSN_Flows(CANnode []int, importantCAN int, unimportantCAN int, hyperperiod int) {
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
	method_name := []string{"fifo","priority", "obo", "wat"}
	for _,name := range method_name {
		method_struct:=&Method{}

		start_time := time.Now()
		can2tsnFlowSet.EncapsulateCAN2TSN(hyperperiod, name)
		method_delay := time.Since(start_time)

		method_struct.Method_Name=name
		method_struct.CAN2TSN_Delay = method_delay
		method_struct.CAN2TSN_O1_Drop = can2tsnFlowSet.O1_Drop

		for _, can2tsn_flow := range can2tsnFlowSet.CAN2TSN_Flows {
			method_struct.CAN2TSNFlows = append(method_struct.CAN2TSNFlows, can2tsn_flow.CAN2TSN_Flow)
		}
		flow_set.Encapsulate = append(flow_set.Encapsulate, method_struct)
	}
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) EncapsulateCAN2TSN(hyperperiod int, method string) {
	datasize_max := 100.
	deadline := 5000
	step 	 := 1000
	if method == "obo"{
		can2tsnFlowSet.Method=method
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			queue := &Queue{}
			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
				// example:
				// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
				// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
				for len(queue.Streams) > 0 {
					s := queue.Streams[0]
					// 逾期 → 丟棄
					if current_time > s.FinishTime {
						can2tsnFlowSet.O1_Drop++
						queue.popQueue(1)
						continue
					}
					// 立即封裝一個 CAN → TSN
					flushStream(can2tsnFlow, current_time, queue.Streams[0].DataSize, deadline)
					if len(queue.Streams) > 0 { queue.popQueue(1) }
				}
			}	
		}
	
	}else if method=="wat"{
		can2tsnFlowSet.Method=method
		
		safe_deadline := 5000
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			queue := &Queue{}
			datasize_count := 0.
			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
				queue.sortQueue(method, current_time)
				soon_count := float64(0)
				// example:
				// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
				// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
				head := 0
    			for head < len(queue.Streams) {
					stream := queue.Streams[head]

					if current_time > stream.FinishTime {
						can2tsnFlowSet.O1_Drop += 1
						head ++
						continue
					} 
					datasize_count += stream.DataSize
					head ++
					if datasize_count >= datasize_max {
						flushStream(can2tsnFlow, current_time, datasize_count , deadline)
						datasize_count = 0
					}
					
				}
				// 剪掉已處理 head 部分
				queue.popQueue(head)
				for len(queue.Streams) > 0 &&
					queue.Streams[0].FinishTime-current_time <= safe_deadline {

					soon_count += queue.Streams[0].DataSize   // 累加這筆
					queue.Streams = queue.Streams[1:]             // pop 這筆
				}
				// 如果累積量 >0，就一次封裝送出
				if soon_count > 0 {

					flushStream(can2tsnFlow, current_time, datasize_count+float64(soon_count), deadline)
					datasize_count = 0
				}		
			}
			if datasize_count > 0 {				
				flushStream(can2tsnFlow, hyperperiod, datasize_count , deadline)
				datasize_count = 0
			}
		}
	}else{
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			can2tsnFlowSet.Method=method
			queue := &Queue{}
			datasize_count := 0.
			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
				queue.sortQueue(method, current_time)
	
				// example:
				// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
				// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
				head := 0
				for head < len(queue.Streams) {
					stream := queue.Streams[head]

					if current_time > stream.FinishTime {
						can2tsnFlowSet.O1_Drop += 1
						head ++
						continue
					} 
					datasize_count += stream.DataSize
					head ++
					if datasize_count >= datasize_max {
						flushStream(can2tsnFlow, current_time, datasize_count , deadline)
						datasize_count = 0
					}
					
				}
				if datasize_count > 0 && current_time % 5000 == 0 {
					flushStream(can2tsnFlow, current_time, datasize_count , deadline)
					datasize_count = 0
				}
			}
			if datasize_count > 0 {				
				flushStream(can2tsnFlow, hyperperiod, datasize_count , deadline)
				datasize_count = 0
			}
		}
	}
	

}

type CAN2TSN_Flow_Set struct {
	Method			string
	CAN2TSN_Flows 	[]*CAN2TSN_Flow
	O1_Drop       	int
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
	can2tsnFlow.CAN2TSN_Flow = createCAN2TSNFlow(f.Source, f.Destination , f.Period, f.Deadline, f.DataSize)
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

func flushStream(flow *CAN2TSN_Flow, now int, packedSize float64, dl int) {
	
	if packedSize > 64{
		stream := createCAN2TSNStream(now, dl, float64(packedSize))
    	flow.CAN2TSN_Flow.Streams = append(flow.CAN2TSN_Flow.Streams, stream)
	}else{
		stream := createCAN2TSNStream(now, dl, 64)
    	flow.CAN2TSN_Flow.Streams = append(flow.CAN2TSN_Flow.Streams, stream)
	}

}

func createCAN2TSNFlow(source int, destination int , period int, deadline int , datasize float64) *Flow {
	newFlow := &Flow{
		Period:      period,
		Deadline:    deadline,
		DataSize:    datasize,
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

func (q *Queue) popQueue(head int) {
	q.Streams = q.Streams[head:]
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

	case "wat":
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
