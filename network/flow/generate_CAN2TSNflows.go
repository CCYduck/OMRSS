package flow

import (
	"fmt"
	"sort"
	"time"
	"encoding/json"
	"log"
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
	for _, name := range []string{"fifo", "priority", "obo", "wat"} {
		fsCopy := can2tsnFlowSet.DeepCopyCAN2TSN()   // <- 自己寫或用 github.com/jinzhu/copier
		if fsCopy == nil { log.Println("deep copy failed"); continue }
		fsCopy.O1_Drop = 0    
		fsCopy.DatasizeCount = 0
		fsCopy.TSNFrameCount = 0

		start := time.Now()
		fsCopy.EncapsulateCAN2TSN(hyperperiod, name)

		method_struct := &Method{
			Method_Name:       name,	

			CAN2TSN_Delay:     time.Since(start),
			CAN2TSN_O1_Drop:   fsCopy.O1_Drop,
			CAN2TSNFlows:      make([]*Flow, 0, len(fsCopy.CAN2TSN_Flows)),
		}
		
		for _, can2tsn_flow := range fsCopy.CAN2TSN_Flows {
			method_struct.CAN2TSNFlows = append(method_struct.CAN2TSNFlows, can2tsn_flow.CAN2TSN_Flow)
		}
		method_struct.BytesSent = fsCopy.DatasizeCount
		method_struct.TSNFrameCount = fsCopy.TSNFrameCount  
		flow_set.Encapsulate = append(flow_set.Encapsulate, method_struct)
		
	}
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set) EncapsulateCAN2TSN(hyperperiod int, method string) {
	// datasize_least := 64.
	datasize_max := 100.
	period := 5000
	deadline := 5000
	step 	 := 1000
	if method == "obo"{
		// can2tsnFlowSet.Method=method
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			queue := &Queue{}
			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))

				// 1) 先把已逾期的丟掉 (跟你原本一樣)
				drop := 0
				for drop < len(queue.Streams) && current_time > queue.Streams[drop].FinishTime {
					can2tsnFlowSet.O1_Drop++
					drop++
				}
				queue.popQueue(drop)
	
				// example:
				// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
				// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
				for len(queue.Streams) > 0 {
					
					// 立即封裝一個 CAN → TSN
					can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_max, deadline)
					// can2tsnFlowSet.DatasizeCount+= queue.Streams[0].DataSize
					// can2tsnFlowSet.TSNFrameCount+=1
					queue.popQueue(1)
				}
			}	
		}
	
	}else if method=="wat"{
		// can2tsnFlowSet.Method=method
		// minLoad := 64.
		safe_deadline := 1500
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			queue := &Queue{}
			datasize_count := 0.
			
			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
				queue.sortQueue(method, current_time)

				// 1) 先把已逾期的丟掉
				drop := 0
				for drop < len(queue.Streams) && current_time > queue.Streams[drop].FinishTime {
					can2tsnFlowSet.O1_Drop++
					drop++
				}
				queue.popQueue(drop)
				
				// 2) 只要佇列裡 **還有** imminent stream，就一直封裝
				for hasImminent(queue, current_time, safe_deadline)  {
					head := 0
					for head < len(queue.Streams) && datasize_count + queue.Streams[head].DataSize < datasize_max{
						datasize_count +=  queue.Streams[head].DataSize
						head ++
					}
					can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_max , deadline)
					datasize_count = 0
					// 剪掉已處理 head 部分
					queue.popQueue(head)
				}
				
			}
			// 4. hyperperiod 結束後，佇列可能還有殘留，都打一包送掉
			if len(queue.Streams) > 0 {
				pack := 0.0
				for _, s := range queue.Streams {
					pack += s.DataSize
				}
			 can2tsnFlowSet.flushStream(can2tsnFlow, hyperperiod, pack, deadline)
			 datasize_count += pack
		 	}
		}
	}else{
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			// can2tsnFlowSet.Method=method
			queue := &Queue{}
			datasize_count := 0.
			
			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
				queue.sortQueue(method, current_time)

				// 1) 先把已逾期的丟掉 (跟你原本一樣)
				drop := 0
				for drop < len(queue.Streams) && current_time > queue.Streams[drop].FinishTime {
					can2tsnFlowSet.O1_Drop++
					drop++
				}
				queue.popQueue(drop)
	
				// example:
				// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
				// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
				head := 0
				for head < len(queue.Streams) {
					stream := queue.Streams[head]				
					datasize_count += stream.DataSize
					head ++
					if datasize_count >= datasize_max{
						can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_max , deadline)
						datasize_count = 0
					}
				}
				queue.popQueue(head)
				// if datasize_count > 0 && len(queue.Streams) > 0 && current_time % period == 0 {
				if datasize_count > 0  && current_time % period == 0 {
					can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_max , deadline)
					datasize_count = 0
				}
			}
			if datasize_count > 0 {				
				can2tsnFlowSet.flushStream(can2tsnFlow, hyperperiod, datasize_max , deadline)
				datasize_count = 0
			}
		}
	}
	

}

type CAN2TSN_Flow_Set struct {
	Method			string
	CAN2TSN_Flows 	[]*CAN2TSN_Flow
	DatasizeCount	float64
	TSNFrameCount	int
	O1_Drop       	int
}

type CAN2TSN_Flow struct {
	Source       int
	Destination  int
	CAN_Streams  []*Stream
	CAN2TSN_Flow *Flow
	Datasize_Count	float64
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set)flushStream(flow *CAN2TSN_Flow, now int, packedSize float64, dl int) {
	
	if packedSize < 64 {packedSize = 64}

	stream := createCAN2TSNStream(now, dl, packedSize)
	flow.CAN2TSN_Flow.Streams = append(flow.CAN2TSN_Flow.Streams, stream)

	can2tsnFlowSet.DatasizeCount+= packedSize
	can2tsnFlowSet.TSNFrameCount+=1
}

// --------------- helper：判斷佇列中是否還有 imminent stream ---------------
func hasImminent(q *Queue, now, safe int) bool {
	for _, s := range q.Streams {
		if s.FinishTime-now <= safe {
			return true
		}
	}
	return false
}

// DeepCopyCAN2TSN returns a deep-cloned *CAN2TSN_Flow_Set.
func (can2tsnFlowSet *CAN2TSN_Flow_Set)DeepCopyCAN2TSN() *CAN2TSN_Flow_Set {
    if can2tsnFlowSet == nil { return nil }

    data, err := json.Marshal(can2tsnFlowSet)
    if err != nil { return nil }

    dst := &CAN2TSN_Flow_Set{}
    if err := json.Unmarshal(data, dst); err != nil {
        return nil
    }
    return dst
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
