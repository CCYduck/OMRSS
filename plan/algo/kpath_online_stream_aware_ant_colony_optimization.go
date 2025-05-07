package algo

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"src/network"
	"src/network/flow"
	"src/plan/algo_timer"
	"src/plan/path"
	"src/plan/schedule"
	"time"
)

var (

	bg_tsn int
	bg_avb int
)

func (osro *OSRO) OSRO_Initial_Settings(network *network.Network) {
	//// osro computing time: Estimate the time it takes to compute routing information

	bg_tsn = network.BG_TSN
	bg_avb = network.BG_AVB

	timer := algo_timer.NewTimer()
	timer.TimerStart()
	osro.KPath = path.KShortestPath(network)
	timer.TimerEnd()

	fmt.Println(len(osro.KPath.TSNPaths),len(osro.KPath.AVBPaths),len(osro.KPath.CAN2TSNPaths))
	osro.InputKPaths 	= osro.KPath.Input_kpath_set(bg_tsn, bg_avb)
	osro.BGKPaths 		= osro.KPath.BG_kpath_set(bg_tsn, bg_avb)

	osro.PRM = compute_prm(osro.KPath)
	osro.VB = compute_vb(osro.KPath, network.Flow_Set)

	osro.Timer[0] = algo_timer.NewTimer()
	osro.Timer[0].TimerMerge(timer)
	osro.Timer[1] = algo_timer.NewTimer()
	osro.Timer[1].TimerMerge(timer)
	osro.Timer[2] = algo_timer.NewTimer()
	osro.Timer[2].TimerMerge(timer)
	osro.Timer[3] = algo_timer.NewTimer()
	osro.Timer[3].TimerMerge(timer)
	osro.Timer[4] = algo_timer.NewTimer()
	osro.Timer[4].TimerMerge(timer)
}

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func (osro *OSRO) OSRO_Run(network *network.Network, timeout_index int) [4]float64 {
	// 6. osro
	// Repeat the execution of epochs within the timeout
	initialobj, initialcost := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath,"obo")
	fmt.Println()
	fmt.Printf("initial value: %d \n", initialcost)
	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", initialobj[0], initialobj[1], initialobj[2], initialobj[3])

	timeout := time.Duration(osro.Timeout) * time.Millisecond
	startTime := time.Now()
	i := 1
	for {
		fmt.Printf("\nepoch%d:\n", i)
		osro.Timer[timeout_index].TimerStart()
		II := epoch(network, osro, timeout_index)
		osro.Timer[timeout_index].TimerStop()

		_, cost1 := schedule.OBJ(network, osro.KPath, II, osro.BGPath,"obo")               	// new
		_, cost2 := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath,"obo") 		// old

		if cost1 < cost2 {
			osro.InputPath = II
			fmt.Println("Change the selected routing !!")
		}
		i += 1

		if time.Since(startTime) >= timeout {
			break
		}
	}

	resultobj, resultcost := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath,"obo")
	fmt.Println()
	fmt.Printf("result value: %d \n", resultcost)
	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", resultobj[0], resultobj[1], initialobj[2], resultobj[3])
	fmt.Println()

	if resultobj[0] != 0 || resultobj[1] != 0 {
		osro.Timer[timeout_index].TimerMax()
	}

	return resultobj
}

