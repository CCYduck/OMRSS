package plan

import (
	// "src/plan/path"
	"fmt"
	"src/plan/algo"
	// "src/plan/path"
	// "src/plan/schedule"
)

func (plan *OMACO) Initiate_Plan() {
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

}

func (plan *OSRO) Initiate_Plan() {
	//Imp50 60 70 80 Unmp 250 300 350 400	
	
	plan.SP.SP_Run(plan.Network)
	// fmt.Println(len(plan.SP.Path.TSNPath),len(plan.SP.Path.AVBPath),len(plan.SP.Path.CAN2TSNPath))
	plan.KP.KP_Run(plan.Network)
	// fmt.Println(len(plan.KP.KPath.TSNPaths),len(plan.KP.KPath.AVBPaths),len(plan.KP.KPath.CAN2TSNPaths))
	
	fmt.Println()
	fmt.Println("OSACO")
	fmt.Println("----------------------------------------")
	method:= []string{"fifo", "priority", "obo", "wat"}

	// for i, k := range plan.KP.KPath.CAN2TSNPaths {
	// 	fmt.Printf("#%02d  k=%v  Method=%q\n", i, k.K, k.Method)
	// }
	for _,m := range method{
		o := &algo.OSRO{}                               // ← 每次 new 一個
		o.OSRO_Initial_Settings(plan.Network, plan.SP.Path, plan.KP.KPath , m)
		res := o.OSRO_Run(plan.Network, 0)              // 只跑自己的 method
		fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", 
		res.Obj[0], res.Obj[1], res.Obj[2], res.Obj[3])
	}
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	
	// method:= []string{"fifo", "priority", "obo", "wat"}
	// plan.SP.Objs_SP = make([]algo.Result, 0, 4)   // 4 種 method：fifo/priority/obo/wat
	
	// for ind,m := range method{
	// 	Objs_sp, cost := schedule.OBJ(
	// 		plan.Network,
	// 		plan.KP.KPath,
	// 		plan.SP.Path.Input_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
	// 		plan.SP.Path.BG_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB),
	// 		m,
	// 	)
	// 	plan.SP.Objs_SP[ind].Obj = Objs_sp //要改SP 變成4個
	// 	plan.SP.Objs_SP[ind].Method = m
	// 	plan.SP.Objs_SP[ind].Cost =cost
	// 	fmt.Println(plan.SP.Objs_SP)
	// }


	// for _, r := range plan.SP.Objs_SP {
    // fmt.Printf("Method=%s  cost=%d  O=%v\n", r.Method, r.Cost, r.Obj)

	// }
}

//func (plan *plan3) Initiate_Plan() {
//
//}
