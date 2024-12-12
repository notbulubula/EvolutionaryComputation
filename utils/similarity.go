package utils

// CommonEdges calculates the ratio of common edges between two solutions.
// Each solution is represented as a list of integers (nodes),
// and edges are implicit as pairs of consecutive nodes in the list.
func CommonEdges(solution1, solution2 []int) float64 {
	edges1 := make(map[[2]int]bool)
	edges2 := make(map[[2]int]bool)

	// Convert solution1 to a set of edges
	for i := 0; i < len(solution1); i++ {
		next := (i + 1) % len(solution1)
		edge := [2]int{solution1[i], solution1[next]}
		if edge[0] > edge[1] {
			edge[0], edge[1] = edge[1], edge[0] // Ensure edge order is consistent
		}
		edges1[edge] = true
	}

	// Convert solution2 to a set of edges
	for i := 0; i < len(solution2); i++ {
		next := (i + 1) % len(solution2)
		edge := [2]int{solution2[i], solution2[next]}
		if edge[0] > edge[1] {
			edge[0], edge[1] = edge[1], edge[0]
		}
		edges2[edge] = true
	}

	// Count common edges
	common := 0
	for edge := range edges1 {
		if edges2[edge] {
			common++
		}
	}

	// Return the number of common edges devided by length of the solution to get the similarity ratio
	return float64(common) / float64(len(solution1)-1)
}

// CommonSelectedNodes calculates the number of common selected nodes between two solutions.
// Each solution is represented as a list of integers (nodes).
func CommonSelectedNodes(solution1, solution2 []int) float64 {
	nodes1 := make(map[int]bool)
	nodes2 := make(map[int]bool)

	// Convert solution1 to a set of nodes
	for _, node := range solution1 {
		nodes1[node] = true
	}

	// Convert solution2 to a set of nodes
	for _, node := range solution2 {
		nodes2[node] = true
	}

	// Count common nodes
	common := 0
	for node := range nodes1 {
		if nodes2[node] {
			common++
		}
	}

	// Return the number of common nodes devided by the length of the solution to get the similarity ratio
	return float64(common) / float64(len(solution1))
}
