package flow

import (
	"fmt"
	"sort"
	"time"
	"encoding/json"
	"log"
	// "math"
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
	for _, name := range []string{"fifo", "priority", "obo", "wst", "mao"} {
		fsCopy := can2tsnFlowSet.DeepCopyCAN2TSN()   // <- 自己寫或用 github.com/jinzhu/copier
		// a := can2tsnFlowSet.CAN2TSN_Flows[0].CAN_Streams[0]
		// b := fsCopy.CAN2TSN_Flows[0].CAN_Streams[0]

		// fmt.Println("Same pointer?", a == b)

		if fsCopy == nil { log.Println("deep copy failed"); continue }
		fsCopy.O1_Encap_Drop = 0    

		fsCopy.DatasizeCount = 0
		fsCopy.TSNFrameCount = 0

		start := time.Now()
		fsCopy.EncapsulateCAN2TSN(hyperperiod, name)

		method_struct := &Method{
			Method_Name:       name,	
			CAN2TSN_Delay:     time.Since(start),
			CAN2TSN_O1_Drop:   fsCopy.O1_Encap_Drop,
			CAN_Area_O1_Drop:  fsCopy.O1_Decap_Drop,
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
	datasize_least := 64.
	datasize_max := 1500.
	period := 5000
	deadline := 5000
	step 	 := 1000
	canSpeedBps  := 1_000_000.0 	//CAN bandwidth 1 Mbps
	bytesPerStep := canSpeedBps / 8 * float64(step) / 1_000_000.0 // = 125B

	if method == "obo"{
		send_queue := map[int]*Queue{}
		getQ := func(dst int) *Queue {
			if q, ok := send_queue[dst]; ok {
				return q
			}
			q := &Queue{}
			send_queue[dst] = q
			return q
		}
		// can2tsnFlowSet.Method=method
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			queue := &Queue{}

			for current_time := 0; current_time < hyperperiod; current_time += step {
				stream := can2tsnFlow.getStreamsByCurrentTime(current_time)
				queue.appendQueue(stream)

				// 1) 先把已逾期的丟掉 (跟你原本一樣)
				drop := 0
				for drop < len(queue.Streams) && current_time > queue.Streams[drop].FinishTime {
					can2tsnFlowSet.O1_Encap_Drop++
					drop++
				}
				if drop > 0 {
					queue.popQueue(drop)
				}
				
				// example:
				// current_time=0 queue=[0_1, 0_2, 0_3, 0_4, 0_5] if 05 datasize_count>datasize_max, createCAN2TSNStream datasize_count = 0
				// current_time=5000 queue=[0_5, 5000_1, 5000_2, 5000_3, 5000_4, 5000_5]
				for len(queue.Streams) > 0 {						
					// 立即封裝一個 CAN → TSN
					can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_least, deadline)
					// fmt.Println(queue.Streams[0].DataSize)
					send_stream := queue.Streams[0]
					send_stream.ArrivalTime = current_time
					sq := getQ(can2tsnFlow.Destination)
					sq.Streams=append(sq.Streams, send_stream)
					// can2tsnFlowSet.DatasizeCount+= queue.Streams[0].DataSize
					// can2tsnFlowSet.TSNFrameCount+=1
					queue.popQueue(1)
				}

			}	
		}
		for dst , sq := range send_queue{
			fmt.Println(method, dst, len(sq.Streams))
			// fmt.Println("Before ",can2tsnFlowSet.O1_Decap_Drop)
			for currentTime := 0; currentTime < hyperperiod; currentTime += step {
				// sq.sortQueue(method, currentTime)
				remaining := bytesPerStep
				
				i := 0
				for i < len(sq.Streams) {
					s := sq.Streams[i]

					if s.FinishTime > 0 && s.FinishTime < currentTime {
						// fmt.Println(s.ArrivalTime)
						can2tsnFlowSet.O1_Decap_Drop++
						sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
						continue
					}

					if s.ArrivalTime > currentTime {
						i++
						continue
					}

					if float64(s.DataSize) > remaining {
						break
					}
					// 傳送封包
					remaining -= float64(s.DataSize)
					// 從 queue 中移除
					sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
				}
				
			}
			// fmt.Println("After ",can2tsnFlowSet.O1_Decap_Drop)
		}
	}else if method=="wst"{
		// can2tsnFlowSet.Method=method
		// minLoad := 64.
		const guardBase = 1200       // µs 讓「何時必須封裝」隨著佇列最緊迫的剩餘時間調整，而保留一段可以把封包真正送出去的 guard
		
		send_queue := map[int]*Queue{}
		getQ := func(dst int) *Queue {
			if q, ok := send_queue[dst]; ok {
				return q
			}
			q := &Queue{}
			send_queue[dst] = q
			return q
		}
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			queue := &Queue{}
			datasize_count := 0.
			a:=0
			for _, stream := range can2tsnFlow.CAN_Streams {
				if stream.ArrivalTime <=  hyperperiod{
					a++
				}
				
			}
			// fmt.Printf("can2tsnFlow S: %v , D: %v , Count: %v \n", can2tsnFlow.Source, can2tsnFlow.Destination, a)

			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
				// queue.sortQueue(method, current_time)

				// 1) 先把已逾期的丟掉
				drop := 0
				for drop < len(queue.Streams) && current_time > queue.Streams[drop].FinishTime {
					can2tsnFlowSet.O1_Encap_Drop++
					drop++
				}
				queue.popQueue(drop)
				
				guard := guardBase + len(queue.Streams) * 520      // 動態 guard

				if len(queue.Streams) == 0 { continue }
				// 2) 只要佇列裡 **還有** imminent stream，就一直封裝
				for hasImminent(queue, current_time, guard)  {
					head := 0
					for head < len(queue.Streams) && datasize_count + queue.Streams[head].DataSize < datasize_max{
						datasize_count +=  queue.Streams[head].DataSize
						send_stream := queue.Streams[head]
						send_stream.ArrivalTime = current_time
						sq := getQ(can2tsnFlow.Destination)
						sq.Streams=append(sq.Streams, send_stream)			
						head ++
					}

					can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_count , deadline)	
					
					// 剪掉已處理 head 部分
					queue.popQueue(head)
					datasize_count = 0
				}			
			}
			// 4. hyperperiod 結束後，佇列可能還有殘留，都打一包送掉
			if len(queue.Streams) > 0 {	
				for _, stream := range queue.Streams {
					stream.ArrivalTime = hyperperiod
					sq := getQ(can2tsnFlow.Destination)
					sq.Streams = append(sq.Streams, stream)
				}	
				can2tsnFlowSet.flushStream(can2tsnFlow, hyperperiod, datasize_count , deadline)
				datasize_count = 0
		 	}
		}
		for dst,sq := range send_queue{
			fmt.Println(method, dst, len(sq.Streams))
			// fmt.Println("Before ",can2tsnFlowSet.O1_Decap_Drop)
			for currentTime := 0; currentTime < hyperperiod; currentTime += step {
				// sq.sortQueue(method, currentTime)
				remaining := bytesPerStep
				
				i := 0
				for i < len(sq.Streams) {
					s := sq.Streams[i]

					if s.FinishTime > 0 && s.FinishTime < currentTime {
						can2tsnFlowSet.O1_Decap_Drop++
						// fmt.Println(s.ArrivalTime)
						sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
						continue
					}

					if s.ArrivalTime > currentTime {
						i++
						continue
					}

					if float64(s.DataSize) > remaining {
						break
					}
					// 傳送封包
					remaining -= float64(s.DataSize)
					// 從 queue 中移除
					sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
				}
				
			}
			// fmt.Println("After ",can2tsnFlowSet.O1_Decap_Drop)
		}
	}else if method=="mao"{
		send_queue := map[int]*Queue{}
		getQ := func(dst int) *Queue {
			if q, ok := send_queue[dst]; ok {
				return q
			}
			q := &Queue{}
			send_queue[dst] = q
			return q
		}
		mtuLimit := 1500.0 / 2          // MTU 限定為 750 Bytes
		
		clusters := make(map[clusterKey][]*Stream)

		// 建立聚合群：同 source, dest, period 的 stream
		for _, flow := range can2tsnFlowSet.CAN2TSN_Flows {
			key := clusterKey{flow.Source, flow.Destination, flow.CAN2TSN_Flow.Period}
			clusters[key] = append(clusters[key], flow.CAN_Streams...)
		}

		for key, streams := range clusters {
			// 依 deadline 排序（急的先封裝）
			sort.Slice(streams, func(i, j int) bool {
				return streams[i].Deadline < streams[j].Deadline
			})
			sq := getQ(key.Dest)
			frame := []*Stream{}
			payload := 0.0

			for _, s := range streams {
				
				if payload+s.DataSize <= mtuLimit {
					frame = append(frame, s)
					payload += s.DataSize
				} else {
					// 送出 TSN Frame
					fullSize := payload
					if !schedulable(frame, fullSize, bytesPerStep) {
						// EMSO 拆解重封裝
						frame1, frame2 := disaggregateByDeadline(frame)
						repackAndInsert(frame1, key, can2tsnFlowSet, bytesPerStep)
						repackAndInsert(frame2, key, can2tsnFlowSet, bytesPerStep)
						
						sq.Streams = append(sq.Streams, frame1...)
						sq.Streams = append(sq.Streams, frame2...)
					} else {
						repackAndInsert(frame, key, can2tsnFlowSet, bytesPerStep)
						sq.Streams = append(sq.Streams, frame...)
					}
					
					// 開始新 frame
					frame = []*Stream{s}
					payload = s.DataSize
				}
			}

			// 收尾未滿 MTU 的 frame
			if len(frame) > 0 {
				fullSize := payload
				if !schedulable(frame, fullSize, bytesPerStep) {
					frame1, frame2 := disaggregateByDeadline(frame)
					repackAndInsert(frame1, key, can2tsnFlowSet, bytesPerStep)
					repackAndInsert(frame2, key, can2tsnFlowSet, bytesPerStep)
					sq.Streams = append(sq.Streams, frame1...)
					sq.Streams = append(sq.Streams, frame2...)
				} else {
					repackAndInsert(frame, key, can2tsnFlowSet, bytesPerStep)
					sq.Streams = append(sq.Streams, frame...)
				}
			}
		}
		for dst, sq := range send_queue{
			fmt.Println(method, dst, len(sq.Streams))
			// fmt.Println("Before ",can2tsnFlowSet.O1_Decap_Drop)
			for currentTime := 0; currentTime < hyperperiod; currentTime += step {
				// sq.sortQueue("fifo", currentTime)
				remaining := bytesPerStep

				i := 0
				for i < len(sq.Streams) {
					s := sq.Streams[i]

					if s.FinishTime > 0 && s.FinishTime < currentTime {
						// fmt.Println(s.ArrivalTime)
						can2tsnFlowSet.O1_Decap_Drop++
						sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
						continue
					}

					if s.ArrivalTime > currentTime {
						i++
						continue
					}

					if float64(s.DataSize) > remaining {
						break
					}
					// 傳送封包
					remaining -= float64(s.DataSize)
					// 從 queue 中移除
					sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
				}
			}
			// fmt.Println("After ",can2tsnFlowSet.O1_Decap_Drop)
		}
		
	}else{
		send_queue := map[int]*Queue{}
		getQ := func(dst int) *Queue {
			if q, ok := send_queue[dst]; ok {
				return q
			}
			q := &Queue{}
			send_queue[dst] = q
			return q
		}
		for _, can2tsnFlow := range can2tsnFlowSet.CAN2TSN_Flows {
			// can2tsnFlowSet.Method=method
			queue := &Queue{}
			datasize_count := 0.
			// i := 0
			for current_time := 0; current_time < hyperperiod; current_time += step {
				queue.appendQueue(can2tsnFlow.getStreamsByCurrentTime(current_time))
				queue.sortQueue(method, current_time)

				// 1) 先把已逾期的丟掉 (跟你原本一樣)
				drop := 0
				for drop < len(queue.Streams) && current_time > queue.Streams[drop].FinishTime {
					can2tsnFlowSet.O1_Encap_Drop++
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
						can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_count , deadline)
						for _, stream := range queue.Streams[:head] {
							stream.ArrivalTime = current_time
							sq := getQ(can2tsnFlow.Destination)
							sq.Streams = append(sq.Streams, stream)
						}
						datasize_count = 0
						queue.popQueue(head)
						head = 0
					}
				}
				
				// if datasize_count > 0 && len(queue.Streams) > 0 && current_time % period == 0 {
				if datasize_count > 0  && current_time % period == 0 {
					can2tsnFlowSet.flushStream(can2tsnFlow, current_time, datasize_count , deadline)
					for _, stream := range queue.Streams[:head] {
							stream.ArrivalTime = current_time
							sq := getQ(can2tsnFlow.Destination)
							sq.Streams = append(sq.Streams, stream)
					}
					datasize_count = 0
					queue.popQueue(head)
					head = 0
				}
			}
			if datasize_count > 0 {				
				can2tsnFlowSet.flushStream(can2tsnFlow, hyperperiod, datasize_count , deadline)
				for _, stream := range queue.Streams {
					stream.ArrivalTime = hyperperiod
					sq := getQ(can2tsnFlow.Destination)
					sq.Streams = append(sq.Streams, stream)
				}	
				datasize_count = 0
			}
		}
		for dst , sq := range send_queue{
			fmt.Println(method, dst, len(sq.Streams))
			// fmt.Println("Before ",can2tsnFlowSet.O1_Decap_Drop)
			for currentTime := 0; currentTime < hyperperiod; currentTime += step {
				// sq.sortQueue("fifo", currentTime)
				remaining := bytesPerStep

				i := 0
				for i < len(sq.Streams) {
					s := sq.Streams[i]

					if s.FinishTime > 0 && s.FinishTime < currentTime {
						// fmt.Println(s.ArrivalTime)
						can2tsnFlowSet.O1_Decap_Drop++
						sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
						continue
					}

					if s.ArrivalTime > currentTime {
						i++
						continue
					}

					if float64(s.DataSize) > remaining {
						break
					}
					// 傳送封包
					remaining -= float64(s.DataSize)
					// 從 queue 中移除
					sq.Streams = append(sq.Streams[:i], sq.Streams[i+1:]...)
				}
			}
			// fmt.Println("After ",can2tsnFlowSet.O1_Decap_Drop)
		}
	}

}

