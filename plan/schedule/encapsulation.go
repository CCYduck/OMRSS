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
	// fmt.Println(len(f.Streams))
}

func (q *Queue)saveCanStream(f *flow.Flow){
	q.CanStreams=append(q.CanStreams,f.Streams...)
	// fmt.Println(len(f.Streams))
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

func EncapsulateCAN2TSN(f *flow.CANFlows)(*MessageQueues,int){
	mq	:=	&MessageQueues{}

	o1_candrop:=0

	
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
		
		// fmt.Println(len(q.CanStreams))
		maxdatasize := 100
		count := 0
		arrivalTime := 0
		deadline := 5000
		// fmt.Println(len(q.CanStreams))
		current_time := 0
		method := "deadline"

		// fmt.Printf("%v\n",len(q.CanStreams))
		stack := &stack{}
		for ind,stream := range q.CanStreams{
			if stream.ArrivalTime==current_time{
				stack.appendstack(stream)
			}else{
				stack.sortstack(method)
				for _,stackstream := range stack.stack{
					if current_time	> stackstream.ArrivalTime+stackstream.Deadline{
						o1_candrop+=1
		
					}else{
						count += int(stackstream.DataSize)
						if count >= maxdatasize{
							// fmt.Println(ind,stream.ArrivalTime,stream.Deadline,current_time)
							q.Flow.Streams=append(q.Flow.Streams,createCAN2TSNStream(arrivalTime,deadline,float64(maxdatasize)))
							
							count = int(stackstream.DataSize)
							arrivalTime		+=	5000
							deadline		+=	5000
							current_time	+=	5000
							
							continue
						}
						// fmt.Println(len(stack.stack))
						stack.popstack()
						// fmt.Println(len(stack.stack))
					}			
				}
				
				stack.appendstack(stream)
				
			}
			if ind == len(q.CanStreams)-1{
				stack.sortstack(method)
				for _,stackstream := range stack.stack{
					if current_time	> stackstream.ArrivalTime+stackstream.Deadline{
						o1_candrop+=1
		
					}else{
						count += int(stackstream.DataSize)
						if count >= maxdatasize{
							// fmt.Println(ind,stream.ArrivalTime,stream.Deadline,current_time)
							q.Flow.Streams=append(q.Flow.Streams,createCAN2TSNStream(arrivalTime,deadline,float64(maxdatasize)))
							
							count = int(stackstream.DataSize)
							arrivalTime		+=	5000
							deadline		+=	5000
							current_time	+=	5000
							continue
						}
						stack.popstack()
					}		
					if  len(stack.stack)>=0 {
						q.Flow.Streams=append(q.Flow.Streams,createCAN2TSNStream(arrivalTime,deadline,float64(maxdatasize)))
					}
				}
			}
		}
	}
	// mq.Show_MQ()
	
	return mq,o1_candrop

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

type stack struct{
	stack		[]*flow.Stream
}

func (s *stack)appendstack(stream *flow.Stream){
	s.stack=append(s.stack, stream)
}

func (s *stack)popstack(){
	s.stack=s.stack[1:]
}

func (s *stack)sortstack(method string){
	switch method {
	case "fifo":
		// 到達時間小 → 大
		sort.Slice(s.stack, func(i, j int) bool {
			return s.stack[i].ArrivalTime <  s.stack[j].ArrivalTime
		})
	case "deadline":
		// Deadline 小 → 大（最急先送）
		sort.Slice( s.stack, func(i, j int) bool {
			return  s.stack[i].Deadline <  s.stack[j].Deadline
		})
	default:
		// 預設 FIFO
		sort.Slice( s.stack, func(i, j int) bool {
			return  s.stack[i].ArrivalTime <  s.stack[j].ArrivalTime
		})
}
}