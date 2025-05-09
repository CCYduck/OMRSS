package plan

import (
	// "src/plan/path"
	"fmt"
	// "path"
	// "src/plan/algo"
	"src/network/flow"
	"src/plan/path"
	"log"
	"src/plan/schedule"
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

type MethodData struct {
    KPS      *path.KPath_Set
    PSInput  *path.Path_set
    PSBG     *path.Path_set
    Flows    []*flow.Flow
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
	method_list := []string{"fifo", "priority", "obo", "wat"}
	mds := make([]MethodData, len(method_list))

	// for i, k := range plan.KP.KPath.CAN2TSNPaths {
	// 	fmt.Printf("#%02d  k=%v  Method=%q\n", i, k.K, k.Method)
	// }
	// for _,m := range method{
	// 	o := &algo.OSRO{}                               // ← 每次 new 一個
	// 	o.OSRO_Initial_Settings(plan.Network, plan.SP.Path, plan.KP.KPath , m)
	// 	res := o.OSRO_Run(plan.Network, 0)              // 只跑自己的 method
	// 	fmt.Printf("O1: %f O2: %f O3: %f O4: %f \n", 
	// 	res.Obj[0], res.Obj[1], res.Obj[2], res.Obj[3])
	// }
	// The timeout of each run is set as 100~1000 ms (200ms, 400ms, 600ms, 800ms, 1000ms)
	
	// method:= []string{"fifo", "priority", "obo", "wat"}
	rawKPS := plan.KP.KPath
	rawPS := plan.SP.Path

	for ind,m := range method_list{

		// 1) 先拿出 KPath 及 Path （整個 struct）
		
	
		// 2) 用你改好的 Get*ByMethod 回傳一個完整的 struct
		//    假設你已經把它改成返回 *KPath_Set / *Path_set
		kpsByM := rawKPS.Getkpathbymethod(m)   
		psByM  := rawPS.Getpathbymethod(m)
	
		// 3) 切出 Input / BG 兩份
		psIn  := psByM.Input_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB)
		psBG  := psByM.BG_Path_set(plan.Network.BG_TSN, plan.Network.BG_AVB)

		// 4) 取出對應的 Flow slice，並驗證長度
		flowsByM := plan.Network.Flow_Set.FindMethod(m)
		if len(flowsByM) != len(psIn.CAN2TSNPath) {
			log.Fatalf("method %s: Flow 數 (%d) vs Path 數 (%d) 不吻合",
				m, len(flowsByM), len(psIn.CAN2TSNPath))
		}
		
		mds[ind] = MethodData{
			KPS:     kpsByM,
			PSInput: psIn,
			PSBG:    psBG,
			Flows:   flowsByM,
		}

	}
	// fmt.Println(plan.SP.Objs_SP)

	// 最後再一口氣跑算 OBJ
	for ind, m := range method_list {
		d := mds[ind]
		objs, cost := schedule.OBJ(
			plan.Network,
			d.KPS,
			d.PSInput,
			d.PSBG,
			m,
		)
		plan.SP.Objs_SP[ind].Obj    = objs
		plan.SP.Objs_SP[ind].Method = m
		plan.SP.Objs_SP[ind].Cost   = cost
	}
	fmt.Println(plan.SP.Objs_SP)

	// for _, r := range plan.SP.Objs_SP {
    // fmt.Printf("Method=%s  cost=%d  O=%v\n", r.Method, r.Cost, r.Obj)

	// }
}

//func (plan *plan3) Initiate_Plan() {
//
//}