type CAN2TSN_Flow_Set struct {
	Method				string
	CAN2TSN_Flows 		[]*CAN2TSN_Flow
	DatasizeCount		float64
	TSNFrameCount		int
	O1_Encap_Drop       int
	O1_Decap_Drop		int
}

type CAN2TSN_Flow struct {
	Source       int
	Destination  int
	CAN_Streams  []*Stream
	CAN2TSN_Flow *Flow
	Datasize_Count	float64
}

func (can2tsnFlowSet *CAN2TSN_Flow_Set)flushStream(flow *CAN2TSN_Flow, now int, packedSize float64, dl int) {
	if packedSize < 64 {
		packedSize = 64
	}
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

type clusterKey struct {
	Source int
	Dest   int
	Period int
}

type FlowKey struct {
	Source      int
	Destination int
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

	case "wst":
		// WAT 根據剩餘時間
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

// ----- 輔助函數 -----
func schedulable(streams []*Stream, fullSize float64, bytesPerUs float64) bool {
	sendTime := int(fullSize / bytesPerUs)
	for _, s := range streams {
		if s.ArrivalTime+s.Deadline < s.ArrivalTime+sendTime {
			return false
		}
	}
	return true
}

func disaggregateByDeadline(streams []*Stream) ([]*Stream, []*Stream) {
	if len(streams) <= 1 {
		return streams, nil
	}
	sort.Slice(streams, func(i, j int) bool {
		return streams[i].Deadline < streams[j].Deadline
	})
	mid := len(streams) / 2
	return streams[:mid], streams[mid:]
}

func repackAndInsert(streams []*Stream, key clusterKey, flowSet *CAN2TSN_Flow_Set, bytesPerUs float64) {
	if len(streams) == 0 {
		return
	}

	payload := 0.0
	for _, s := range streams {
		payload += s.DataSize
	}
	fullSize := payload

	f := createCAN2TSNFlow(key.Source, key.Dest, key.Period, 5000, fullSize)
	f.Streams = append(f.Streams, streams...)
	// flowSet.addNewCAN2TSNFlowToSet(f)
	flowSet.searchCAN2TSNFlow(f)
	flowSet.DatasizeCount += fullSize
	flowSet.TSNFrameCount++
}
