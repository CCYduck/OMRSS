package schedule

// import (
// 	"fmt"
// 	"src/network"
// 	"src/network/flow"
// 	"src/plan/path"
// 	"time"
// )

// // Objectives
// func OBJP(network *network.Network, X *path.KPath_Set, II *path.Path_set, II_prime *path.Path_set, can2tsn_paths *path.Path_set) ([4]float64, int) {
// 	S := network.Flow_Set.Input_TSNflow_set()
// 	S_prime := network.Flow_Set.BG_flow_set()
// 	var (
// 		obj                  [4]float64
// 		cost                 int
// 		tsn_can_failed_count int           = 0 // O1
// 		avb_failed_count     int           = 0 // O2
// 		bandwidth_userate    int           = 0 // O3 ... pass
// 		wcd_sum              time.Duration     // O4
// 	)
// 	linkmap := map[string]float64{}

// 	// Round1: Schedule BG flow
// 	// O1
// 	for nth, path := range II_prime.TSNPath {
// 		schedulability := path_schedulability(0, S_prime.TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
// 		tsn_can_failed_count += 1 - schedulability
// 		//fmt.Printf("BackGround TSN route%d: %b \n", nth, schedulability)
// 	}

// 	// O2 and O4
// 	for nth, path := range II_prime.AVBPath {
// 		wcd := WCDP(path, X, S_prime.AVBFlows[nth], network.Flow_Set)
// 		wcd_sum += wcd
// 		schedulability := path_schedulability(wcd, S_prime.AVBFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
// 		avb_failed_count += 1 - schedulability
// 		//fmt.Printf("BackGround AVB route%d: %b \n", nth, schedulability)
// 	}

// 	// Round2: Schedule Input flow
// 	// O1
// 	for nth, path := range II.TSNPath {
// 		schedulability := path_schedulability(0, S.TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
// 		tsn_can_failed_count += 1 - schedulability
// 		//fmt.Printf("Input TSN route%d: %b \n", nth, schedulability)
// 	}

// 	for nth, path := range II.CAN2TSNPath {
// 		//fmt.Printf("BackGround TSN route%d: %b \n", nth, schedulability)
// 		for _ , method :=range S.Encapsulate{
// 			schedulability := path_schedulability(0, method.CAN2TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
// 			tsn_can_failed_count += 1 - schedulability
// 		}
// 	}
// 	// O1 and O4
// 	// tsn_can_failed_count += network.Flow_Set.CAN2TSN_O1_Drop
// 	// wcd_sum += network.Flow_Set.CAN2TSN_Delay
// 	// for nth, path := range can2tsn_paths.CAN2TSNPath {
// 	// 	schedulability := path_schedulability(0, network.Flow_Set.CAN2TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
// 	// 	tsn_can_failed_count += 1 - schedulability
// 	// 	//fmt.Printf("Input TSN route%d: %b \n", nth, schedulability)
// 	// }

// 	// O2 and O4
// 	for nth, path := range II.AVBPath {
// 		wcd := WCDP(path, X, S.AVBFlows[nth], network.Flow_Set)
// 		wcd_sum += wcd
// 		schedulability := path_schedulability(wcd, S.AVBFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
// 		avb_failed_count += 1 - schedulability
// 		//fmt.Printf("Input AVB route%d: %b \n", nth, schedulability)
// 	}

// 	// O3 bandwidth_userate

// 	obj[0] = float64(tsn_can_failed_count)       // O1
// 	obj[1] = float64(avb_failed_count)           // O2
// 	obj[2] = float64(bandwidth_userate)          // O3 ... pass
// 	obj[3] = float64(wcd_sum / time.Microsecond) // O4

// 	cost += int(wcd_sum/time.Microsecond) * 1
// 	cost += avb_failed_count * 1000000
// 	cost += tsn_can_failed_count * 100000000

// 	return obj, cost
// }

// func path_schedulability(wcd time.Duration, flow *flow.Flow, path *path.Path, linkmap map[string]float64, bandwidth float64, hyperPeriod int) int {
// 	r := wcd <= time.Duration(flow.Deadline)*time.Microsecond
// 	node := path.GetNodeByID(flow.Source)
// 	schedulable, _ := path_schedulable(node, -1, flow, path, linkmap, bandwidth, hyperPeriod)

// 	if r && schedulable {
// 		return 1
// 	}
// 	return 0
// }

// func path_schedulable(node *path.Node, parentID int, flow *flow.Flow, route *path.Path, linkmap map[string]float64, bandwidth float64, hyperPeriod int) (bool, map[string]float64) {
// 	for _, link := range node.Connections {
// 		if link.ToNodeID == parentID {
// 			continue

// 		} else {
// 			//// Duplex
// 			if !(link.FromNodeID == flow.Source || link.ToNodeID == flow.Destination) {
// 				key := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
// 				linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
// 				if linkmap[key] > bandwidth {
// 					return false, linkmap
// 				}
// 			}

// 			// Simplex
// 			if !(link.FromNodeID == flow.Source || link.ToNodeID == flow.Destination) {
// 				key := ""
// 				key1 := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
// 				key2 := fmt.Sprintf("%d>%d", link.ToNodeID, link.FromNodeID)
// 				if _, ok := linkmap[key1]; !ok {
// 					if _, ok := linkmap[key2]; !ok {
// 						key = key1
// 					} else {
// 						key = key2
// 					}
			
// 				} else {
// 					key = key1
// 				}
			
// 				linkmap[key] += flow.DataSize * float64((hyperPeriod / flow.Period))
// 				if linkmap[key] > bandwidth {
// 					return false, linkmap
// 				}
// 			}

// 			nextnode := route.GetNodeByID(link.ToNodeID)
// 			schedulable, updatedLinkmap := path_schedulable(nextnode, node.ID, flow, route, linkmap, bandwidth, hyperPeriod)
// 			if !schedulable {
// 				return false, updatedLinkmap
// 			}
// 		}
// 	}
// 	return true, linkmap
// }

