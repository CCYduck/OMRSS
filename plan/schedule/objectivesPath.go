package schedule

import (
	"fmt"
	"src/network"
	"src/network/flow"
	"src/plan/path"
	
	"time"
)

// Objectives
func OBJP(network *network.Network, X *path.KPath_Set, II *path.Path_set, II_prime *path.Path_set) ([4]float64, int) {
	S := network.Flow_Set.Input_TSNflow_set()
	S_prime := network.Flow_Set.BG_flow_set()
	var (
		obj                [4]float64
		cost               int
		tsn_can_failed_count   int           = 0 // O1 tsn和can的AR (掉包率)
		avb_failed_count   int           = 0 // O2 AVB的AR
		bandwidth_userate int           = 0 // O3 頻寬使用率
		wcd_sum        time.Duration     // O4 TSN 和 CAN的 Delay
	)
	linkmap := map[string]float64{}

	// Round1: Schedule BG flow
	// O1 
	for nth, path := range II_prime.TSNPath {
		schedulability := path_schedulability(0, S_prime.TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_can_failed_count += 1 - schedulability
		//fmt.Printf("BackGround TSN route%d: %b \n", nth, schedulability)
	}

	// O2 and O4
	for nth, path := range II_prime.AVBPath {
		wcd := WCDP(path, X, S_prime.AVBFlows[nth], network.Flow_Set)
		wcd_sum += wcd
		schedulability := path_schedulability(wcd, S_prime.AVBFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		avb_failed_count += 1 - schedulability
		//fmt.Printf("BackGround AVB route%d: %b \n", nth, schedulability)
	}
	// O3 ... pass
	bandwidth_userate =0

	//解封裝 WCD

	// Round2: Schedule Input flow
	// O1

	for nth, path := range II.TSNPath {
		schedulability := path_schedulability(0, S.TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_can_failed_count += 1 - schedulability
		//fmt.Printf("Input TSN route%d: %b \n", nth, schedulability)
	}

	//封裝 這邊要return delay,can2tsn封包
	
	can2tsnflow, o1_can_drop_fifo , o1_can_drop_priority, o1_can_drop_mat , delay_fifo , delay_priority , delay_mat := EncapsulateCAN2TSN(network.Flow_Set, network.HyperPeriod)
	// fmt.Printf("\n%v\n",)
	tsn_can_failed_count += o1_can_drop_fifo

	for nth, path := range II.TSNPath {
		schedulability := path_schedulability(0, S.TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_can_failed_count += 1 - schedulability
		//fmt.Printf("Input TSN route%d: %b \n", nth, schedulability)
	}

	can2tsnflow.Show_CAN2TSNFlowSet()
	fmt.Printf("O1_CAN Drop: M1: %v , M2: %v\n ,  Encapsulate Delay: %v", o1_can_drop_fifo , o1_can_drop_priority,o1_can_drop_mat , delay_fifo , delay_priority ,delay_mat)


	
	// // O2 and O4
	// for nth, path := range II.AVBPath {
	// 	wcd := WCDP(path, X, S.AVBFlows[nth], network.Flow_Set)
	// 	avb_wcd_sum += wcd
	// 	schedulability := path_schedulability(wcd, S.AVBFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
	// 	avb_failed_count += 1 - schedulability
	// 	//fmt.Printf("Input AVB route%d: %b \n", nth, schedulability)
	// }
	// O3 ... pass

	// obj[0] = float64(tsn_failed_count)               // O1
	// // obj[1] = float64(avb_failed_count)               // O2
	// obj[2] = float64(all_rerouted_count)             // O3 ... pass linkmap
	// obj[3] = float64(avb_wcd_sum / time.Microsecond) // O4

	// cost += int(avb_wcd_sum/time.Microsecond) * 1
	// cost += avb_failed_count * 1000000
	// cost += tsn_failed_count * 100000000

	return obj, cost
}

func path_schedulability(wcd time.Duration, flow *flow.Flow, path *path.Path, linkmap map[string]float64, bandwidth float64, hyperPeriod int) int {
	r := wcd <= time.Duration(flow.Deadline)*time.Microsecond
	node := path.GetNodeByID(flow.Source)
	schedulable, _ := path_schedulable(node, -1, flow, path, linkmap, bandwidth, hyperPeriod)

	if r && schedulable {
		return 1
	}
	return 0
}

func path_schedulable(node *path.Node, parentID int, flow *flow.Flow, route *path.Path, linkmap map[string]float64, bandwidth float64, hyperPeriod int) (bool, map[string]float64) {
	for _, link := range node.Connections {
		if link.ToNodeID == parentID {
			continue

		} else {
			//// Duplex
			if !(link.FromNodeID == flow.Source || path_loopcompare(link.ToNodeID, flow.Destination)) {
				key := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
				linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
				if linkmap[key] > bandwidth {
					return false, linkmap
				}
			}

			//// Simplex
			//if !(link.FromNodeID == flow.Source || loopcompare(link.ToNodeID, flow.Destinations)) {
			//	key := ""
			//	key1 := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
			//	key2 := fmt.Sprintf("%d>%d", link.ToNodeID, link.FromNodeID)
			//	if _, ok := linkmap[key1]; !ok {
			//		if _, ok := linkmap[key2]; !ok {
			//			key = key1
			//		} else {
			//			key = key2
			//		}
			//
			//	} else {
			//		key = key1
			//	}
			//
			//	linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
			//	if linkmap[key] > bandwidth {
			//		return false, linkmap
			//	}
			//}

			nextnode := route.GetNodeByID(link.ToNodeID)
			schedulable, updatedLinkmap := path_schedulable(nextnode, node.ID, flow, route, linkmap, bandwidth, hyperPeriod)
			if !schedulable {
				return false, updatedLinkmap
			}
		}
	}
	return true, linkmap
}

func path_loopcompare(a int, b int) bool {
	if a == b {
		return true
	}
	return false
}

func Testqueue(network *network.Network) {
	can2tsnflow, o1_can_drop_fifo , o1_can_drop_priority,o1_can_drop_mat , delay_fifo , delay_priority ,delay_mat := EncapsulateCAN2TSN(network.Flow_Set, network.HyperPeriod)
	// fmt.Printf("\n%v\n",)

	can2tsnflow.Show_CAN2TSNFlowSet()
	fmt.Printf("O1_CAN Drop FIFO: %v Delay %v , Priority: %v Delay: %v , MAT: %v Delay: %v \n",
	o1_can_drop_fifo , delay_fifo, o1_can_drop_priority , delay_priority , o1_can_drop_mat   ,delay_mat)
}
