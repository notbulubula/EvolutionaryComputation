package utils

import (
	"fmt"
	"strings"
)
// Returns number of total nodes, number of nodes to select
// solution list with initial Node and a map of visited nodes
// Use to remove duplicated code
func GetInitialState(costMatrix [][]int, startNode int) (int, int, []int, map[int]bool) {
	numNodes := len(costMatrix)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	solution := make([]int, 0, numToSelect)

	solution = append(solution, startNode)
	visited := make(map[int]bool)

	for i := 0; i < numNodes; i++ {
		visited[i] = false
	}
	visited[startNode] = true

	return numNodes, numToSelect, solution, visited
}

func GetSuggestedState(costMatrix [][]int, solution []int) (int, int, map[int]bool) {
	numNodes := len(costMatrix)
	numToSelect := (numNodes + 1) / 2 // Rounds up if odd
	visited := make(map[int]bool)

	for i := 0; i < numNodes; i++ {
		visited[i] = false
	}

	for _, node := range solution {
		visited[node] = true
	}

	return numNodes, numToSelect, visited
}

func InsertAt(slice []int, index int, element int) []int {
	if index < 0 || index > len(slice) {
		// Handle index out of bounds
		return slice
	}
	return append(slice[:index], append([]int{element}, slice[index:]...)...)
}

func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func SolutionToString(solution []int) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(solution)), ","), "[]")
}
