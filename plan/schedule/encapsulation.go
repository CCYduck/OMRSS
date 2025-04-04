package schedule

import (
	// "src/network"
	"src/network/flow"
	"sync"
	// "src/plan/algo_timer"
	// "fmt"
	"sort"
)

//network b 1.堆疊 2.根據週期進入堆疊 3.根據最大封裝數量12個進行封裝 4.然後生成TSN flow Datasize 100bytes , period 5000 , deadline 5000us
// 全域堆疊結構(視你需求放哪裡)，用 map[destID][]flow.Flow
var (
    flowStack = make(map[int][]flow.Flow)
    stackLock sync.Mutex
	// timer = algo_timer.NewTimer()
)

func EncapsulateCAN2TSN(source int, target int, datasize float64 , deadline int) (float64, *flow.Flow){
	stackLock.Lock()
    defer stackLock.Unlock()

	// timer.TimerStart()
    dest := target // 假設只考慮單一目的地

	newFlow := flow.Flow{
		Source: 		source,
		Destination:	target,
		// DataSize: 		datasize,	
		Deadline: 		deadline,
	}
    // push flow 進入對應目的地堆疊
    flowStack[dest] = append(flowStack[dest], newFlow)

    // 若堆疊長度達到 5，就執行封裝
    if len(flowStack[dest]) >= 12 {
		
		sort.Slice(flowStack[dest], func(i, j int) bool {
            return flowStack[dest][i].Deadline < flowStack[dest][j].Deadline
        })
        flowsToEncap := flowStack[dest][:12]       // 取前 12
        flowStack[dest] = flowStack[dest][12:]     // 移除前 12
        // 建立封裝產物
		
		var pkt *flow.Flow
		pkt.Source 		=	source
		pkt.Destination = 	target
		pkt.Deadline	=	5000


        // 計算大小
        var sumSize float64
        for _, cf := range flowsToEncap {
            sumSize += cf.DataSize // flow 中的 DataSize
			// pkt.Deadline=append(pkt.Deadline, cf.Deadline)
        }
        pkt.DataSize = sumSize

        // 計算 delay
        // 例如：delay = sumSize / 125
        delay := sumSize / 125.0
		// timer.TimerEnd()
        return delay, pkt
    }
    // 尚未達 5 條，不封裝
    return 0, nil
}