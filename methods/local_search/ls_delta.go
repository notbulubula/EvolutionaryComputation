package local_search

import (
	"evolutionary_computation/methods"
	"math"
	"sort"
)

func LS_Delta(distanceMatrix [][]int, startNode int) []int {
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

	moves := getMovesDelta(distanceMatrix, solution)

	improved := true
	for improved {
		solution, improved = SteepestDelta(solution, visted, unselected, distanceMatrix, moves)
		continue
	}
	return solution
}

type MoveDelta struct {
	i, j  int // indices of nodes involved
	delta int // change in cost
}

func getMovesDelta(distanceMatrix [][]int, solution []int) []MoveDelta {
	var moves []MoveDelta
	n := len(solution)

	// Generate all combinations of moves and add None as delta
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			moves = append(moves, MoveDelta{i, j, 0})
		}
	}

	for move := range moves {
		move_i_index := findIndex(solution, moves[move].i)
		move_j_index := findIndex(solution, moves[move].j)

		//if i and j are in solution and not neighbours
		// calculate TwoEdgesExchange delata
		if move_i_index != -1 && move_j_index != -1 && math.Abs(float64(move_i_index-move_j_index)) != 1 {
			delta := deltaTwoEdgesExchange(solution, move_i_index, move_j_index, distanceMatrix)
			if delta < 0 {
				moves[move].delta = delta
			}

		}

		//if i is in solution and j is not
		// calculate InterRouteExchange delta
		if move_i_index != -1 && move_j_index == -1 {
			delta := deltaInterRouteExchange(solution, move_i_index, moves[move].j, distanceMatrix)
			if delta < 0 {
				moves[move].delta = delta
			}

		}
	}

	// Sort based on delta in ascending order
	sort.Slice(moves, func(a, b int) bool {
		return moves[a].delta < moves[b].delta
	})

	return moves
}
