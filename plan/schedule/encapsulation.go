package schedule

import (
	// "src/network"
	"src/network/flow"
	"sync"
	// "time"
	// "fmt"
)

//network b 1.堆疊 2.根據週期進入堆疊 3.根據最大封裝數量12個進行封裝 4.然後生成TSN flow Datasize 100bytes , period 5000 , deadline 5000us
// 全域堆疊結構(視你需求放哪裡)，用 map[destID][]flow.Flow
var (
    flowStack = make(map[int][]flow.Flow)
    stackLock sync.Mutex
)

func EncapsulateCAN2TSN(source int, target int, datasize float64 , deadline int) (float64, *flow.CAN2TSNFlow){

	stackLock.Lock()
    defer stackLock.Unlock()

    dest := target // 假設只考慮單一目的地

	var pkt *flow.CAN2TSNFlow
	flow := flow.Flow{
		Source: 		source,
		Destination:	target,
		// DataSize: 		datasize,	
		Deadline: 		deadline,
	}
    // push flow 進入對應目的地堆疊
    flowStack[dest] = append(flowStack[dest], flow)

    // 若堆疊長度達到 5，就執行封裝
    if len(flowStack[dest]) >= 5 {
        flowsToEncap := flowStack[dest][:5]       // 取前 5
        flowStack[dest] = flowStack[dest][5:]     // 移除前 5
        // 建立封裝產物
		

		// 然後在 struct literal 裡使用 flowsPtr
		pkt.Source =	source
		pkt.Destination = target


        // 計算大小
        var sumSize float64
        for _, cf := range flowsToEncap {
            sumSize += cf.DataSize // flow 中的 DataSize
			pkt.Deadline=append(pkt.Deadline, cf.Deadline)
        }
        pkt.DataSize = sumSize

        // 計算 delay
        // 例如：delay = sumSize / 125
        delay := sumSize / 125.0

        return delay, pkt
    }
    // 尚未達 5 條，不封裝
    return 0, nil
}