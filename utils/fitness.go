package utils

// Fitness calculates the total cost (sum of distances and node costs) of a given solution.
func Fitness(solution []int, nodes []Node, distanceMatrix [][]int) int {
	totalDistance := 0
	totalCost := 0

	numNodes := len(solution)

	for i := 0; i < numNodes; i++ {
		currentNode := solution[i]
		nextNode := solution[(i+1)%numNodes] // Ensure wrap-around for the cycle

		totalDistance += distanceMatrix[currentNode][nextNode]

		totalCost += nodes[currentNode].Cost
	}

	// The fitness is the sum of total distance and total cost
	return totalDistance + totalCost
}
