package methods

import (
	"evolutionary_computation/utils"
	"math/rand"
)

// NearestNeighborFlexible generates a solution using the nearest neighbor heuristic, starting from a random node.
// The algorithm selects the nearest neighbor of any node in the solution until half of the nodes are selected.
func NearestNeighborFlexible(nodes []utils.Node, distanceMatrix [][]int) []int {
	numNodes := len(nodes)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	selectedIDs := make([]int, 0, numToSelect)

	// Select a random starting node and add it to the solution
	startNode := rand.Intn(len(nodes))
	selectedIDs = append(selectedIDs, nodes[startNode].ID)

	// Keep track of visited nodes
	visited := make(map[int]bool)
	visited[startNode] = true

	// Continue adding the nearest neighbor until half of the nodes are selected
	for len(selectedIDs) < numToSelect {
		// Find the nearest neighbor that has not been visited
		nearestNeighbor := -1
		minInsertionCost := -1
		insertPosition := -1

		for i := 0; i < numNodes; i++ {
			if !visited[i] {
				// Calculate insertion costs for each position in the selectedIDs
				for j := 0; j <= len(selectedIDs); j++ {
					var cost int
					if j == 0 {
						// Insert at the beginning
						cost = distanceMatrix[i][selectedIDs[0]]
					} else if j == len(selectedIDs) {
						// Insert at the end
						cost = distanceMatrix[selectedIDs[len(selectedIDs)-1]][i]
					} else {
						// Insert between selectedIDs[j-1] and selectedIDs[j]
						cost = distanceMatrix[selectedIDs[j-1]][i] + distanceMatrix[i][selectedIDs[j]]
					}

					// Update nearest neighbor if a lower cost is found
					if minInsertionCost == -1 || cost < minInsertionCost {
						minInsertionCost = cost
						nearestNeighbor = i
						insertPosition = j // Store the position for insertion
					}
				}
			}
		}

		// Insert the nearest neighbor at the best position found
		selectedIDs = insertAtPosition(selectedIDs, nodes[nearestNeighbor].ID, insertPosition)

		// Mark the nearest neighbor as visited
		visited[nearestNeighbor] = true
	}

	// Return list of visited node IDs in the order of the cycle
	return selectedIDs
}

// insertAtPosition inserts an element at a specified position in a slice
func insertAtPosition(slice []int, value int, position int) []int {
	// If the position is at the end, just append the value
	if position == len(slice) {
		return append(slice, value)
	}
	// If the position is at the start or middle, create a new slice with the value inserted
	slice = append(slice[:position+1], slice[position:]...) // Resize the slice to accommodate the new value
	slice[position] = value
	return slice
}
