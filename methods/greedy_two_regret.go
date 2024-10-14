package methods

import (
	"evolutionary_computation/utils"
	"math"
	"sort"
	"fmt"
)

func GreedyTwoRegret(distanceMatrix [][]int, startNode int) []int {
	_, numToSelect, solution, visited := utils.GetInitialState(distanceMatrix, startNode)

	for len(solution) < numToSelect {
		best1, best2 := twoBestCandidates(visited, solution, distanceMatrix)

		bestCost1, secondBest1, insertPos1 := getBestInsertionCost(best1, solution, distanceMatrix)
		regret1 := bestCost1 - secondBest1
		bestCost2, secondBest2, insertPos2 := getBestInsertionCost(best2, solution, distanceMatrix)
		regret2 := bestCost2 - secondBest2

		if regret1 >= regret2 {
			solution = utils.InsertAt(solution, insertPos1, best1)
			visited[best1] = true
		} else {
			solution = utils.InsertAt(solution, insertPos2, best2)
			visited[best2] = true
		}
	}

	return solution
}

func GreedyTwoRegretWithWeights(distanceMatrix [][]int, startNode int) []int {
	_, numToSelect, solution, visited := utils.GetInitialState(distanceMatrix, startNode)


	for len(solution) < numToSelect {
		currentFitness := utils.Fitness(solution, distanceMatrix)

		best1, best2 := twoBestCandidates(visited, solution, distanceMatrix)
		
		bestCost1, secondBest1, insertPos1 := getBestInsertionCost(best1, solution, distanceMatrix)
		regret1 := bestCost1 - secondBest1
		best1Solution := utils.InsertAt(solution, insertPos1, best1)
		best1Fitness := utils.Fitness(best1Solution, distanceMatrix)
		
		bestCost2, secondBest2, insertPos2 := getBestInsertionCost(best2, solution, distanceMatrix)
		regret2 := bestCost2 - secondBest2
		best2Solution := utils.InsertAt(solution, insertPos2, best2)
		best2Fitness := utils.Fitness(best2Solution, distanceMatrix)

		weight1 := calculateWeight(regret1, best1Fitness, currentFitness)
		weight2 := calculateWeight(regret2, best2Fitness, currentFitness)

		if weight1 >= weight2 {
			solution = utils.InsertAt(solution, insertPos1, best1)
			visited[best1] = true
			fmt.Printf("Appended %v, visited? %v \n", best1, visited[best1])
		} else {
			solution = utils.InsertAt(solution, insertPos2, best2)
			visited[best2] = true
			fmt.Printf("Appended %v, visited? %v \n", best2, visited[best2])
		}

	}

	return solution
}

func calculateWeight(regret int, newFitness int, currentFitness int) int {
	regretWeight := 1
	changeWeight := -1 //We normally minimize the change, but want to maximize the equation
	return regretWeight * regret + changeWeight * (newFitness - currentFitness)
}

func twoBestCandidates(visited map[int]bool, solution []int, distanceMatrix [][]int) (int, int) {
	type candidate struct {
		node   int
		cost   int
		insert int
	}

	var candidates []candidate

	// Evaluate all unvisited nodes
	for i, seen := range visited {
		if !seen {
			cost, _, insertPos := getBestInsertionCost(i, solution, distanceMatrix)
			candidates = append(candidates, candidate{node: i, cost: cost, insert: insertPos})
		}
	}

	// Sort candidates by their insertion cost (ascending order)
	sort.Slice(candidates, func(a, b int) bool {
		return candidates[a].cost < candidates[b].cost
	})
	return candidates[0].node, candidates[1].node
}

// Returns the best insertion cost and the position for the given node in the solution
func getBestInsertionCost(node int, solution []int, distanceMatrix [][]int) (int, int, int) {
	bestCost := math.MaxInt32
	secondBestCost := math.MaxInt32
	bestPos := 0

	for j := 0; j <= len(solution); j++ {
		var cost int

		if j == 0 {
			// Insert at the beginning
			cost = distanceMatrix[node][solution[0]]
		} else if j == len(solution) {
			// Insert at the end
			cost = distanceMatrix[solution[len(solution)-1]][node]
		} else {
			// Insert between solution[j-1] and solution[j]
			cost = distanceMatrix[solution[j-1]][node] +
				distanceMatrix[node][solution[j]] -
				distanceMatrix[solution[j-1]][solution[j]]
		}

		// Update best and second-best costs and track the best position
		if cost < bestCost {
			secondBestCost = bestCost
			bestCost = cost
			bestPos = j
		} else if cost < secondBestCost {
			secondBestCost = cost
		}
	}

	return bestCost, secondBestCost, bestPos
}

