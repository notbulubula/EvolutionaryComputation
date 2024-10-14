package methods

import (
	"evolutionary_computation/utils"
	"math"
	"sort"
)

func GreedyTwoRegret(distanceMatrix [][]int, startNode int) []int {
	_, numToSelect, solution, visited := utils.GetInitialState(distanceMatrix, startNode)

	for len(solution) < numToSelect {
		best1, best2 := twoBestCandidates(visited, solution, distanceMatrix)

		regret1, insertPos1 := getBestInsertionCost(best1, solution, distanceMatrix)
		regret2, insertPos2 := getBestInsertionCost(best2, solution, distanceMatrix)

		if regret1 >= regret2 {
			solution = append(solution[:insertPos1], append([]int{best1}, solution[insertPos1:]...)...)
			visited[best1] = true
		} else {
			solution = append(solution[:insertPos2], append([]int{best2}, solution[insertPos2:]...)...)
			visited[best2] = true
		}
	}

	return solution
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
			cost, insertPos := getBestInsertionCost(i, solution, distanceMatrix)
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
func getBestInsertionCost(node int, solution []int, distanceMatrix [][]int) (int, int) {
	bestCost := math.MaxInt32
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

		// Track the best insertion position and cost
		if cost < bestCost {
			bestCost = cost
			bestPos = j
		}
	}

	return bestCost, bestPos
}
