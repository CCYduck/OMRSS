package schedule

import (
	"src/network/flow"
	"src/plan/path"
	"time"
)

// Worse-Case Delay
func WCDP(z *path.Path, kpath_set *path.KPath_Set, flow *flow.Flow, flow_set *flow.Flows) time.Duration {
	end2end := time.Duration(0)
	node := z.GetNodeByID(flow.Source)
	wcd := path_end2end_delay(node, -1, end2end, z, kpath_set, flow, flow_set)
	//fmt.Printf("max wcd: %v \n", wcd)

	return wcd
}

// Use DFS to find all dataflow paths in the Route
// Calculate the End to End Delay for each dataflow path and select the maximum one
func path_end2end_delay(node *path.Node, parentID int, end2end time.Duration, z *path.Path, kpath_set *path.KPath_Set, flow *flow.Flow, flow_set *flow.Flows) time.Duration {
	//fmt.Printf("%d: %v \n", node.ID, end2end)
	maxE2E := end2end
	for _, link := range node.Connections {
		per_hop := time.Duration(0)
		if link.ToNodeID == parentID {
			continue

		} else {
			// Calculation of latency for a single link
			per_hop += path_transmit_avb_itself(flow.DataSize, link.Cost)
			//per_hop += interfere_from_be(conn.Cost)
			per_hop += path_interfere_from_avb(link, kpath_set, flow.DataSize)
			per_hop += path_interfere_from_tsn(link, kpath_set, flow_set)
			end2end += per_hop

			nextnode := z.GetNodeByID(link.ToNodeID)
			nextE2E := path_end2end_delay(nextnode, node.ID, end2end, z, kpath_set, flow, flow_set)

			if maxE2E < nextE2E {
				maxE2E = nextE2E
			}
		}

		end2end -= per_hop
	}
	return maxE2E
}

// Calculate the transmission time of AVB
func path_transmit_avb_itself(datasize float64, bytes_rate float64) time.Duration {
	/// Maximum proportion of bandwidth that AVB streams can occupy.
	const MAX_AVB_SETTING float64 = 0.75
	nanoseconds := datasize * bytes_rate * MAX_AVB_SETTING
	duration := time.Duration(int64(nanoseconds))

	return duration
}

// The time occupied by a BE packet before transmission
//func interfere_from_be(bytes_rate float64) time.Duration {
//	// Maximum number of bytes in a frame.
//	const MTU float64 = 1500.
//	nanoseconds := MTU * bytes_rate
//	duration := time.Duration(int64(nanoseconds))
//
//	return duration
//}

// The time occupied by other AVB packets during transmission
func path_interfere_from_avb(link *path.Connection, KPath_set *path.KPath_Set, datasize float64) time.Duration {
	// Occupied bytes by other AVB
	var occupiedbytes float64
	for _, avb_kpath := range KPath_set.AVBPaths {
		for _, path := range avb_kpath.Paths {
			node := path.GetNodeByID(link.FromNodeID)
			if node != nil {
				for _, conn := range node.Connections {
					if conn.ToNodeID == link.ToNodeID {
						occupiedbytes += datasize
					}
				}
			}
		}
	}
	occupiedbytes -= datasize // Deducting its own datasize

	return path_transmit_avb_itself(occupiedbytes, link.Cost)
}

// The known time occupied by TSN packets during transmission
func path_interfere_from_tsn(link *path.Connection, kpath_set *path.KPath_Set, flow_set *flow.Flows) time.Duration {
	// Occupied bytes by TSN
	var occupiedbytes float64
	for nth, tsn_path := range kpath_set.TSNPaths {
		for _, path := range tsn_path.Paths {
			node := path.GetNodeByID(link.FromNodeID)
			if node != nil {
				for _, conn := range node.Connections {
					if conn.ToNodeID == link.ToNodeID {
						// occupiedbytes += datasize * (hyperPeriod / period)
						occupiedbytes += flow_set.TSNFlows[nth].DataSize *
							(float64(flow_set.TSNFlows[nth].HyperPeriod) / float64(flow_set.TSNFlows[nth].Period))
					}
				}
			}
		}
	}

	return path_transmit_avb_itself(occupiedbytes, link.Cost)
}
