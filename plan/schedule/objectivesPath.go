package schedule

import (
	"fmt"
	"src/network"
	"src/network/flow"
	"src/plan/path"
	"time"
)

// Objectives
func OBJ(network *network.Network, X *path.KPath_Set, II *path.Path_set, II_prime *path.Path_set,m string) ([4]float64, int) {

	S := network.Flow_Set.Input_TSNflow_set()
	S_prime := network.Flow_Set.BG_flow_set()

	fmt.Println(len(II.TSNPath),len(II.AVBPath),len(II.CAN2TSNPath),len(II_prime.TSNPath),len(II_prime.AVBPath))
	// fmt.Println(len(S.TSNFlows),len(S.AVBFlows),len(S.Encapsulate[0].CAN2TSNFlows),len(S_prime.TSNFlows),len(S_prime.AVBFlows))
	var (
		obj                  		[4]float64
		cost                 		int
		tsn_can_failed_count		int           = 0
		avb_failed_count     		int           = 0 // O2
		bandwidth_userate    		int           = 0 // O3 ... pass
		wcd_sum              		time.Duration     // O4
	)
	linkmap := map[string]float64{}

	// Round1: Schedule BG flow
	// O1
	for nth, path := range II_prime.TSNPath {
		schedulability := schedulability(0, S_prime.TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_can_failed_count += 1 - schedulability
		//fmt.Printf("BackGround TSN route%d: %b \n", nth, schedulability)
	}

	// O2 and O4
	for nth, path := range II_prime.AVBPath {
		wcd := WCD(path, X, S_prime.AVBFlows[nth], network.Flow_Set)
		wcd_sum += wcd
		schedulability := schedulability(wcd, S_prime.AVBFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		avb_failed_count += 1 - schedulability
		//fmt.Printf("BackGround AVB route%d: %b \n", nth, schedulability)
	}

	// Round2: Schedule Input flow
	// O1
	for nth, path := range II.TSNPath {
		schedulability := schedulability(0, S.TSNFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_can_failed_count += 1 - schedulability
		//fmt.Printf("Input TSN route%d: %b \n", nth, schedulability)
	}
	
	method_path := II.Getpathbymethod(m)
	fmt.Println(len(method_path))
	for nth, path := range method_path {
		// fmt.Printf("BackGround TSN route%d: %b \n", nth, schedulability)
		// fmt.Println(len(S.Encapsulate))
		// fmt.Println()
		schedulability := schedulability(0, S.FindMethod(m)[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		tsn_can_failed_count += 1 - schedulability
	}

	// O2 and O4
	for nth, path := range II.AVBPath {
		wcd := WCD(path, X, S.AVBFlows[nth], network.Flow_Set)
		wcd_sum += wcd
		schedulability := schedulability(wcd, S.AVBFlows[nth], path, linkmap, network.Bandwidth, network.HyperPeriod)
		avb_failed_count += 1 - schedulability
		//fmt.Printf("Input AVB route%d: %b \n", nth, schedulability)
	}

	// O3 bandwidth_userate
	totalBytes := 0
	for _, used := range linkmap {
		totalBytes += int(used)           // linkmap 存的是 bytes
	}


	obj[0] = float64(tsn_can_failed_count)       // O1
	obj[1] = float64(avb_failed_count)           // O2
	obj[2] = float64(bandwidth_userate)          // O3 ... pass
	obj[3] = float64(wcd_sum / time.Microsecond) // O4

	cost += int(wcd_sum/time.Microsecond) * 1
	cost += avb_failed_count * 1000000
	cost += tsn_can_failed_count * 100000000
	fmt.Println(linkmap)
	return obj, cost
}

func schedulability(wcd time.Duration, flow *flow.Flow, path *path.Path, linkmap map[string]float64, bandwidth float64, hyperPeriod int) int {
	r := wcd <= time.Duration(flow.Deadline)*time.Microsecond
	node := path.GetNodeByID(flow.Source)
	schedulable, _ := schedulable(node, -1, flow, path, linkmap, bandwidth, hyperPeriod)

	if r && schedulable {
		return 1
	}
	return 0
}

func schedulable(node *path.Node, parentID int, flow *flow.Flow, route *path.Path, linkmap map[string]float64, bandwidth float64, hyperPeriod int) (bool, map[string]float64) {
	for _, link := range node.Connections {
		if link.ToNodeID == parentID {
			continue

		} else {
			// // Duplex
			// if !(link.FromNodeID == flow.Source || link.ToNodeID == flow.Destination) {
			// 	key := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
			// 	for _ ,stream := range flow.Streams{
			// 		linkmap[key] += stream.DataSize
			// 	} 
			// 	if linkmap[key] > bandwidth {
			// 		return false, linkmap
			// 	}
			// }
			
			// Simplex
			if !(link.FromNodeID == flow.Source || link.ToNodeID == flow.Destination) {
				key := ""
				key1 := fmt.Sprintf("%d>%d", link.FromNodeID, link.ToNodeID)
				key2 := fmt.Sprintf("%d>%d", link.ToNodeID, link.FromNodeID)
				if _, ok := linkmap[key1]; !ok {
					if _, ok := linkmap[key2]; !ok {
						key = key1
					} else {
						key = key2
					}
			
				} else {
					key = key1
				}
			
				for _ ,stream := range flow.Streams{
					linkmap[key] += stream.DataSize
				} 
				if linkmap[key] > bandwidth {
					return false, linkmap
				}
			}

			nextnode := route.GetNodeByID(link.ToNodeID)
			schedulable, updatedLinkmap := schedulable(nextnode, node.ID, flow, route, linkmap, bandwidth, hyperPeriod)
			if !schedulable {
				return false, updatedLinkmap
			}
		}
	}
	return true, linkmap
}