func compute_prm(X *path.KPath_Set) *Pheromone {
	pheromone := &Pheromone{}

	for nth, kpath := range X.TSNPaths {
		var prm []float64
		for i := 0; i < len(kpath.Paths); i++ {
			if nth < bg_tsn {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.TSN_PRM = append(pheromone.TSN_PRM, prm)
	}

	for nth, kpath := range X.CAN2TSNPaths {
		var prm []float64
		for i := 0; i < len(kpath.Paths); i++ {
			if nth < bg_tsn {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.C2T_PRM = append(pheromone.C2T_PRM, prm)
	}

	for nth, kpath := range X.AVBPaths {
		var prm []float64
		for i := 0; i < len(kpath.Paths); i++ {
			if nth < bg_avb {
				prm = append(prm, 0.5)
			} else {
				prm = append(prm, 1.)
			}
		}
		pheromone.AVB_PRM = append(pheromone.AVB_PRM, prm)
	}

	return pheromone
}

func compute_vb(X *path.KPath_Set, flow_set *flow.Flows) *Visibility {
	var preference float64 = 2.
	Input_flow_set := flow_set.Input_TSNflow_set()
	BG_flow_set := flow_set.BG_flow_set()

	visibility := &Visibility{}
	// osro CompVB
	// TSN flow
	for nth, tsn_kpath := range X.TSNPaths {
		var v []float64
		for kth := range tsn_kpath.Paths {
			mult := 1.
			if nth < bg_tsn && kth == 0 {
				mult = preference
			}

			//value := mult / float64(tsn_ktree.Trees[0].Weight) // mult / Tree weight
			value := mult / math.Exp(float64(tsn_kpath.Paths[0].Weight)) // mult / exponential function( Tree weight )
			v = append(v, value)
		}
		visibility.TSN_VB = append(visibility.TSN_VB, v)
	}

	// osro CompVB
	// CAN2TSN flow
	for nth, c2t_kpath := range X.CAN2TSNPaths {
		var v []float64
		for kth := range c2t_kpath.Paths {
			mult := 1.
			if nth < bg_tsn && kth == 0 {
				mult = preference
			}

			//value := mult / float64(tsn_ktree.Trees[0].Weight) // mult / Tree weight
			value := mult / math.Exp(float64(c2t_kpath.Paths[0].Weight)) // mult / exponential function( Tree weight )
			v = append(v, value)
		}
		visibility.C2T_VB = append(visibility.C2T_VB, v)
	}

	// osro CompVB
	// AVB flow
	for nth, avb_kpath:= range X.AVBPaths {
		var v []float64
		for kth, z := range avb_kpath.Paths {
			mult := 1.
			if nth < bg_avb && kth == 0 {
				mult = preference
			}

			if nth >= bg_avb {
				//fmt.Printf("Input flow%d tree%d \n", nth, kth)
				value := mult / float64(schedule.WCD(z, X, Input_flow_set.AVBFlows[nth-bg_avb], flow_set))
				v = append(v, value)

			} else {
				//fmt.Printf("Backgourd flow%d tree%d \n", nth, kth)
				value := mult / float64(schedule.WCD(z, X, BG_flow_set.AVBFlows[nth], flow_set))
				v = append(v, value)
			}
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	return visibility
}

func probability(osro *OSRO) (*path.Path_set, *path.Path_set, [2][]int, [2][]int) {
	var (
		input_k_location [2][]int // (tsn k index, avb k index)
		bg_k_location    [2][]int // (tsn k index, avb k index)
	)

	II := &path.Path_set{}
	II_prime := &path.Path_set{}
	
	for nth, kpath := range osro.KPath.TSNPaths{
		Denominator := 0.
		if len(kpath.Paths) == 0 { continue }
		for kth := range kpath.Paths {
			Denominator += osro.VB.TSN_VB[nth][kth] * osro.PRM.TSN_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range kpath.Paths {
			probability := (osro.VB.TSN_VB[nth][kth] * osro.PRM.TSN_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				// if kth == 5 => arr[0,0,0,0,0,0,...,1,1,1,...,2,2,2,2,..,3,3,3,...,4,4,4,4,...] len(arr) ~ 100
				arr = append(arr, kth)
			}
		}
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t :=kpath.Paths[n]

		if nth < bg_tsn {
			bg_k_location[0] = append(bg_k_location[0], n)
			II_prime.TSNPath = append(II_prime.TSNPath, t)
		} else {
			input_k_location[0] = append(input_k_location[0], n)
			II.TSNPath = append(II.TSNPath, t)
		}
	}

	for nth, kpath := range osro.KPath.CAN2TSNPaths{
		Denominator := 0.
		for kth := range kpath.Paths {
			Denominator += osro.VB.C2T_VB[nth][kth] * osro.PRM.C2T_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range kpath.Paths {
			probability := (osro.VB.C2T_VB[nth][kth] * osro.PRM.C2T_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				// if kth == 5 => arr[0,0,0,0,0,0,...,1,1,1,...,2,2,2,2,..,3,3,3,...,4,4,4,4,...] len(arr) ~ 100
				arr = append(arr, kth)
			}
		}
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t :=kpath.Paths[n]

		if nth < bg_tsn {
			bg_k_location[0] = append(bg_k_location[0], n)
			II_prime.CAN2TSNPath = append(II_prime.TSNPath, t)
		} else {
			input_k_location[0] = append(input_k_location[0], n)
			II.CAN2TSNPath = append(II.CAN2TSNPath, t)
		}
	}

	for nth, kpath := range osro.KPath.AVBPaths{
		Denominator := 0.
		for kth := range kpath.Paths {
			Denominator += osro.VB.AVB_VB[nth][kth] * osro.PRM.AVB_PRM[nth][kth]
		}

		n := 0
		var arr []int
		for kth := range kpath.Paths {
			probability := (osro.VB.AVB_VB[nth][kth] * osro.PRM.AVB_PRM[nth][kth]) / Denominator
			for j := 0; j < int(probability*100); j++ {
				// if kth == 5 => arr[0,0,0,0,0,0,...,1,1,1,...,2,2,2,2,..,3,3,3,...,4,4,4,4,...] len(arr) ~ 100
				arr = append(arr, kth)
			}
		}
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t := kpath.Paths[n]

		if nth < bg_avb {
			bg_k_location[1] = append(bg_k_location[1], n)
			II_prime.AVBPath = append(II_prime.AVBPath, t)
		} else {
			input_k_location[1] = append(input_k_location[1], n)
			II.AVBPath = append(II.AVBPath, t)
		}
	}

	return II, II_prime, input_k_location, bg_k_location
}

func epoch(network *network.Network, osro *OSRO, timeout_index int) *path.Path_set {
	II, _, input_k_location, _ := probability(osro)
	//II, II_prime, input_k_location, bg_k_location := Probability(osro.KTrees, osro.VB, osro.PRM) // BG ... pass
	fmt.Printf("Select input routing %v \n", input_k_location)
	//fmt.Printf("Select background routing %v \n", bg_k_location) // BG ... pass
	osro.Timer[timeout_index].TimerStop()
	obj_list, cost := schedule.OBJ(network, osro.KPath, II, osro.BGPath,"obo")
	//obj, cost := Obj(network, X, II, II_prime) // BG ... pass
	osro.Timer[timeout_index].TimerStart()

	if obj_list[0] == 0 && obj_list[1] == 0 {
		osro.Timer[timeout_index].TimerEnd()
	}

	for nth, kpath := range osro.KPath.TSNPaths {
		for kth := range kpath.Paths {
			if nth < bg_tsn { // BG ... pass
				//osro.PRM.TSN_PRM[nth][kth] *= osro.P
				//if kth == bg_k_location[0][nth] {
				//	osro.PRM.TSN_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				osro.PRM.TSN_PRM[nth][kth] *= osro.P
				if kth == input_k_location[0][nth-bg_tsn] {
					osro.PRM.TSN_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	for nth, kpath := range osro.KPath.CAN2TSNPaths{
		for kth := range kpath.Paths {
			if nth < bg_tsn { // BG ... pass
				//osro.PRM.TSN_PRM[nth][kth] *= osro.P
				//if kth == bg_k_location[0][nth] {
				//	osro.PRM.TSN_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				osro.PRM.C2T_PRM[nth][kth] *= osro.P
				if kth == input_k_location[0][nth-bg_tsn] {
					osro.PRM.C2T_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	for nth, kpath := range osro.KPath.AVBPaths {
		for kth := range kpath.Paths {
			if nth < bg_avb { // BG ... pass
				//osro.PRM.AVB_PRM[nth][kth] *= osro.P
				//if kth == bg_k_location[1][nth] {
				//	osro.PRM.AVB_PRM[nth][kth] += (1 / cost[3])
				//}
			} else { // Input
				osro.PRM.AVB_PRM[nth][kth] *= osro.P
				if kth == input_k_location[1][nth-bg_avb] {
					osro.PRM.AVB_PRM[nth][kth] += float64(1 / cost)
				}
			}
		}
	}

	return II
}
