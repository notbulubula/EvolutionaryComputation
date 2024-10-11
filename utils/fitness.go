package utils

// Fitness calculates the total cost (sum of distances and node costs) of a given solution.
func Fitness(solution []int, distanceMatrix [][]int) int {
	totalCost := 0

	numNodes := len(solution)

	for i := 0; i < numNodes; i++ {
		currentNode := solution[i]
		nextNode := solution[(i+1)%numNodes] // Ensure wrap-around for the cycle

		totalCost += distanceMatrix[currentNode][nextNode]

	}

	// The fitness is the sum of total distance and total cost
	return totalCost
}
