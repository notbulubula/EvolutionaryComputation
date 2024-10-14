package utils

// Returns number of total nodes, number of nodes to select
// solution list with initial Node and a map of visited nodes
// Use to remove duplicated code
func GetInitialState(costMatrix [][]int, startNode int) (int, int, []int, map[int]bool) {
	numNodes := len(costMatrix)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	solution := make([]int, 0, numToSelect)

	solution = append(solution, startNode)
	visited := make(map[int]bool)
	visited[startNode] = true

	return numNodes, numToSelect, solution, visited
}
