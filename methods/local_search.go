package methods

import (
	"math/rand"
)

type Move struct {
	moveType string
	i, j     int // indices of nodes involved
}


func generateMoves(solution []int, visited map[int]bool, unselectedNodes []int) []Move {
	var moves []Move
	n := len(solution)

	// Intra-route: two-nodes exchange
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			moves = append(moves, Move{"twoNodesExchange", i, j})
		}
	}

	// Intra-route: two-edges exchange
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			moves = append(moves, Move{"twoEdgesExchange", i, j})
		}
	}

	// Inter-route: exchange between selected and unselected nodes
	for i := 0; i < n; i++ {
		for _, unselected := range unselectedNodes {
			moves = append(moves, Move{"interRouteExchange", i, unselected})
		}
	}

	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })

	return moves
}


func deltaTwoNodesExchange(solution []int, i, j int, distanceMatrix [][]int) int {
	n := len(solution)

	prevI := solution[(i-1+n)%n]
	nextI := solution[(i+1)%n]
	prevJ := solution[(j-1+n)%n]
	nextJ := solution[(j+1)%n]

	costBefore := distanceMatrix[prevI][solution[i]] + distanceMatrix[solution[i]][nextI] +
		distanceMatrix[prevJ][solution[j]] + distanceMatrix[solution[j]][nextJ]

	costAfter := distanceMatrix[prevI][solution[j]] + distanceMatrix[solution[j]][nextI] +
		distanceMatrix[prevJ][solution[i]] + distanceMatrix[solution[i]][nextJ]

	return costAfter - costBefore
}


func deltaTwoEdgesExchange(solution []int, i, j int, distanceMatrix [][]int) int {
	n := len(solution)

	nodeI := solution[i]
	nextI := solution[(i+1)%n]
	nodeJ := solution[j]
	nextJ := solution[(j+1)%n]

	costBefore := distanceMatrix[nodeI][nextI] + distanceMatrix[nodeJ][nextJ]
	costAfter := distanceMatrix[nodeI][nodeJ] + distanceMatrix[nextI][nextJ]

	return costAfter - costBefore
}


func deltaInterRouteExchange(solution []int, selectedIndex int, unselectedNode int, distanceMatrix [][]int) int {
	n := len(solution)

	prevSelected := solution[(selectedIndex-1+n)%n]
	nextSelected := solution[(selectedIndex+1)%n]

	costBefore := distanceMatrix[prevSelected][solution[selectedIndex]] + distanceMatrix[solution[selectedIndex]][nextSelected]
	costAfter := distanceMatrix[prevSelected][unselectedNode] + distanceMatrix[unselectedNode][nextSelected]

	return costAfter - costBefore
}

func applyMove(solution []int, move Move, unselectedNodes *[]int) {
	switch move.moveType {
	case "twoNodesExchange":
		solution[move.i], solution[move.j] = solution[move.j], solution[move.i]
	case "twoEdgesExchange":
		solution[move.i+1], solution[move.j] = solution[move.j], solution[move.i+1]
	case "interRouteExchange":
		solution[move.i], (*unselectedNodes)[move.j] = (*unselectedNodes)[move.j], solution[move.i]
	}
}

// GreedyMove evaluates moves until an improvement is found
func GreedyMove(solution []int, visited map[int]bool, unselectedNodes []int, distanceMatrix [][]int) []int {
	moves := generateMoves(solution, visited, unselectedNodes)

	for _, move := range moves {
		var delta int
		switch move.moveType {
		case "twoNodesExchange":
			delta = deltaTwoNodesExchange(solution, move.i, move.j, distanceMatrix)
		case "twoEdgesExchange":
			delta = deltaTwoEdgesExchange(solution, move.i, move.j, distanceMatrix)
		case "interRouteExchange":
			delta = deltaInterRouteExchange(solution, move.i, move.j, distanceMatrix)
		}

		if delta < 0 {
			applyMove(solution, move, &unselectedNodes)
			return solution
		}
	}

	return solution
}

func SteepestMove(solution []int, visited map[int]bool, unselectedNodes []int, distanceMatrix [][]int) []int {
	moves := generateMoves(solution, visited, unselectedNodes)
	bestDelta := 0
	var bestMove Move

	// Evaluate all moves
	for _, move := range moves {
		var delta int
		switch move.moveType {
		case "twoNodesExchange":
			delta = deltaTwoNodesExchange(solution, move.i, move.j, distanceMatrix)
		case "twoEdgesExchange":
			delta = deltaTwoEdgesExchange(solution, move.i, move.j, distanceMatrix)
		case "interRouteExchange":
			delta = deltaInterRouteExchange(solution, move.i, move.j, distanceMatrix)
		}

		if delta < bestDelta {
			bestDelta = delta
			bestMove = move
		}
	}

	if bestDelta < 0 {
		applyMove(solution, bestMove, &unselectedNodes)
	}

	return solution
}
