package plan

import (
	// "src/plan/path"
	// "fmt"
	// "src/plan/algo"
	// "src/plan/schedule"
)

// func (plan *OMACO) Initiate_Plan() {
// algo run
// fmt.Println("Steiner Tree")
// fmt.Println("----------------------------------------")
// plan.SMT.SMT_Run(plan.Network)

// fmt.Println()
// fmt.Println("MDTC")
// fmt.Println("----------------------------------------")
// plan.MDTC.MDTC_Run(plan.Network)

// fmt.Println()
// fmt.Println("OSACO")
// fmt.Println("----------------------------------------")
// plan.OSACO.OSACO_Initial_Settings(plan.Network, plan.SMT.Trees)
// // The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
// for i := 0; i < 5; i++ {
// 	plan.OSACO.Objs_osaco[i] = plan.OSACO.OSACO_Run(plan.Network, i)
// }

// fmt.Println()
// fmt.Println("OSACO_IAS")
// fmt.Println("----------------------------------------")
// plan.OSACO_IAS.OSACO_Initial_Settings(plan.Network, plan.SMT.Trees)
// // The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
// for i := 0; i < 5; i++ {
// 	plan.OSACO_IAS.Objs_osaco[i] = plan.OSACO_IAS.OSACO_Run(plan.Network, i)
// }

// obj_smt, _ := schedule.OBJ(
// 	plan.Network,
// 	plan.OSACO.KTrees,
// 	plan.SMT.Trees.Input_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
// 	plan.SMT.Trees.BG_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
// )

// obj_mdt, _ := schedule.OBJ(
// 	plan.Network,
// 	plan.OSACO.KTrees,
// 	plan.MDTC.Trees.Input_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
// 	plan.MDTC.Trees.BG_Tree_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
// )

// plan.SMT.Objs_smt = obj_smt
// plan.MDTC.Objs_mdtc = obj_mdt

// if obj_mdt[0] != 0 || obj_mdt[1] != 0 {
// 	plan.MDTC.Timer.TimerMax()
// }

// }

func (plan *OSRO) Initiate_Plan() {
	//Imp50 60 70 80 Unmp 250 300 350 400

	// schedule.Testqueue(plan.Network)
	
	// fmt.Println("KPath")
	// fmt.Println("----------------------------------------")
	// plan.KP.KP_Run(plan.Network)

	// // fmt.Println()
	// // fmt.Println("MDTC")
	// fmt.Println("----------------------------------------")
	// // plan.MDTC.MDTC_Run(plan.Network)
	
	// fmt.Println(len(plan.SP.Path.TSNPath),len(plan.SP.Path.AVBPath),len(plan.SP.Path.CAN2TSNPath))
	// path_set := path.BestPath(plan.Network)
	// path_set.Show_Path_Set()
	
	plan.SP.SP_Run(plan.Network)
	// plan.KP.KP_Run(plan.Network)
	// fmt.Println()
	// fmt.Println("OSACO")
	// fmt.Println("----------------------------------------")
	// method:= []string{"fifo", "priority", "obo", "wat"}
	// for ind,m := range method{
	// 	plan.OSRO.OSRO_Initial_Settings(plan.Network, plan.SP.Path, m)
	// 	plan.OSRO.OSRO_Run(plan.Network,0)
	// 	fmt.Printf("result value: %v \n", m)
	// 	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", 
	// 	plan.OSRO.Objs_osro[ind].Obj[0], plan.OSRO.Objs_osro[ind].Obj[1], plan.OSRO.Objs_osro[ind].Obj[2], plan.OSRO.Objs_osro[ind].Obj[3])
	// 	fmt.Println()
	// }
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	// plan.OSRO.OSRO_Run(plan.Network,0)
	// fmt.Println(len(plan.OSRO.Objs_osro))
	// for ind,m := range plan.OSRO.Objs_osro[0].Method{
	// 	fmt.Printf("result value: %v \n", m)
	// 	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", 
	// 	plan.OSRO.Objs_osro[ind].Obj[0], plan.OSRO.Objs_osro[ind].Obj[1], plan.OSRO.Objs_osro[ind].Obj[2], plan.OSRO.Objs_osro[ind].Obj[3])
	// 	fmt.Println()
	// }
	// fmt.Printf("result value: %v \n", plan.OSRO.Objs_osro[0].Method)
	// fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", 
	// plan.OSRO.Objs_osro[0].Obj[0], plan.OSRO.Objs_osro[0].Obj[1], plan.OSRO.Objs_osro[0].Obj[2], plan.OSRO.Objs_osro[0].Obj[3])
	// fmt.Println()
	// for i := 0; i < 5; i++ {
	// 	plan.OSRO.Objs_osro[i] = plan.OSRO.OSRO_Run(plan.Network,i)
	// 	fmt.Println()
	// 	fmt.Printf("result value: %d \n", plan.OSRO.Objs_osro[i].Cost)
	// 	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", 
	// 	plan.OSRO.Objs_osro[i].Obj[0], plan.OSRO.Objs_osro[i].Obj[1], plan.OSRO.Objs_osro[i].Obj[2], plan.OSRO.Objs_osro[i].Obj[3])
	// 	fmt.Println()
	// }

	// plan.SP.SP_Run(plan.Network)
	// plan.KP.KP_Run(plan.Network)
	
	// method:= []string{"fifo", "priority", "obo", "wat"}
	// // plan.SP.Objs_SP = make([]*algo.Result, 0, 4)   // 4 種 method：fifo/priority/obo/wat
	// // // fmt.Println(len(plan.SP.Path.TSNPath), len(plan.SP.Path.Input_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB).TSNPath), len(plan.SP.Path.BG_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB).TSNPath))
	// // plan.SP.InputPath = plan.SP.Path.Input_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB)
	// // plan.SP.BGPath = plan.SP.Path.BG_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB)
	// for ind,m := range method{
	// 	plan.OSRO_method[ind].OSRO_Initial_Settings(plan.Network, plan.SP.Path, m)

	// 	Objs_sp, cost := schedule.OBJ(
	// 		plan.Network,
	// 		plan.OSRO_method[ind].KPath,
	// 		plan.OSRO_method[ind].InputPath ,
	// 		plan.OSRO_method[ind].BGPath ,
	// 		m,
	// 	)
	// 	plan.OSRO_method[ind].Objs_osro = append(plan.OSRO_method[ind].Objs_osro, &algo.Result{
	// 		Obj:    Objs_sp,
	// 		Method: m,
	// 		Cost:   cost,
	// 	})

	// 	fmt.Printf(" %v : O1: %f O2: %f O3: %f O4: %f \n",m , Objs_sp[0], Objs_sp[1], Objs_sp[2], Objs_sp[3])

	// 	// plan.SP.Objs_SP=append(plan.SP.Objs_SP, result)//要改SP 變成4個
	// 	// fmt.Printf("method=%s obj=%v\n", m, Objs_sp)
	// }
	// file := "osro_history-safe_deadline guard--important_can 100 --unimportant_can 500 --input_tsn 30 --input_avb 70 --bg_tsn 30 --bg_avb 18 .xlsx"
	// var all []*algo.Result
	// for _, osro := range plan.OSRO_method {
	// 	all = append(all, osro.Objs_osro...) // 每種方法可能 append 多筆 epoch 結果
	// }
	// algo.SaveOSROExcel(file, all)
	// fmt.Println("Results appended to", file)
}

//func (plan *plan3) Initiate_Plan() {
//
//}
