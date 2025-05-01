package algo

import (
	"src/plan/algo_timer"
	"src/plan/routes"
	"src/plan/path"
)

type SP struct{
	Path		*path.Path_set
	InputFlow	*path.Path_set
	BGFlow		*path.Path_set
	Objs_SP		*[4]float64
	Timer 		*algo_timer.Timer
}

type SMT struct {
	Trees      *routes.Trees_set
	InputTrees *routes.Trees_set
	BGTrees    *routes.Trees_set
	Objs_smt   [4]float64
	Timer      *algo_timer.Timer
}

type MDTC struct {
	Trees      *routes.Trees_set
	InputTrees *routes.Trees_set
	BGTrees    *routes.Trees_set
	Objs_mdtc  [4]float64
	Timer      *algo_timer.Timer
}

type KP struct {
	KPath      	*path.KPath_Set
	InputKPath 	*path.KPath_Set
	BGKPath    	*path.KPath_Set
	Objs_kp   	[4]float64
	Timer      	*algo_timer.Timer
}

type OSACO struct {
	Timeout       	int
	K             	int
	P             	float64
	KTrees        	*routes.KTrees_set
	KPath		  	*path.KPath_Set
	VB            	*Visibility
	PRM           	*Pheromone
	InputTrees    	*routes.Trees_set
	InputPath	  	*path.Path_set
	InputKPaths   	*path.KPath_Set
	BGTrees       	*routes.Trees_set
	BGPath		  	*path.Path_set
	BGKPaths 	  	*path.KPath_Set
	Objs_osaco    	[5][4]float64        // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	Timer         	[5]*algo_timer.Timer // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}
	Method_Number 	int                  // 0: TOP K minimum weight 1: Increasing Arithmetic Sequence 2: Average Arithmetic Sequence
}

type Visibility struct {
	TSN_VB [][]float64
	AVB_VB [][]float64
	C2T_VB [][]float64
}

type Pheromone struct {
	TSN_PRM [][]float64
	AVB_PRM [][]float64
	C2T_PRM [][]float64
}
