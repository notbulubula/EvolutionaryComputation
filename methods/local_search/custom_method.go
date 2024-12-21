package local_search

import (
	"evolutionary_computation/methods"
	"evolutionary_computation/utils"
	"math"
	"math/rand"
	"time"
)

// A custom Large Neighbourhood Search method that uses more randomization in destroying
// it utilizes tabu search to avoid revisiting the same solutions and to explore more of the solution space
// BestSolution is approved with the use of simulated annealing to improve exploration
func CustomMethod(costMatrix [][]int, startNode int) []int {
	var bestFitness int
	var bestSolution []int
	var currentFitness int
	var currentSolution []int

	// Parameters
	percentage := 0.43
	temperature := 1500.0
	coolingRate := 0.995
	tabuTenure := 20
	tabuList := make(map[string]int) // Tabu list as a map of solution hashes to iteration count
	callCount := 0

	currentSolution = methods.RandomSolution(costMatrix, startNode)
	currentFitness = utils.Fitness(currentSolution, costMatrix)
	bestFitness = currentFitness
	bestSolution = currentSolution

	startTime := time.Now()

	for time.Since(startTime) < 3*time.Second {
		// Destroy and repair solution
		destroyedSolution := DestroySolutionRandom(currentSolution, percentage)
		repairedSolution := methods.NearestNeighborFlexibleFromSolution(costMatrix, destroyedSolution)

		newSolution := SteepestIntraEdgeFromSolution(repairedSolution, costMatrix, startNode)
		newFitness := utils.Fitness(newSolution, costMatrix)

		// Convert solution to string (or hash) for tabu list
		solutionKey := utils.SolutionToString(newSolution)

		// Check tabu list
		isTabu := tabuList[solutionKey] > callCount
		aspiration := newFitness < bestFitness // Override tabu if fitness is better than the best

		if (!isTabu || aspiration) && (newFitness < currentFitness || rand.Float64() < math.Exp(-float64(newFitness-currentFitness)/temperature)) {
			currentSolution = newSolution
			currentFitness = newFitness

			// Update the tabu list
			tabuList[solutionKey] = callCount + tabuTenure

			// Update the best solution
			if currentFitness < bestFitness {
				bestFitness = currentFitness
				bestSolution = currentSolution
			}
		}

		// Cool down the temperature
		temperature *= coolingRate
		callCount++

		// Periodic cleanup of the tabu list
		if callCount%50 == 0 {
			for k, v := range tabuList {
				if v <= callCount {
					delete(tabuList, k)
				}
			}
		}
	}

	println("Number of calls:", callCount)
	return bestSolution
}

func DestroySolutionRandom(solution []int, percentage float64) []int {
    // Calculate the total number of nodes to remove
    numNodesToRemove := int(float64(len(solution)) * percentage)
    if numNodesToRemove <= 0 || len(solution) <= 1 {
        return solution // Nothing to remove
    }

    modifiedSolution := make([]int, len(solution))
    copy(modifiedSolution, solution)

    // Randomly decide the number of groups (2 to 12)
    numGroups := rand.Intn(3) + 2 // Generates a random number in [2, 5]

    // Distribute the nodes to remove across groups
    groupSizes := make([]int, numGroups)
    for i := 0; i < numNodesToRemove; i++ {
        groupSizes[rand.Intn(numGroups)]++ // Increment a random group's size
    }

    // Randomly remove nodes for each group
    for _, groupSize := range groupSizes {
        if groupSize > 0 {
            start := rand.Intn(len(modifiedSolution) - groupSize + 1) // Select random start index
            modifiedSolution = append(modifiedSolution[:start], modifiedSolution[start+groupSize:]...)
        }
    }

    return modifiedSolution
}
