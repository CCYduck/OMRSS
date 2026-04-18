package plan

import (
	// "src/plan/path"
	"fmt"
	"src/plan/algo"
	"src/plan/schedule"

)

// import (
// 	// "src/plan/path"
// 	"src/plan/algo"
// )

func (plan *OSRO) Initiate_Plan()(sp []*algo.Result, kp []*algo.Result) {
	// //Imp50 60 70 80 Unmp 250 300 350 400

	// // schedule.Testqueue(plan.Network)
	plan.SP.SP_Run(plan.Network)
	plan.KP.KP_Run(plan.Network)
	
	// // plan.Network.Graph_Set.Show_Graphs()


	fmt.Println("Shortest Path")
	fmt.Println("----------------------------------------")
	method:= []string{"fifo", "priority", "obo", "wst", "mao"}
	plan.SP.Objs_SP = make([]*algo.Result, 0, len(method))   // 4 種 method：fifo/priority/obo/wat
	// fmt.Println(len(plan.SP.Path.TSNPath), len(plan.SP.Path.Input_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB).TSNPath), len(plan.SP.Path.BG_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB).TSNPath))
	plan.SP.InputPath = plan.SP.Path.Input_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB)
	plan.SP.BGPath = plan.SP.Path.BG_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB)
	
	for _,m := range method{
		plan.SP.InputPath.CAN2TSNPath = plan.SP.Path.Getpathbymethod(m)
		plan.KP.KPath.CAN2TSNPaths =plan.KP.KPath.Getkpathbymethod(m)
		// plan.OSRO_method[ind].OSRO_Run(plan.Network, 0, ind, m)

		Objs_sp, cost := schedule.OBJ(
			plan.Network, 
			plan.KP.KPath, 
			plan.SP.InputPath ,
			plan.SP.BGPath,
			m,
		)

		plan.SP.Objs_SP = append(plan.SP.Objs_SP, &algo.Result{
			Obj:    Objs_sp,
			Method: m,
			Cost:   cost,
		})

		fmt.Printf("%v : O1: %f O2: %f O3: %f O4: %f Cost: %v \n",m , Objs_sp[0], Objs_sp[1], Objs_sp[2], Objs_sp[3], cost)

		// plan.SP.Objs_SP=append(plan.SP.Objs_SP, result)//要改SP 變成4個
		// fmt.Printf("method=%s obj=%v\n", m, Objs_sp)
	}
	// file_sp := "sp_history0629 --important_can 30 --unimportant_can 150 --input_tsn 9 --input_avb 21.xlsx"
	//fmt.Println()
	var all_sp []*algo.Result

	for _, sp := range plan.SP.Objs_SP {
		all_sp = append(all_sp, sp) // 每種方法可能 append 多筆 epoch 結果
	}
	// algo.SaveOSROExcel(file_sp, all_sp)
	// fmt.Println("Results appended to", file_sp)

	fmt.Println()
	fmt.Println("OSRO")
	fmt.Println("----------------------------------------")
	plan.KP.Objs_kp = make([]*algo.Result, 0, len(method))   // 4 種 method：fifo/priority/obo/wat
	for ind,m := range method{
		
		plan.OSRO_method[ind].OSRO_Initial_Settings(plan.Network, plan.SP.Path, m)
		// fmt.Printf("%v \n", m)
		plan.OSRO_method[ind].OSRO_Run(plan.Network, 0, ind, m)
		// fmt.Println(plan.Network.Flow_Set.Encapsulate[ind].Method_Name,m)
		
		Objs_kp, cost := schedule.OBJ(
			plan.Network,
			plan.OSRO_method[ind].KPath,
			plan.OSRO_method[ind].InputPath ,
			plan.OSRO_method[ind].BGPath ,
			m,
		)

		plan.OSRO_method[ind].Objs_osro = append(plan.OSRO_method[ind].Objs_osro, &algo.Result{
			Obj:    Objs_kp,
			Method: m,
			Cost:   cost,
		})

		//fmt.Printf("%v : O1: %f O2: %f O3: %f O4: %f Cost: %v \n",m , Objs_kp[0], Objs_kp[1], Objs_kp[2], Objs_kp[3], cost)
		fmt.Printf("%v : O1: %f O2: %f O3: %f O4: %f  \n",m , Objs_kp[0], Objs_kp[1], Objs_kp[2], Objs_kp[3])
		// plan.SP.Objs_SP=append(plan.SP.Objs_SP, result)//要改SP 變成4個
		// fmt.Printf("method=%s obj=%v\n", m, Objs_sp)
	}
	// file_osro := "osro_history0629 --important_can 30 --unimportant_can 150 --input_tsn 9 --input_avb 21.xlsx"
	var all []*algo.Result
	for _, osro := range plan.OSRO_method {
		all = append(all, osro.Objs_osro...) // 每種方法可能 append 多筆 epoch 結果
	}
	// algo.SaveOSROExcel(file_osro, all)
	// fmt.Println("Results appended to", file_osro)
	return plan.SP.Objs_SP, all
}

// func (plan *OSRO) Initiate_Plan()(sp []*algo.Result, kp []*algo.Result)  {
// 	return sp, kp
// }
