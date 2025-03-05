package path

// Verify the existence of a Node in a Tree using its ID
func (path *Path) CheckNodeByID(id int) (*Node, bool) {
	for _, node := range path.Nodes {
		if node.ID == id {
			return node, true
		}
	}
	return &Node{ID: id}, false
}

// Find the nodes in the tree by id
func (path *Path) GetNodeByID(id int) *Node {
	for _, node := range path.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

func Compare_Connections(conn1, conn2 []*Connection) bool {
	i := 0
	for _, c1 := range conn1 {
		for _, c2 := range conn2 {
			if c2.ToNodeID == c1.ToNodeID {
				i += 1
			}
		}
	}

	if i == len(conn1) {
		return true
	} else {
		return false
	}
}

func loopcompare_simplex(a int, b []int) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

func loopcompare_complex(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}