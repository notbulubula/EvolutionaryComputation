package local_search

import (
	"evolutionary_computation/methods"
	"sort"
)

func LS_Candidates(distanceMatrix [][]int, startNode int) []int {
	// Run random function to get the initial solution
	solution := methods.RandomSolution(distanceMatrix, startNode)
	visted := make(map[int]bool)
	// Make visited map
	for _, id := range solution {
		visted[id] = true
	}

	// Make unselected list out of all nodes from matrix without selected nodes
	m := len(distanceMatrix)
	unselected := make([]int, 0, m)
	for i := 0; i < m; i++ {
		if _, ok := visted[i]; !ok {
			unselected = append(unselected, i)
		}
	}

	moves := getCandidateMoves(distanceMatrix, 10)
	// Run local search function as long as there is improvement
	improved := true
	// for k:=0; k<5 && improved; k++ {
	for improved {
		solution, improved = SteepestCandidate(solution, visted, unselected, distanceMatrix, moves)

	}
	return solution
}

func getCandidateMoves(distanceMatrix [][]int, N int) []Move {
	var moves []Move

	for i, row := range distanceMatrix {
		// Collect indices and values of each row, ignoring the diagonal (i == j)
		type IndexedValue struct {
			j, value int
		}
		var indexedValues []IndexedValue
		for j, value := range row {
			if i != j { // Ignore diagonal
				indexedValues = append(indexedValues, IndexedValue{j, value})
			}
		}

		// Sort based on values in ascending order
		sort.Slice(indexedValues, func(a, b int) bool {
			return indexedValues[a].value < indexedValues[b].value
		})

		// Select the N smallest values and create Moves
		for k := 0; k < N && k < len(indexedValues); k++ {
			moves = append(moves, Move{"", i, indexedValues[k].j})
		}
	}

	return moves
}
