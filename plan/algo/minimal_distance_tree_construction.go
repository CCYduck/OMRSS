package algo

import (
	// "fmt"
	"src/network"
	"src/plan/algo_timer"
	"src/plan/path"
	"src/plan/routes"
)

func (mtdc *MDTC) MDTC_Run(network *network.Network) {
	// 5. DistanceTree
	mtdc.Timer = algo_timer.NewTimer()
	mtdc.Timer.TimerStart()
	mtdc.Trees = routes.Get_DistanceTree_Routing(network)
	// mtdc.Trees = path.BestPath(network)
	mtdc.Timer.TimerStop()
}

func (SP *SP) SP_Run(network *network.Network) {
	// 5. DistanceTree
	SP.Timer = algo_timer.NewTimer()
	SP.Timer.TimerStart()
	// mtdc.Trees = routes.Get_DistanceTree_Routing(network)
	SP.Path = path.BestPath(network)
	// fmt.Println(SP.Path)
	SP.Timer.TimerStop()
}

func (KP *KP) KP_Run(network *network.Network) {
	// 5. DistanceTree
	KP.Timer = algo_timer.NewTimer()
	KP.Timer.TimerStart()
	// mtdc.Trees = routes.Get_DistanceTree_Routing(network)
	
	KP.KPath = path.KShortestPath(network)
	KP.Timer.TimerStop()
}

