package methods

import "evolutionary_computation/utils"

func GreedyTwoRegret(distanceMatrix [][]int, startNode int) []int {
	numNodes, numToSelect, solution, visited := utils.GetInitialState(distanceMatrix, startNode)

	return solution
}