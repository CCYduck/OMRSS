package osaco

import (
	"src/flow"
	"src/routes"
	"time"
)

func CompVB(X *routes.KTrees_set, flow_set *flow.Flows) *Visibility {
	BG_trees_set := X.BG_SteinerTree_set()
	Input_flow_set := flow_set.Input_flow_set()
	BG_flow_set := flow_set.BG_flow_set()

	var (
		preference   float64 = 2.
		bg_tsn_start int     = len(X.TSNTrees) / 2
		bg_avb_start int     = len(X.AVBTrees) / 2
	)

	visibility := &Visibility{}
	// OSACO CompVB line 9
	// TSN flow

	for nth, tsn_ktree := range X.TSNTrees {
		var v []float64
		for kth := range tsn_ktree.Trees {
			mult := 1.
			if nth >= bg_tsn_start && kth == 0 {
				mult = preference
			}

			value := mult / float64(tsn_ktree.Trees[0].Weight) // mult / SteinerTree weight
			v = append(v, value)
		}
		visibility.TSN_VB = append(visibility.TSN_VB, v)
	}

	// OSACO CompVB line 11
	// AVB flow
	for nth, avb_ktree := range X.AVBTrees {
		var v []float64
		for kth, z := range avb_ktree.Trees {
			mult := 1.
			if nth >= bg_avb_start && kth == 0 {
				mult = preference
			}

			if nth < bg_avb_start {
				//fmt.Printf("Input flow%d tree%d \n", nth, kth)
				value := mult / WCD(z, BG_trees_set, Input_flow_set.AVBFlows[nth])
				v = append(v, value)

			} else {
				//fmt.Printf("Backgourd flow%d tree%d \n", nth, kth)
				value := mult / WCD(z, BG_trees_set, BG_flow_set.AVBFlows[nth-bg_avb_start])
				v = append(v, value)
			}
		}
		visibility.AVB_VB = append(visibility.AVB_VB, v)
	}

	return visibility
}

func WCD(z *routes.Tree, BG_trees_set *routes.Trees_set, flow *flow.Flow) float64 {
	end2end := time.Duration(0)
	node := z.GetNodeByID(flow.Source)
	wcd := EndtoEndDelay(node, -1, end2end, z, BG_trees_set, flow)
	//fmt.Printf("max wcd: %v \n", wcd)

	return float64(wcd)
}

// Use DFS to find all dataflow paths in the Route
// Calculate the End to End Delay for each dataflow path and select the maximum one
func EndtoEndDelay(node *routes.Node, parentID int, end2end time.Duration, z *routes.Tree, BG_trees_set *routes.Trees_set, flow *flow.Flow) time.Duration {
	//fmt.Printf("%d: %v \n", node.ID, end2end)
	maxE2E := end2end
	for _, link := range node.Connections {
		per_hop := time.Duration(0)
		if link.ToNodeID == parentID {
			continue

		} else {
			// Calculation of latency for a single link
			per_hop += transmit_avb_itself(flow.Streams[0].DataSize, link.Cost)
			//per_hop += interfere_from_be(conn.Cost)
			per_hop += interfere_from_avb(link, BG_trees_set, flow.Streams[0].DataSize)
			per_hop += interfere_from_tsn(per_hop, link, flow)
			end2end += per_hop

			nextnode := z.GetNodeByID(link.ToNodeID)
			nextE2E := EndtoEndDelay(nextnode, node.ID, end2end, z, BG_trees_set, flow)

			if maxE2E < nextE2E {
				maxE2E = nextE2E
			}
		}

		end2end -= per_hop
	}
	return maxE2E
}

// Calculate the transmission time of AVB
func transmit_avb_itself(datasize float64, bytes_rate float64) time.Duration {
	/// Maximum proportion of bandwidth that AVB streams can occupy.
	const MAX_AVB_SETTING float64 = 0.75
	nanoseconds := datasize * bytes_rate * MAX_AVB_SETTING
	duration := time.Duration(int64(nanoseconds))

	return duration
}

// The time occupied by a BE packet before transmission
func interfere_from_be(bytes_rate float64) time.Duration {
	// Maximum number of bytes in a frame.
	const MTU float64 = 1500.
	nanoseconds := MTU * bytes_rate
	duration := time.Duration(int64(nanoseconds))

	return duration
}

// The time occupied by other AVB packets during transmission
func interfere_from_avb(link *routes.Connection, BG_trees_set *routes.Trees_set, datasize float64) time.Duration {
	// Occupied bytes by other AVB
	var occupiedbytes float64
	for _, tree := range BG_trees_set.AVBTrees {
		node := tree.GetNodeByID(link.FromNodeID)
		if node != nil {
			for _, conn := range node.Connections {
				if conn.ToNodeID == link.ToNodeID {
					occupiedbytes += datasize
				}
			}
		}
	}
	occupiedbytes -= datasize // Deducting its own datasize

	return transmit_avb_itself(occupiedbytes, link.Cost)
}

// Sune Mølgaard Laursen, Paul Pop, Wilfried Steiner, "Routing Optimization of AVB Streams in TSN Networks"
// The known time occupied by TSN packets during transmission
func interfere_from_tsn(twc time.Duration, link *routes.Connection, linkflow *flow.Flow) time.Duration {

	return time.Duration(0)
}
