package methods

import "evolutionary_computation/utils"

// NearestNeighborFlexible generates a solution using the nearest neighbor heuristic, starting from a random node.
// The algorithm selects the nearest neighbor of any node in the solution until half of the nodes are selected.
func NearestNeighborFlexible(distanceMatrix [][]int, startNode int) []int {
	numNodes, numToSelect, solution, visited := utils.GetInitialState(distanceMatrix, startNode)

	// Continue adding the nearest neighbor until half of the nodes are selected
	for len(solution) < numToSelect {
		// Find the nearest neighbor that has not been visited
		nearestNeighbor := -1
		minInsertionCost := int(^uint(0) >> 1) // Max int value
		insertPosition := -1

		for i := 0; i < numNodes; i++ {
			if !visited[i] {
				// Calculate insertion costs for each position in the solution
				for j := 0; j <= len(solution); j++ {
					var cost int
					if j == 0 {
						// Insert at the beginning
						cost = distanceMatrix[solution[0]][i]
					} else if j == len(solution) {
						// Insert at the end
						cost = distanceMatrix[solution[len(solution)-1]][i]
					} else {
						// Insert between solution[j-1] and solution[j]
						cost = distanceMatrix[solution[j-1]][i] +
							distanceMatrix[i][solution[j]] -
							distanceMatrix[solution[j-1]][solution[j]]
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
		if insertPosition == len(solution) {
			// If the insert position is at the end, simply append the element
			solution = append(solution, nearestNeighbor)
		} else {
			// Otherwise, insert the element at the correct position
			solution = append(solution[:insertPosition+1], solution[insertPosition:]...) // Resize the slice
			solution[insertPosition] = nearestNeighbor
		}

		// Mark the nearest neighbor as visited
		visited[nearestNeighbor] = true
	}

	// Return list of visited node IDs in the order of the cycle
	return solution
}

func NearestNeighborFlexibleFromSolution(distanceMatrix [][]int, solution []int) []int {
	numNodes, numToSelect, visited := utils.GetSuggestedState(distanceMatrix, solution)

	// Continue adding the nearest neighbor until half of the nodes are selected
	for len(solution) < numToSelect {
		// Find the nearest neighbor that has not been visited
		nearestNeighbor := -1
		minInsertionCost := int(^uint(0) >> 1) // Max int value
		insertPosition := -1

		for i := 0; i < numNodes; i++ {
			if !visited[i] {
				// Calculate insertion costs for each position in the solution
				for j := 0; j <= len(solution); j++ {
					var cost int
					if j == 0 {
						// Insert at the beginning
						cost = distanceMatrix[solution[0]][i]
					} else if j == len(solution) {
						// Insert at the end
						cost = distanceMatrix[solution[len(solution)-1]][i]
					} else {
						// Insert between solution[j-1] and solution[j]
						cost = distanceMatrix[solution[j-1]][i] +
							distanceMatrix[i][solution[j]] -
							distanceMatrix[solution[j-1]][solution[j]]
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
		if insertPosition == len(solution) {
			// If the insert position is at the end, simply append the element
			solution = append(solution, nearestNeighbor)
		} else {
			// Otherwise, insert the element at the correct position
			solution = append(solution[:insertPosition+1], solution[insertPosition:]...) // Resize the slice
			solution[insertPosition] = nearestNeighbor
		}

		// Mark the nearest neighbor as visited
		visited[nearestNeighbor] = true
	}

	// Return list of visited node IDs in the order of the cycle
	return solution
}
