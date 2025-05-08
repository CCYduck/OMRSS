package path

import (
	// "fmt"
	"src/network"
)

func (path_set *Path_set) Input_Path_set(bg_tsn_end int, bg_avb_end int) *Path_set {
	Input_path_set := new_Path_Set()


	Input_path_set.TSNPath= append(Input_path_set.TSNPath, path_set.TSNPath[bg_tsn_end:]...)
	Input_path_set.AVBPath = append(Input_path_set.AVBPath, path_set.AVBPath[bg_avb_end:]...)
	Input_path_set.CAN2TSNPath =append(Input_path_set.CAN2TSNPath,path_set.CAN2TSNPath... )

	
	return Input_path_set
}

func (path_set *Path_set) BG_Path_set(bg_tsn_end int, bg_avb_end int) *Path_set {
	BG_path_set := new_Path_Set()

	BG_path_set .TSNPath = append(BG_path_set .TSNPath , path_set.TSNPath[:bg_tsn_end]...)
	BG_path_set .AVBPath = append(BG_path_set .AVBPath, path_set.AVBPath[:bg_avb_end]...)

	return BG_path_set 
}


func Get_OSRO_Routing(Network *network.Network) *KPath_Set {
	const k = 3 
	kpath_set := new_KPath_Set()

	for idx, flow := range Network.Flow_Set.TSNFlows {
		topo := Network.Graph_Set.TSNGraphs[idx]

		kp   := BuildKPath(k, flow.Source, flow.Destination, topo)
		kpath_set.TSNPaths = append(kpath_set.TSNPaths, kp)
	}
	
	// -------- AVB ----------
	for idx, flow := range Network.Flow_Set.AVBFlows {
		topo := Network.Graph_Set.AVBGraphs[idx]

		kp   := BuildKPath(k, flow.Source, flow.Destination, topo)
		kpath_set.AVBPaths = append(kpath_set.AVBPaths, kp)
	}
	
	// -------- CAN→TSN (封裝流) ----------
	for _, m := range Network.Flow_Set.Encapsulate {   // 每種封裝方法
		for _, f := range m.CAN2TSNFlows {             // 每條 CAN→TSN flow
			topo := Network.Graph_Set.GetGarphBySD(f.Source, f.Destination)

			kp   := BuildKPath(k, f.Source, f.Destination, topo)		
			kpath_set.CAN2TSNPaths = append(kpath_set.CAN2TSNPaths, kp)
		}
	}
	
	return kpath_set
	
}


func (kpath_set *KPath_Set) Input_kpath_set(bg_tsn_end int, bg_avb_end int) *KPath_Set {
	Input_kpath_set :=  new_KPath_Set()

	Input_kpath_set.TSNPaths = append(Input_kpath_set.TSNPaths, kpath_set.TSNPaths[bg_tsn_end:]...)
	Input_kpath_set.AVBPaths = append(Input_kpath_set.AVBPaths, kpath_set.AVBPaths[bg_avb_end:]...)
	Input_kpath_set.CAN2TSNPaths =append(Input_kpath_set.CAN2TSNPaths,kpath_set.CAN2TSNPaths... )

	return Input_kpath_set
}

func (kpath_set *KPath_Set) BG_kpath_set(bg_tsn_end int, bg_avb_end int) *KPath_Set {
	BG_kpath_set := new_KPath_Set()

	BG_kpath_set.TSNPaths = append(BG_kpath_set.TSNPaths, kpath_set.TSNPaths[:bg_tsn_end]...)
	BG_kpath_set.AVBPaths = append(BG_kpath_set.AVBPaths, kpath_set.AVBPaths[:bg_avb_end]...)

	return BG_kpath_set
}

func (path_set *Path_set)Getpathbymethod(method string) []*Path{
	method_path_set := new_Path_Set()
	for _,path := range path_set.CAN2TSNPath{
		// fmt.Println(path.Method,method)
		if path.Method == method {
			method_path_set.CAN2TSNPath =append(method_path_set.CAN2TSNPath, path)
		}
	}
	return method_path_set.CAN2TSNPath
}