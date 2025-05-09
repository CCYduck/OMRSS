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

// 回傳節點 id slice
func (p *Path) NodeIDs() []int {
	out := make([]int, len(p.Nodes))
	for i, n := range p.Nodes { out[i] = n.ID }
	return out
}

// root prefix（含 i）
func (p *Path) PrefixIDs(i int) []int {
	out := make([]int, i+1)
	for k := 0; k <= i; k++ { out[k] = p.Nodes[k].ID }
	return out
}

func equalSlice(a, b []int) bool {
	if len(a) != len(b) { return false }
	for i := range a {
		if a[i] != b[i] { return false }
	}
	return true
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

// 回傳這條路徑從 Source → Destination 的有序連線切片
func (p *Path) GetLinks() []*Connection {
    var links []*Connection
    // IDs 裡儲存了節點走訪順序
    for i := 0; i < len(p.IDs)-1; i++ {
        fromID := p.IDs[i]
        toID   := p.IDs[i+1]
        node := p.GetNodeByID(fromID)
        if node == nil {
            continue
        }
        // 在這個節點的 Connections 找到正確的那一條
        for _, c := range node.Connections {
            if c.ToNodeID == toID {
                links = append(links, c)
                break
            }
        }
    }
    return links
}