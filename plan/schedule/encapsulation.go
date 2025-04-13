package schedule

import (
	// "src/network"
	"src/network/flow"
	// "sync"
	// "src/plan/algo_timer"
	"fmt"
	"sort"
)


type MessageQueues struct{
	Queue 		[]*Queue
}

type Queue struct{
	Source			int
	Destination 	int
	Flow			*flow.Flow		// CAN→TSN 聚合後的 TSN Flow
	CanStreams 		[]*flow.Stream	// 原始 CAN Streams
}

func (mq *MessageQueues)createQueue(f *flow.Flow){
	queue := &Queue{}

	queue.Source=f.Source
	queue.Destination=f.Destination
	queue.Flow = createCAN2TSNFlow(f.Source,f.Destination)
	queue.CanStreams=append(queue.CanStreams,f.Streams...)

	mq.Queue=append(mq.Queue, queue)
}

func (q *Queue)saveCanStream(f *flow.Flow){
	
	q.CanStreams=append(q.CanStreams,f.Streams...)
}

func (mq *MessageQueues)searchQueue(f *flow.Flow){

	s :=f.Source
	d :=f.Destination

	
	for _,queue:=range mq.Queue{
		if queue.Source ==s && queue.Destination==d{
			queue.saveCanStream(f)
			return
		}
	}
	mq.createQueue(f)
	
}

func EncapsulateCAN2TSN(f *flow.CANFlows) *MessageQueues{
	mq	:=	&MessageQueues{}

	for _, impf := range f.ImportantCANFlows{
		// fmt.Printf("Source: %v ,Destinatione: %v , Datasize: %v \n",impf.Source, impf.Destination, impf.DataSize)
		mq.searchQueue(impf)
    }

	for _, unimpf := range f.UnimportantCANFlows{
		// fmt.Printf("Source: %v ,Destinatione: %v , Datasize: %v \n",unimpf.Source, unimpf.Destination, unimpf.DataSize)
		mq.searchQueue(unimpf)
    }
	//檢查CAN NODE 檢查stream
	//現在有每個queue了 接下來就是咬把每個queue進行封裝 
	
	for _, q:= range mq.Queue{

		sortCANStreams(q.CanStreams, "fifo")

		maxdatasize := 64
		count := 0
		arrivalTime := 0
		deadline := 5000

		// fmt.Printf("%v\n",len(q.CanStreams))
		for ind,stream := range q.CanStreams{
			if count >= maxdatasize{
				q.Flow.Streams=append(q.Flow.Streams,createCAN2TSNStream(arrivalTime,deadline,float64(maxdatasize)))

				count = int(stream.DataSize)
				arrivalTime	+=	5000
				deadline	+=	5000
			}else{
				count += int(stream.DataSize)

			}
			if ind == len(q.CanStreams)-1{
				q.Flow.Streams=append(q.Flow.Streams,createCAN2TSNStream(arrivalTime,deadline,float64(maxdatasize)))
				count = int(stream.DataSize)

			}
			
		}
	}
	// mq.Show_MQ()
	return mq

}

func sortCANStreams(streams []*flow.Stream, strategy string) {

	switch strategy {
	case "fifo":
		// 到達時間小 → 大
		sort.Slice(streams, func(i, j int) bool {
			return streams[i].ArrivalTime < streams[j].ArrivalTime
		})
	case "deadline":
		// Deadline 小 → 大（最急先送）
		sort.Slice(streams, func(i, j int) bool {
			return streams[i].Deadline < streams[j].Deadline
		})
	case "datasize":
		// 資料量大 → 小
		sort.Slice(streams, func(i, j int) bool {
			return streams[i].DataSize > streams[j].DataSize
		})
	default:
		// 預設 FIFO
		sort.Slice(streams, func(i, j int) bool {
			return streams[i].ArrivalTime < streams[j].ArrivalTime
		})
	}
}

func (mq *MessageQueues) Show_MQ() {

	for _, q := range mq.Queue {
		fmt.Printf("Queue (%d→%d) streams=%d\n",
			q.Source, q.Destination, len(q.CanStreams))
			for ind, stream := range q.Flow.Streams{
				fmt.Printf("	Stream %v ,ArrivalTime: %v  ,Deadline: %v ,Datasize: %v\n" ,ind , stream.ArrivalTime , stream.Deadline , stream.DataSize)
			}
		
	}
}

func createCAN2TSNFlow(source int, destination int) *flow.Flow{

	newFlow := &flow.Flow{
		Period:      5000,
		Deadline:    5000,
		DataSize:    100,
		Source:      source,
		Destination: destination, 
	}

	return newFlow
}

func createCAN2TSNStream(arrivaltime int, deadline int,datasize float64) *flow.Stream{

	newStream := &flow.Stream{
		ArrivalTime: arrivaltime,
		Deadline:    deadline,       
		DataSize:    datasize,
	}

	return newStream
}