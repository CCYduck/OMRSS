package schedule

import (
	// "src/network"
	"src/network/flow"
	// "sync"
	// "src/plan/algo_timer"
	"fmt"
	// "sort"
)


type MessageQueues struct{
	Queue 		[]*Queue
}

type Queue struct{
	Source			int
	Destination 	int
	Streams 		[]*flow.Stream
}

func (mq *MessageQueues)createQueue(f *flow.Flow){
	queue := &Queue{}

	queue.Source=f.Source
	queue.Destination=f.Destination
	queue.Streams=append(queue.Streams,f.Streams...)

	mq.Queue=append(mq.Queue, queue)
}

func (q *Queue)saveStream(f *flow.Flow){
	q.Streams=append(q.Streams,f.Streams...)
}

func (mq *MessageQueues)searchQueue(f *flow.Flow){
	s :=f.Source
	d :=f.Destination

	for _,queue:=range mq.Queue{
		if queue.Source ==s && queue.Destination==d{
			queue.saveStream(f)
		}
	}
	mq.createQueue(f)
}

func EncapsulateCAN2TSN(f *flow.CANFlows) (flow.Flow){
	
	mq	:=	&MessageQueues{}

	for _, impf := range f.ImportantCANFlows{
		fmt.Printf("Source: %v ,Destinatione: %v , Datasize: %v \n",impf.Source, impf.Destination, impf.DataSize)

		mq.searchQueue(impf)
    }

	for _, unimpf := range f.UnimportantCANFlows{
		fmt.Printf("Source: %v ,Destinatione: %v , Datasize: %v \n",unimpf.Source, unimpf.Destination, unimpf.DataSize)
		mq.searchQueue(unimpf)
    }//檢查CAN NODE 檢查stream
	//現在有每個queue了 接下來就是咬把每個queue進行封裝 
	mq.Show_MQ()

	
    return flow.Flow{}
}

func (mq *MessageQueues) Show_MQ() {

	for index, q := range mq.Queue {
		fmt.Printf("Queue %d Source: %v  ,Destination: %v\n",index , q.Source , q.Destination)
		
	}
}