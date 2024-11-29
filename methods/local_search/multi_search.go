package local_search

import (
	"evolutionary_computation/utils"
	"math/rand"
	"time"
)

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

func IterativeLocalSearch(costMatrix [][]int, pointless_value int) []int {
	var bestFitness int
	var bestSolution []int
	var callCount int
	var solution []int

	var percentage float64 = 0.3

	startTime := time.Now()

	//TODO: change the time to average from MultiLocalSearch
	for time.Since(startTime) < 30*time.Second {
		startNode := callCount % len(costMatrix)

		if callCount == 0 {
			solution = RandomSteepestIntraEdge(costMatrix, startNode)
		} else {
			bestSolutionCopy := make([]int, len(bestSolution))
			copy(bestSolutionCopy, bestSolution)

			permutatedSolution := PermuteSolution(bestSolutionCopy, percentage)

			solution = SteepestIntraEdgeFromSolution(permutatedSolution, costMatrix, startNode)
		}

		callCount++
		fitness := utils.Fitness(solution, costMatrix)

		if callCount == 1 || fitness < bestFitness {
			bestFitness = fitness
			bestSolution = solution
		}
	}
	// TODO: add callCount to results dict
	println("Number of RandomSteepestIntraEdge calls:", callCount)
	return bestSolution
}

func PermuteSolution(solution []int, percentage float64) []int {
	// take 2 random indexes wchich will cover n percent of the solution
	m := len(solution)
	start := rand.Intn(m)
	end := start + int(percentage*float64(m))

	// permute the solution
	for i := start; i < end; i++ {
		//change the node to random value outside of the solution

		new_node := rand.Intn(m)
		for utils.Contains(solution, new_node) {
			new_node = rand.Intn(m)
		}
		index := i % m
		solution[index] = new_node
	}

	return solution
}
