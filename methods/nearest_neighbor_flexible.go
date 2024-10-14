package methods

// NearestNeighborFlexible generates a solution using the nearest neighbor heuristic, starting from a random node.
// The algorithm selects the nearest neighbor of any node in the solution until half of the nodes are selected.
func NearestNeighborFlexible(distanceMatrix [][]int, startNode int) []int {
	numNodes := len(distanceMatrix)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	selectedIDs := make([]int, 0, numToSelect)

	selectedIDs = append(selectedIDs, startNode)

	// Keep track of visited nodes
	visited := make(map[int]bool)
	visited[startNode] = true

	// Continue adding the nearest neighbor until half of the nodes are selected
	for len(selectedIDs) < numToSelect {
		// Find the nearest neighbor that has not been visited
		nearestNeighbor := -1
		minInsertionCost := int(^uint(0) >> 1) // Max int value
		insertPosition := -1

		for i := 0; i < numNodes; i++ {
			if !visited[i] {
				// Calculate insertion costs for each position in the selectedIDs
				for j := 0; j <= len(selectedIDs); j++ {
					var cost int
					if j == 0 {
						// Insert at the beginning
						cost = distanceMatrix[selectedIDs[0]][i]
					} else if j == len(selectedIDs) {
						// Insert at the end
						cost = distanceMatrix[selectedIDs[len(selectedIDs)-1]][i]
					} else {
						// Insert between selectedIDs[j-1] and selectedIDs[j]
						cost = distanceMatrix[selectedIDs[j-1]][i] +
							distanceMatrix[i][selectedIDs[j]] -
							distanceMatrix[selectedIDs[j-1]][selectedIDs[j]]
					}

					// Update nearest neighbor if a lower cost is found
					if cost < minInsertionCost {
						minInsertionCost = cost
						nearestNeighbor = i
						insertPosition = j // Store the position for insertion
					}
				}
			}
		}

		// Insert the nearest neighbor at the best position found
		if insertPosition == len(selectedIDs) {
			// If the insert position is at the end, simply append the element
			selectedIDs = append(selectedIDs, nearestNeighbor)
		} else {
			// Otherwise, insert the element at the correct position
			selectedIDs = append(selectedIDs[:insertPosition+1], selectedIDs[insertPosition:]...) // Resize the slice
			selectedIDs[insertPosition] = nearestNeighbor
		}

		// Mark the nearest neighbor as visited
		visited[nearestNeighbor] = true
	}

	// Return list of visited node IDs in the order of the cycle
	return selectedIDs
}
