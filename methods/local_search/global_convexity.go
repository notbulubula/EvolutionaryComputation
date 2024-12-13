package local_search

import (
	"evolutionary_computation/utils"
)

func GlobalConvexityLS(costMatrix [][]int, bestsolution []int, similarity_measure string, similarity_to string) ([]float64, []int) {
	iterations := 1000
	localOptima := make([][]int, 0, iterations)
	fitnesses := make([]int, 0, iterations)
	similarities := make([]float64, 0, iterations)

	// Generate 1000 random solutions and optimize them using greedy local search
	for i := 0; i < iterations; i++ {
		startNode := i % len(costMatrix)
		solution := RandomGreedyIntraEdge(costMatrix, startNode)
		fitness := utils.Fitness(solution, costMatrix)

		localOptima = append(localOptima, solution)
		fitnesses = append(fitnesses, fitness)
	}

	// Compute similarity
	for i, solution := range localOptima {
		var sim float64

		if similarity_to == "best" {
			// Compare similarity to the best solution
			if similarity_measure == "common_nodes" {
				sim = utils.CommonNodes(solution, bestsolution)
			} else {
				sim = utils.CommonEdges(solution, bestsolution)
			}
		} else {
			// Compare similarity to the average of all other local optima
			numOptima := len(localOptima) - 1 // Exclude current solution from comparison

			var totalSimilarity float64
			for j, otherSolution := range localOptima {
				if i == j {
					continue // Skip the current solution
				}

				if similarity_measure == "common_nodes" {
					totalSimilarity += utils.CommonNodes(solution, otherSolution)
				} else {
					totalSimilarity += utils.CommonEdges(solution, otherSolution)
				}

			}
			sim = totalSimilarity / float64(numOptima)

		}

		similarities = append(similarities, sim)
	}

	return similarities, fitnesses
}
