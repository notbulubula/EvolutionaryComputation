package local_search

import (
	"evolutionary_computation/methods"
	"evolutionary_computation/utils"
	"math/rand"
	"time"
)


func LargeNeighbourhoodWithLS(costMatrix [][]int, pointless_value int) []int {
	var bestFitness int
	var bestSolution []int
	var callCount int
	var solution []int

	percentage := 0.2

	startTime := time.Now()

	//TODO: change the time to average from MultiLocalSearch
	for time.Since(startTime) < 24*time.Second {
		startNode := callCount % len(costMatrix)

		if callCount == 0 {
			solution = RandomSteepestIntraEdge(costMatrix, startNode)
		} else {
			solution = bestSolution // Always use the best solution to perform operations
		}
		callCount++

		destroyedSolution := DestroySolution(solution, percentage)
		repairedSolution := methods.NearestNeighborFlexibleFromSolution(costMatrix, destroyedSolution)
		solution = SteepestIntraEdgeFromSolution(repairedSolution, costMatrix, startNode)

		fitness := utils.Fitness(solution, costMatrix)
		if callCount == 1 || fitness < bestFitness {
			bestFitness = fitness
			bestSolution = solution
		}
	}
	// TODO: add callCount to results dict
	println("Number of calls:", callCount)
	return bestSolution
}


func LargeNeighbourhood(costMatrix [][]int, pointless_value int) []int {
	var bestFitness int
	var bestSolution []int
	var callCount int
	var solution []int

	percentage := 0.2

	startTime := time.Now()

	//TODO: change the time to average from MultiLocalSearch
	for time.Since(startTime) < 24*time.Second {
		startNode := callCount % len(costMatrix)

		if callCount == 0 {
			solution = RandomSteepestIntraEdge(costMatrix, startNode)
		} else {
			solution = bestSolution // Always use the best solution to perform operations
		}
		callCount++

		destroyedSolution := DestroySolution(solution, percentage)
		//NOTE: Culprit number 2 if things break
		solution := methods.NearestNeighborFlexibleFromSolution(costMatrix, destroyedSolution)

		fitness := utils.Fitness(solution, costMatrix)
		if callCount == 1 || fitness < bestFitness {
			bestFitness = fitness
			bestSolution = solution
		}
	}
	// TODO: add callCount to results dict
	println("Number of calls:", callCount)
	return bestSolution
}

func DestroySolution(solution []int, percentage float64) []int {
    numNodesToRemove := int(float64(len(solution)) * percentage)
    if numNodesToRemove == 0 {
        return solution // Nothing to remove
    }

    // Split into 3 groups (as evenly as possible)
    groupSizes := []int{
        numNodesToRemove / 3,
        numNodesToRemove / 3,
        numNodesToRemove - 2*(numNodesToRemove/3), // Remaining nodes go to the last group
    }

    modifiedSolution := make([]int, len(solution))
    copy(modifiedSolution, solution)

    for _, groupSize := range groupSizes {
        if groupSize > 0 {
            // Randomly select a starting index and create a subpath of groupSize
            start := rand.Intn(len(modifiedSolution) - groupSize + 1)
			// NOTE: If things break, this is the most likely culprit
            modifiedSolution = append(modifiedSolution[:start], modifiedSolution[start+groupSize:]...)
        }
    }

    return modifiedSolution
}
