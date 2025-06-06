package routes

// Detect Cycle in a an Undirected Graph
// https://www.tutorialspoint.com/Detect-Cycle-in-a-an-Undirected-Graph
func (tree *Tree) FindCyCle() (bool, []int) {
	var cyclelist []int
	for _, node := range tree.Nodes {
		visited := make(map[int]bool)
		if DFSCyCle(tree, node, visited, -1, node.ID) {
			cyclelist = append(cyclelist, node.ID)
		}
	}
	if len(cyclelist) != 0 {
		return true, cyclelist
	} else {
		return false, cyclelist
	}
}

func DFSCyCle(tree *Tree, node *Node, visited map[int]bool, parentID int, startID int) bool {
	visited[node.ID] = true
	for _, conn := range node.Connections {
		if conn.ToNodeID == parentID {
			continue
		}
		if conn.ToNodeID == startID {
			return true
		}
		if visited[conn.ToNodeID] {
			continue
		}
		toNode := tree.GetNodeByID(conn.ToNodeID)
		if DFSCyCle(tree, toNode, visited, node.ID, startID) {
			return true
		}
	}
	return false
}

// How to find feedback edge set in undirected graph
// https://stackoverflow.com/questions/10791689/how-to-find-feedback-edge-set-in-undirected-graph
func (tree *Tree) GetFeedbackEdgeSet(cyclelist []int, E []int) [][2]int {
	var E_prime [][2]int
	for _, cycle := range cyclelist {
		node := tree.GetNodeByID(cycle)
		for _, conn := range node.Connections {
			if !(InCycleList(cyclelist, conn.ToNodeID)) {
				continue
			} else {
				var nodeconn [2]int
				nodeconn[0] = conn.FromNodeID
				nodeconn[1] = conn.ToNodeID
				if InE(E, nodeconn) {
					continue
				}
				if !(InEPrime(E_prime, nodeconn)) {
					E_prime = append(E_prime, nodeconn)
				}
			}
		}
	}

	return E_prime
}

func InCycleList(cyclelist []int, ToNodeID int) bool {
	for _, id := range cyclelist {
		if id == ToNodeID {
			return true
		}
	}
	return false
}

func InE(E []int, nodeconn [2]int) bool {
	for index, id := range E {
		if id == nodeconn[0] {
			if E[index+1] == nodeconn[1] {
				return true
			}
			if E[index-1] == nodeconn[1] {
				return true
			}
		}
	}
	return false
}

func InEPrime(E_prime [][2]int, nodeconn [2]int) bool {
	for _, e_prime := range E_prime {
		if e_prime[0] == nodeconn[0] && e_prime[1] == nodeconn[1] {
			return true
		}
		if e_prime[0] == nodeconn[1] && e_prime[1] == nodeconn[0] {
			return true
		}
	}
	return false
}
