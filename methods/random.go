package methods

import (
	"evolutionary_computation/utils"
	"math/rand"
)

// RandomSolution generates a random solution and returns a list of node IDs in the selected order.
func RandomSolution(nodes []utils.Node, distanceMatrix [][]int) []int {
	numNodes := len(nodes)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	selectedIDs := make([]int, 0, numToSelect)

	perm := rand.Perm(numNodes)
	for i := 0; i < numToSelect; i++ {
		selectedIDs = append(selectedIDs, nodes[perm[i]].ID)
	}

	// Construct a random cycle with the selected nodes' IDs
	return selectedIDs // Return list of visited node IDs in the order of the cycle
}
