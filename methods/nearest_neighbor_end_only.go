package methods

import (
	"math/rand"
)

// NearestNeighborEndOnly generates a solution using the nearest neighbor heuristic, starting from a random node.
// The algorithm selects the nearest neighbor of the last node in the solution until half of the nodes are selected.
func NearestNeighborEndOnly(distanceMatrix [][]int) []int {
	numNodes := len(distanceMatrix)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	selectedIDs := make([]int, 0, numToSelect)

	// Select a random starting node and add it to the solution
	startNode := rand.Intn(numNodes)
	selectedIDs = append(selectedIDs, startNode)

	// Keep track of visited nodes
	visited := make(map[int]bool)
	visited[startNode] = true

	// Continue adding the nearest neighbor of last node until half of the nodes are selected
	for len(selectedIDs) < numToSelect {
		// Get the last node in the solution
		lastNode := selectedIDs[len(selectedIDs)-1]

		// Find the nearest neighbor that has not been visited
		nearestNeighbor := -1
		minDistance := -1

		for i := 0; i < numNodes; i++ {
			if !visited[i] {
				distance := distanceMatrix[lastNode][i]
				if minDistance == -1 || distance < minDistance {
					minDistance = distance
					nearestNeighbor = i
				}
			}
		}

		// Add the nearest neighbor to the solution
		selectedIDs = append(selectedIDs, nearestNeighbor)
		visited[nearestNeighbor] = true
	}

	// Construct a cycle with the selected nodes' IDs
	return selectedIDs // Return list of visited node IDs in the order of the cycle
}
