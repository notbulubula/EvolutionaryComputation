package local_search

import (
	"evolutionary_computation/methods"
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

	moves := getMovesDelta(distanceMatrix, solution, unselected)

	improved := true

	for improved {
		solution, improved = SteepestDelta(solution, visted, unselected, distanceMatrix, moves)
		continue
	}
	return solution
}

type MoveDelta struct {
	moveType string
	i, j     int // indices of nodes involved
	delta    int // change in cost
}

func getMovesDelta(distanceMatrix [][]int, solution []int, unselectedNodes []int) []MoveDelta {
	var moves []MoveDelta
	n := len(solution)

	// Intra-route: two-edges exchange
	for i := 0; i < n; i++ {
		for j := i + 2; j < n; j++ {
			delta := deltaTwoEdgesExchange(solution, i, j, distanceMatrix)
			if delta < 0 {
				moves = append(moves, MoveDelta{"twoEdgesExchange", i, j, delta})
			}
		}
	}

	// Inter-route: exchange between selected and unselected nodes
	for i := 0; i < n; i++ {
		for _, unselected := range unselectedNodes {
			delta := deltaInterRouteExchange(solution, i, unselected, distanceMatrix)
			if delta < 0 {
				moves = append(moves, MoveDelta{"interRouteExchange", i, unselected, delta})
			}
		}
	}

	// Sort based on delta in ascending order
	sort.Slice(moves, func(a, b int) bool {
		return moves[a].delta < moves[b].delta
	})

	return moves
}
