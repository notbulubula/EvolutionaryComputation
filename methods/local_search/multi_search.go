package local_search

import "evolutionary_computation/utils"

func MultiLocalSearch(costMatrix [][]int, pointless_value int) []int {
	var bestFitness int
	var bestSolution []int

	for i := 0; i < 200; i++ {
		startNode := i % len(costMatrix)

		solution := RandomSteepestIntraEdge(costMatrix, startNode)
		fitness := utils.Fitness(solution, costMatrix)

		if i == 0 || fitness < bestFitness {
			bestFitness = fitness
			bestSolution = solution
		}
	}

	return bestSolution
}
