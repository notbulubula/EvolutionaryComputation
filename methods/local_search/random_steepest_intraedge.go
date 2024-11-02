package local_search

import (
	"evolutionary_computation/methods"
)

func RandomSteepestIntraEdge(distanceMatrix [][]int, startNode int) []int {
	// Run random function to get the initial solution
	selectedIDs := methods.RandomSolution(distanceMatrix, startNode)
	visted := make(map[int]bool)
	// Make visited map
	for _, id := range selectedIDs {
		visted[id] = true
	}
	// Make unselected list out of all nodes form matrix without selected nodes
	m := len(distanceMatrix)
	unselected := make([]int, 0, m)
	for i := 0; i < m; i++ {
		if _, ok := visted[i]; !ok {
			unselected = append(unselected, i)
		}
	}

	// Run local search function as long as there is improvement
	solution, improved := SteepestMove(selectedIDs, visted, unselected, distanceMatrix, "EdgeExchange")
	for improved {
		solution, improved = SteepestMove(solution, visted, unselected, distanceMatrix, "EdgeExchange")

	}
	return solution
}