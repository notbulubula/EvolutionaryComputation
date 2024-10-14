package methods

import (
	"math/rand"
)

// RandomSolution generates a random solution and returns a list of node IDs in the selected order.
func RandomSolution(distanceMatrix [][]int) []int {
	numNodes := len(distanceMatrix)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	selectedIDs := make([]int, 0, numToSelect)

	perm := rand.Perm(numNodes)
	for i := 0; i < numToSelect; i++ {
		selectedIDs = append(selectedIDs, perm[i])
	}

	// Construct a random cycle with the selected nodes' IDs
	return selectedIDs // Return list of visited node IDs in the order of the cycle
}
