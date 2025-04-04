package schedule

import (
	"fmt"
	"src/network"
	"src/network/flow"
	"src/plan/routes"
	// "src/plan/path"

	"time"
)

// Objectives
func OBJ(network *network.Network, X *routes.KTrees_set, II *routes.Trees_set, II_prime *routes.Trees_set) ([4]float64, int) {
	S := network.TSNFlow_Set.Input_TSNflow_set()
	S_prime := network.TSNFlow_Set.BG_flow_set()
	var (
		obj                [4]float64
		cost               int
		tsn_failed_count   int           = 0 // O1
		avb_failed_count   int           = 0 // O2
		all_rerouted_count int           = 0 // O3 ... pass
		avb_wcd_sum        time.Duration     // O4
	)
	linkmap := map[string]float64{}

	// Round1: Schedule BG flow
	// O1
	for nth, route := range II_prime.TSNTrees {
		schedulability := schedulability(0, S_prime.TSNFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_failed_count += 1 - schedulability
		//fmt.Printf("BackGround TSN route%d: %b \n", nth, schedulability)
	}

	//封裝 這邊要return delay,can2tsn封包
	var totalDelay float64
	importantCANFlows := network.CANFlow_Set
    for _, flow := range importantCANFlows.ImportantCANFlows{
		fmt.Printf("Source: %v ,Destinatione: %v , Datasize: %v ",flow.Source, flow.Destination, flow.DataSize)
        d, pkt := EncapsulateCAN2TSN(flow.Source, flow.Destination, flow.DataSize, flow.Deadline)
        if pkt != nil {
            fmt.Printf("封裝了一個CAN2TSN packet, datasize=%.2f, delay=%.2f\n", pkt.DataSize, d)
            totalDelay += d
        }
    }


	// O2 and O4
	for nth, route := range II_prime.AVBTrees {
		wcd := WCD(route, X, S_prime.AVBFlows[nth], network.TSNFlow_Set)
		avb_wcd_sum += wcd
		schedulability := schedulability(wcd, S_prime.AVBFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		avb_failed_count += 1 - schedulability
		//fmt.Printf("BackGround AVB route%d: %b \n", nth, schedulability)
	}
	// O3 ... pass

	
	//解封裝 WCD



	// Round2: Schedule Input flow
	// O1
	for nth, route := range II.TSNTrees {
		schedulability := schedulability(0, S.TSNFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_failed_count += 1 - schedulability
		//fmt.Printf("Input TSN route%d: %b \n", nth, schedulability)
	}

	// O2 and O4
	for nth, route := range II.AVBTrees {
		wcd := WCD(route, X, S.AVBFlows[nth], network.TSNFlow_Set)
		avb_wcd_sum += wcd
		schedulability := schedulability(wcd, S.AVBFlows[nth], route, linkmap, network.Bandwidth, network.HyperPeriod)
		avb_failed_count += 1 - schedulability
		//fmt.Printf("Input AVB route%d: %b \n", nth, schedulability)
	}
	// O3 ... pass

	obj[0] = float64(tsn_failed_count)               // O1
	obj[1] = float64(avb_failed_count)               // O2
	obj[2] = float64(all_rerouted_count)             // O3 ... pass
	obj[3] = float64(avb_wcd_sum / time.Microsecond) // O4

	cost += int(avb_wcd_sum/time.Microsecond) * 1
	cost += avb_failed_count * 1000000
	cost += tsn_failed_count * 100000000

	return obj, cost
}




func schedulability(wcd time.Duration, flow *flow.Flow, route *routes.Tree, linkmap map[string]float64, bandwidth float64, hyperPeriod int) int {
	r := wcd <= time.Duration(flow.Deadline)*time.Microsecond
	node := route.GetNodeByID(flow.Source)
	schedulable, _ := schedulable(node, -1, flow, route, linkmap, bandwidth, hyperPeriod)

	if r && schedulable {
		return 1
	}
	return 0
}

func schedulable(node *routes.Node, parentID int, flow *flow.Flow, route *routes.Tree, linkmap map[string]float64, bandwidth float64, hyperPeriod int) (bool, map[string]float64) {
	// for _, link := range node.Connections {
	// 	if link.ToNodeID == parentID {
	// 		continue

	// 	} else {
	// 		//// Duplex
	// 		if !(link.FromNodeID == flow.Source || loopcompare(link.ToNodeID, flow.Destination)) {
	// 			key := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
	// 			linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
	// 			if linkmap[key] > bandwidth {
	// 				return false, linkmap
	// 			}
	// 		}

	// 		//// Simplex
	// 		//if !(link.FromNodeID == flow.Source || loopcompare(link.ToNodeID, flow.Destinations)) {
	// 		//	key := ""
	// 		//	key1 := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
	// 		//	key2 := fmt.Sprintf("%d>%d", link.ToNodeID, link.FromNodeID)
	// 		//	if _, ok := linkmap[key1]; !ok {
	// 		//		if _, ok := linkmap[key2]; !ok {
	// 		//			key = key1
	// 		//		} else {
	// 		//			key = key2
	// 		//		}
	// 		//
	// 		//	} else {
	// 		//		key = key1
	// 		//	}
	// 		//
	// 		//	linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
	// 		//	if linkmap[key] > bandwidth {
	// 		//		return false, linkmap
	// 		//	}
	// 		//}

	// 		nextnode := route.GetNodeByID(link.ToNodeID)
	// 		schedulable, updatedLinkmap := schedulable(nextnode, node.ID, flow, route, linkmap, bandwidth, hyperPeriod)
	// 		if !schedulable {
	// 			return false, updatedLinkmap
	// 		}
	// 	}
	// }
	return true, linkmap
}

func loopcompare(a int, b []int) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}
