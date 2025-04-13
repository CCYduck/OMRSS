package path

import (
	// "fmt"
	"src/network"
)

func (path_set *Path_set) Input_Path_set(bg_tsn_end int, bg_avb_end int) *Path_set {
	Input_path_set := new_Path_Set()

	Input_path_set.TSNPath= append(Input_path_set.TSNPath, Input_path_set.TSNPath[bg_tsn_end:]...)
	Input_path_set.AVBPath = append(Input_path_set.AVBPath, Input_path_set.AVBPath[bg_avb_end:]...)

	return Input_path_set
}

func (path_set *Path_set) BG_Path_set(bg_tsn_end int, bg_avb_end int) *Path_set {
	BG_path_set := new_Path_Set()

	BG_path_set .TSNPath = append(BG_path_set .TSNPath , path_set.TSNPath[:bg_tsn_end]...)
	BG_path_set .AVBPath = append(BG_path_set .AVBPath, path_set.AVBPath[:bg_avb_end]...)

	return BG_path_set 
}

func Get_OSACO_Routing(network *network.Network, SP *Path_set, K int, Method_Number int) *KPath_Set {
	kpath_set := new_KPath_Set()

	// for nth, flow := range network.TSNFlow_Set.TSNFlows {
	// 	Ktrees := KSpanningTree(v2v, SMT.TSNTrees[nth], K, flow.Source, flow.Destination, network.BytesRate, Method_Number)
	// 	ktrees_set.TSNTrees = append(ktrees_set.TSNTrees, Ktrees)
	// }
	// fmt.Printf("Finish OSACO %d TSN streams routing\n", len(ktrees_set.TSNTrees))

	// for nth, flow := range network.TSNFlow_Set.AVBFlows {
	// 	Ktrees := KSpanningTree(v2v, SMT.AVBTrees[nth], K, flow.Source, flow.Destination, network.BytesRate, Method_Number)
	// 	ktrees_set.AVBTrees = append(ktrees_set.AVBTrees, Ktrees)
	// }
	// fmt.Printf("Finish OSACO %d AVB streams routing\n", len(ktrees_set.AVBTrees))

	return kpath_set
}

func (kpath_set *KPath_Set) Input_kpath_set(bg_tsn_end int, bg_avb_end int) *KPath_Set {
	Input_kpath_set :=  new_KPath_Set()

	Input_kpath_set.TSNPaths = append(Input_kpath_set.TSNPaths, kpath_set.TSNPaths[bg_tsn_end:]...)
	Input_kpath_set.AVBPaths = append(Input_kpath_set.AVBPaths, kpath_set.AVBPaths[bg_tsn_end:]...)

	return Input_kpath_set
}

func (kpath_set *KPath_Set) BG_kpath_set(bg_tsn_end int, bg_avb_end int) *KPath_Set {
	BG_kpath_set := new_KPath_Set()

	BG_kpath_set.TSNPaths = append(BG_kpath_set.TSNPaths, kpath_set.TSNPaths[:bg_tsn_end]...)
	BG_kpath_set.AVBPaths = append(BG_kpath_set.TSNPaths, kpath_set.AVBPaths[:bg_tsn_end]...)

	return BG_kpath_set
}
