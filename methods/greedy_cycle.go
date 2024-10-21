package methods

import "evolutionary_computation/utils"

// GreedyCycle function starts by selecting a random vertex as the starting point.
// It builds a cycle by repeatedly inserting the nearest vertex that minimizes the cycle length increase.
// The process continues until all vertices are added to form a complete cycle.
func GreedyCycle(distanceMatrix [][]int, startNode int) []int {

	numNodes, numToSelect, selectedIDs, visited := utils.GetInitialState(distanceMatrix, startNode)

	// Continue adding the vertices until all are selected
	for len(selectedIDs) < numToSelect {
		bestNode := -1
		bestPosition := -1
		bestIncrease := int(^uint(0) >> 1) // Max int value

		// Find the best unvisited node to insert and the best position in the cycle
		for i := 0; i < numNodes; i++ {
			if !visited[i] {
				for j := 0; j < len(selectedIDs); j++ {
					// Calculate the increase in cycle length by inserting node i between selectedIDs[j] and selectedIDs[(j+1) % len(selectedIDs)]
					next := (j + 1) % len(selectedIDs)
					increase := distanceMatrix[selectedIDs[j]][i] +
						distanceMatrix[i][selectedIDs[next]] -
						distanceMatrix[selectedIDs[j]][selectedIDs[next]]

					// Find the minimum increase
					if increase < bestIncrease {
						bestIncrease = increase
						bestNode = i
						bestPosition = next
					}
				}
			}
		}

		// Insert the bestNode in the bestPosition found
		selectedIDs = append(selectedIDs[:bestPosition], append([]int{bestNode}, selectedIDs[bestPosition:]...)...)
		visited[bestNode] = true
	}

	return selectedIDs
}
