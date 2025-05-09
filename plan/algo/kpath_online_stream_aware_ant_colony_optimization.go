package algo

import (
	"crypto/rand"
	"fmt"
	// "fmt"
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

func (osro *OSRO) OSRO_Initial_Settings(network *network.Network, sp *path.Path_set, kp *path.KPath_Set, method string) {
	//// osro computing time: Estimate the time it takes to compute routing information

	bg_tsn = network.BG_TSN
	bg_avb = network.BG_AVB
	
	timer := algo_timer.NewTimer()
	timer.TimerStart()
	
	timer.TimerEnd()

	// osro.InputPath 	= sp.Input_Path_set(bg_tsn, bg_avb)
	// osro.InputPath.CAN2TSNPath = sp.GetPathByMethod(method)
	// osro.BGPath 	= sp.BG_Path_set(bg_tsn, bg_avb)
	// fmt.Println("sp",len(sp.TSNPath),len(sp.AVBPath),len(sp.CAN2TSNPath))
	// fmt.Println("osro",len(osro.InputPath.TSNPath),len(osro.InputPath.AVBPath),len(osro.InputPath.CAN2TSNPath))


	// osro.KPath		= kp.Input_kpath_set(bg_tsn, bg_avb)
	// osro.KPath.CAN2TSNPaths = kp.GetKPathByMethod(method)
	// osro.BGKPaths   = kp.BG_kpath_set(bg_tsn,bg_avb)

	// fmt.Println("kp",len(kp.TSNPaths),len(kp.AVBPaths),len(kp.CAN2TSNPaths))
	// fmt.Println("osro",len(osro.KPath.TSNPaths),len(osro.KPath.AVBPaths),len(osro.KPath.CAN2TSNPaths))


	/* --- Input / BG 路徑 -------------------- */
    osro.InputPath = sp.Input_Path_set(bg_tsn, bg_avb)
    osro.InputPath.CAN2TSNPath = sp.GetPathByMethod(method) // full-len (含 nil)
    osro.BGPath    = sp.BG_Path_set(bg_tsn, bg_avb)

    /* --- K-Paths --------------------------- */
    osro.KPath   = kp.Input_kpath_set(bg_tsn, bg_avb)   // full slice
    // osro.BGKPaths = kp.BG_kpath_set(bg_tsn, bg_avb)

    // 如果要快取「同 method 的 k-path」可另外存
    osro.MethodKPaths = map[string][]*path.KPath{
        method: osro.KPath.GetKPathByMethod(method),
    }
	// fmt.Println(len(osro.MethodKPaths[method]), osro.MethodKPaths[method])

	fmt.Println("sp",len(sp.TSNPath),len(sp.AVBPath),len(sp.CAN2TSNPath))
	fmt.Println("kp",len(kp.TSNPaths),len(kp.AVBPaths),len(kp.CAN2TSNPaths))
	fmt.Println("osro",len(osro.InputPath.TSNPath),len(osro.InputPath.AVBPath),len(osro.InputPath.CAN2TSNPath))
	fmt.Println("osro",len(osro.KPath.TSNPaths),len(osro.KPath.AVBPaths),len(osro.KPath.CAN2TSNPaths))


	osro.PRM = compute_prm(osro.KPath)
	osro.VB = compute_vb(osro.KPath,network.Flow_Set)

	osro.Timer[0] = algo_timer.NewTimer()
	fmt.Println(method)
	osro.Timer[0].TimerMerge(timer)
	// osro.Timer[1] = algo_timer.NewTimer()
	// osro.Timer[1].TimerMerge(timer)
	// osro.Timer[2] = algo_timer.NewTimer()
	// osro.Timer[2].TimerMerge(timer)
	// osro.Timer[3] = algo_timer.NewTimer()
	// osro.Timer[3].TimerMerge(timer)
	// osro.Timer[4] = algo_timer.NewTimer()
	// osro.Timer[4].TimerMerge(timer)
	
}

// Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"
func (osro *OSRO) OSRO_Run(network *network.Network, timeout_index int) Result {
	
	var result Result
	// 把每種 Method 各自的 InputPath 存在 map〈method → path set〉
	inputPathMap := make(map[string]*path.Path_set)

    for _, enc := range network.Flow_Set.Encapsulate {
        m := enc.Method_Name

		// kpath4m := osro.KPath.Getkpathbymethod(m)   // ★臨時切片
		// if len(kpath4m) == 0 {                      // 沒路徑就用 nil 佔位
		// 	kpath4m = make([]*path.KPath, len(enc.CAN2TSNFlows))
		// }

		// 深拷貝一份 BG＋TSN＋AVB 路徑
		// ps := &path.Path_set{}
		// ps.TSNPath = append(ps.TSNPath, osro.InputPath.TSNPath...)
		// ps.AVBPath = append(ps.AVBPath, osro.InputPath.AVBPath...)
		// ps.CAN2TSNPath = osro.InputPath.GetPathByMethod(m)

		// inputPathMap[m] = ps
		// input := ps

		timeout := time.Duration(osro.Timeout) * time.Millisecond
		start   := time.Now()

		for {
			
			II := epoch(network, osro, timeout_index,m)   // 產生 candidate
			fmt.Println(len(II.TSNPath),len(II.AVBPath),len(II.CAN2TSNPath))
			_, newCost := schedule.OBJ(network, osro.KPath, II, osro.BGPath, m)
			_, oldCost := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath, m)
	
			if newCost < oldCost {
				inputPathMap[m] = II     // ★只更新自己的
				osro.InputPath = II
			}
			if time.Since(start) >= timeout {
				break
			}
		}
	
		// ----------- 統計結果 -----------
		res,cost := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath, m)
		result.Method = m
		result.Obj    = res
		result.Cost   = cost
	
		fmt.Printf("\nresult value: %s\n", m)
		fmt.Printf("O1:%f O2:%f O3:%f O4:%f\n",
		result.Obj[0], result.Obj[1], result.Obj[2], result.Obj[3])
    }
	
	// // 6. osro
	// // Repeat the execution of epochs within the timeout
	// for _, method := range network.Flow_Set.Encapsulate {
	// 	method_name := method.Method_Name

	// 	// 1. 先把 InputPath 取出；第一次用 DeepCopy 初始化
    //     input, ok := inputPathMap[method_name]
    //     if !ok {
    //         input = osro.InputPath // 假設有 DeepCopy 方法
    //         inputPathMap[method_name] = input
    //     }

	// 	// initialobj, initialcost := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath, method_name)
	// 	// fmt.Println()
	// 	// fmt.Printf("initial value: %d \n", initialcost)
	// 	// fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", initialobj[0], initialobj[1], initialobj[2], initialobj[3])

	// 	timeout := time.Duration(osro.Timeout) * time.Millisecond
	// 	startTime := time.Now()
	// 	// i := 1
	// 	// for {
	// 	// 	fmt.Printf("\nepoch%d:\n", i)
	// 	// 	osro.Timer[timeout_index].TimerStart()
	// 	// 	II := epoch(network, osro, timeout_index)
	// 	// 	osro.Timer[timeout_index].TimerStop()

	// 	// 	_, cost1 := schedule.OBJ(network, osro.KPath, II, osro.BGPath,method.Method_Name)               	// new
	// 	// 	_, cost2 := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath,method.Method_Name) 		// old

	// 	// 	if cost1 < cost2 {
	// 	// 		osro.InputPath = II
	// 	// 		fmt.Println("Change the selected routing !!")
	// 	// 	}
	// 	// 	i += 1

	// 	// 	if time.Since(startTime) >= timeout {
	// 	// 		break
	// 	// 	}
	// 	// }
	// 	for i := 1; ; i++ {
    //         II := epoch(network, osro, timeout_index)

    //         _, newCost := schedule.OBJ(network, osro.KPath, II, osro.BGPath, method_name)
    //         _, oldCost := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath, method_name)
    //         if newCost < oldCost {
    //             inputPathMap[method_name] = II // 只更新自己的 InputPath
    //             osro.InputPath = II
    //         }
    //         if time.Since(startTime) >= timeout {
    //             break
    //         }
    //     }
	// 	// resultobj, resultcost := schedule.OBJ(network, osro.KPath, osro.InputPath, osro.BGPath,method.Method_Name)
		
	// 	// 3. 記錄最終結果
    //     resultObj, resultCost := schedule.OBJ(network, osro.KPath, input, osro.BGPath, method_name)
    //     result.Method= method_name
	// 	result.Obj=  resultObj
	// 	result.Cost=   resultCost

		
	// 	fmt.Println()
	// 	fmt.Printf("result value: %v \n", result.Method)
	// 	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", resultObj[0], resultObj[1], resultObj[2], resultObj[3])
	// 	fmt.Println()

	// 	if resultObj[0] != 0 || resultObj[1] != 0 {
	// 		osro.Timer[timeout_index].TimerMax()
	// 	}
	// }
	// fmt.Println()
	// 	fmt.Printf("result value: %d \n", result.Cost)
	// 	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", result.Obj[0], result.Obj[1], result.Obj[2], result.Obj[3])
	// 	fmt.Println()
	return result
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

	bg_c2t := 0
	for nth, kpath := range X.CAN2TSNPaths {
		var prm []float64
		for i := 0; i < len(kpath.Paths); i++ {
			if nth < bg_c2t {
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
	bg_c2t := 0
	for nth, c2t_kpath := range X.CAN2TSNPaths {
		
		var v []float64
		for kth := range c2t_kpath.Paths {
			mult := 1.
			if nth < bg_c2t && kth == 0 {
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
				value := mult / float64(schedule.WCD(z, X, flow_set.AVBFlows[nth-bg_avb], flow_set))
				v = append(v, value)

			} else {
				//fmt.Printf("Backgourd flow%d tree%d \n", nth, kth)
				value := mult / float64(schedule.WCD(z, X, flow_set.AVBFlows[nth], flow_set))
				v = append(v, value)
			}
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	return visibility
}

func probability(osro *OSRO, method string) (*path.Path_set, *path.Path_set, [2][]int, [2][]int) {
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

		// ★ 在 probability() 裡，每處 build arr 的 for 迴圈「之後」加這段
		if len(arr) == 0 {
			// -> Denominator 為 0，或是所有機率經四捨五入後都小於 0.01
			for kth := range kpath.Paths {          // 用均勻分布 fallback
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

	
	// kpaths := osro.MethodKPaths[method] 
	fmt.Println("kpaths",len(osro.MethodKPaths[method]))
	for nth, kpath := range osro.MethodKPaths[method] {
		if kpath == nil || kpath.Method != method {
			continue                                 // 跳過不屬於此 method 的項目
		}

		// fmt.Printf("kpath.Method: %v ,  method: %v \n",kpath.Method,method)
		
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

		// ★ 在 probability() 裡，每處 build arr 的 for 迴圈「之後」加這段
		if len(arr) == 0 {
			// -> Denominator 為 0，或是所有機率經四捨五入後都小於 0.01
			for kth := range kpath.Paths {          // 用均勻分布 fallback
				arr = append(arr, kth)
			}
		}

		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
		n = arr[int(randomIndex.Int64())]
		t := kpath.Paths[n]
		t.Method = kpath.Method
		// fmt.Println(kpath.Method)
		input_k_location[0] = append(input_k_location[0], n)
		II.CAN2TSNPath = append(II.CAN2TSNPath, t)

	}
	fmt.Println("II",len(II.CAN2TSNPath))

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

		// ★ 在 probability() 裡，每處 build arr 的 for 迴圈「之後」加這段
		if len(arr) == 0 {
			// -> Denominator 為 0，或是所有機率經四捨五入後都小於 0.01
			for kth := range kpath.Paths {          // 用均勻分布 fallback
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

func epoch(network *network.Network, osro *OSRO, timeout_index int, method string) *path.Path_set {
	II, _, input_k_location, _ := probability(osro,method)
	//II, II_prime, input_k_location, bg_k_location := Probability(osro.KTrees, osro.VB, osro.PRM) // BG ... pass
	// fmt.Printf("Select input routing %v \n", input_k_location)
	// fmt.Printf("Select background routing %v \n", bg_k_location) // BG ... pass
	osro.Timer[timeout_index].TimerStop()
	obj_list, cost := schedule.OBJ(network, osro.KPath, II, osro.BGPath, method)
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
	tsnN := len(osro.KPath.TSNPaths)
	for nth, kpath := range osro.KPath.CAN2TSNPaths{
		for kth := range kpath.Paths {
			// Input
			osro.PRM.C2T_PRM[nth][kth] *= osro.P
			if kth == input_k_location[0][tsnN +nth] {
				osro.PRM.C2T_PRM[nth][kth] += float64(1 / cost)
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
