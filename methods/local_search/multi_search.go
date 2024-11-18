package local_search

import (
    "time"
    "evolutionary_computation/utils"
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

    startTime := time.Now()

    //TODO: change the time to average from MultiLocalSearch
    for time.Since(startTime) < 30*time.Second {
        startNode := callCount % len(costMatrix)

        if callCount == 0 {
            solution = RandomSteepestIntraEdge(costMatrix, startNode)
        } else {
            bestSolutionCopy := make([]int, len(bestSolution))
            copy(bestSolutionCopy, bestSolution)
            //TODO add permutation of bestSolutionCopy


            solution = SteepestIntraEdgeFromSolution(bestSolutionCopy, costMatrix, startNode)
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
