package local_search

import (
	"fmt"
	"math/rand"
)

type Move struct {
	moveType string
	i, j     int // indices of nodes involved
}

func generateMoves(solution []int, unselectedNodes []int, intraMoveType string) []Move {
	var moves []Move
	n := len(solution)

	// Intra-route: two-nodes exchange
	if intraMoveType == "NodeExchange" {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				moves = append(moves, Move{"twoNodesExchange", i, j})
			}
		}
	}

	// Intra-route: two-edges exchange
	if intraMoveType == "EdgeExchange" {
		for i := 0; i < n; i++ {
			for j := i + 2; j < n; j++ {
				moves = append(moves, Move{"twoEdgesExchange", i, j})
			}
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

	//if i and j are neighbours
	if (i+1)%n == j {
		costBefore := distanceMatrix[prevI][solution[i]] + distanceMatrix[solution[i]][solution[j]] + distanceMatrix[solution[j]][nextJ] - distanceMatrix[nextJ][nextJ]
		costAfter := distanceMatrix[prevI][solution[j]] + distanceMatrix[solution[j]][solution[i]] + distanceMatrix[solution[i]][nextJ] - distanceMatrix[nextJ][nextJ]
		return costAfter - costBefore
	}
	if (j+1)%n == i {
		costBefore := distanceMatrix[prevJ][solution[j]] + distanceMatrix[solution[j]][solution[i]] + distanceMatrix[solution[i]][nextI] - distanceMatrix[nextI][nextI]
		costAfter := distanceMatrix[prevJ][solution[i]] + distanceMatrix[solution[i]][solution[j]] + distanceMatrix[solution[j]][nextI] - distanceMatrix[nextI][nextI]
		return costAfter - costBefore
	}

	costBefore := distanceMatrix[prevI][solution[i]] + distanceMatrix[solution[i]][nextI] - distanceMatrix[nextI][nextI] +
		distanceMatrix[prevJ][solution[j]] + distanceMatrix[solution[j]][nextJ] - distanceMatrix[nextJ][nextJ]

	costAfter := distanceMatrix[prevI][solution[j]] + distanceMatrix[solution[j]][nextI] - distanceMatrix[nextI][nextI] +
		distanceMatrix[prevJ][solution[i]] + distanceMatrix[solution[i]][nextJ] - distanceMatrix[nextJ][nextJ]

	return costAfter - costBefore
}

func deltaTwoEdgesExchange(solution []int, i int, j int, distanceMatrix [][]int) int {
	n := len(solution)

	nodeI := solution[i]
	nextI := solution[(i+1)%n]
	nodeJ := solution[j]
	nextJ := solution[(j+1)%n]

	costBefore := distanceMatrix[nodeI][nextI] - distanceMatrix[nextI][nextI] +
		distanceMatrix[nodeJ][nextJ] - distanceMatrix[nextJ][nextJ]

	costAfter := distanceMatrix[nodeI][nodeJ] - distanceMatrix[nodeJ][nodeJ] +
		distanceMatrix[nextI][nextJ] - distanceMatrix[nextJ][nextJ]

	return costAfter - costBefore
}

// i, j are indexes in solution NOT IN MATRIX
// calculates delta fitness that adds a i->j edge
// e.g. transforms 12345 -> 12354
func deltaTwoEdgesExchangeCandidate(solution []int, i int, j int, dM [][]int) int {
	n := len(solution)

	I := solution[i]
	J := solution[j]

	nextI := solution[(i+1)%n]
	prevJ := solution[(j-1+n)%n]
	nextJ := solution[(j+1)%n]

	return dM[I][J] + dM[J][nextI] + dM[prevJ][nextJ] - dM[I][nextI] - dM[prevJ][J] - dM[J][nextJ]

}

// Moves j to i+1 position
func moveAfter(slice []int, i int, j int) []int {
	if i == j || i < 0 || i >= len(slice) || j < 0 || j >= len(slice) {
		return slice
	}

	// Remove the element at index `j` and save it
	elem := slice[j]
	if j > i {
		// If j is after i, remove it before inserting at i+1
		slice = append(slice[:j], slice[j+1:]...)
		slice = append(slice[:i+1], append([]int{elem}, slice[i+1:]...)...)
	} else {
		// If j is before i, insert at i+1 first, then remove old j
		slice = append(slice[:i+1], append([]int{elem}, slice[i+1:]...)...)
		slice = append(slice[:j], slice[j+1:]...)
	}

	return slice
}

// Return delta
func deltaInterCandidate(solution []int, i int, j int, dM [][]int) int {
	n := len(solution)
	iInSolution := findIndex(solution, i)
	nextI := solution[(iInSolution+1)%n]
	nextNextI := solution[(iInSolution+2)%n]

	delta := dM[i][j] + dM[j][nextNextI] - dM[i][nextI] - dM[nextI][nextNextI]

	return delta
}

func deltaInterRouteExchange(solution []int, selectedIndex int, unselectedNode int, distanceMatrix [][]int) int {
	n := len(solution)

	prevSelected := solution[(selectedIndex-1+n)%n]
	nextSelected := solution[(selectedIndex+1)%n]

	costBefore := distanceMatrix[prevSelected][solution[selectedIndex]] +
		distanceMatrix[solution[selectedIndex]][nextSelected] -
		distanceMatrix[nextSelected][nextSelected]

	costAfter := distanceMatrix[prevSelected][unselectedNode] +
		distanceMatrix[unselectedNode][nextSelected] -
		distanceMatrix[nextSelected][nextSelected]

	return costAfter - costBefore
}

func applyMove(solution []int, move Move, unselectedNodes *[]int) {
	debugPrints := false
	if debugPrints {
		fmt.Println("Solution before move: ", solution)
		fmt.Println("Move: ", solution[move.i], move.j, move.moveType)
	}
	switch move.moveType {
	case "twoNodesExchange":
		solution[move.i], solution[move.j] = solution[move.j], solution[move.i]
	case "twoEdgesExchange":
		reverseSegment(solution, move.i+1, move.j)
	case "interRouteExchange":
		// Find the index of `move.j` (unique node) in `unselectedNodes`
		jIndex := findIndex(*unselectedNodes, move.j)
		if jIndex != -1 {
			// Perform the swap using the found index
			solution[move.i], (*unselectedNodes)[jIndex] = (*unselectedNodes)[jIndex], solution[move.i]
		}
	}
	if debugPrints {
		fmt.Println("Solution after move: ", solution)
	}
}

func reverseSegment(solution []int, i, j int) {
	for i < j {
		solution[i], solution[j] = solution[j], solution[i]
		i++
		j--
	}
}

func findIndex(slice []int, value int) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

// GreedyMove evaluates moves until an improvement is found
func GreedyMove(solution []int, visited map[int]bool, unselectedNodes []int, distanceMatrix [][]int, intraMoveType string) ([]int, bool) {
	moves := generateMoves(solution, unselectedNodes, intraMoveType)
	improved := false

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
			improved = true
			return solution, improved
		}
	}

	return solution, improved
}

func SteepestMove(solution []int, visited map[int]bool, unselectedNodes []int, distanceMatrix [][]int, intraMoveType string) ([]int, bool) {
	moves := generateMoves(solution, unselectedNodes, intraMoveType)
	bestDelta := 0
	improved := false
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
		improved = true
		return solution, improved
	}

	return solution, improved
}

type CandidateMove struct {
	moveType             string
	i, j                 int
	solutionI, solutionJ int
}

func SteepestCandidate(solution []int,
	visited map[int]bool,
	unselectedNodes []int,
	distanceMatrix [][]int,
	moves []Move,
) ([]int, bool) {

	bestDelta := 0
	improved := false
	var bestMove Move

	for _, move := range moves {
		var delta int

		tempMove := CandidateMove{
			moveType:  move.moveType,
			i:         move.i,
			j:         move.j,
			solutionI: findIndex(solution, move.i),
			solutionJ: findIndex(solution, move.j),
		}

		if tempMove.solutionI == -1 && tempMove.solutionJ == -1 {
			continue
		}

		if tempMove.solutionJ == -1 {
			// inter move here
			tempMove.moveType = "interRouteExchange"
			delta = deltaInterCandidate(solution, tempMove.i, tempMove.j, distanceMatrix)
		} else if tempMove.solutionI == -1{
			tempMove.moveType = "interRouteExchange"

			tempMove.i, tempMove.j = tempMove.j, tempMove.i
			tempMove.solutionI, tempMove.solutionJ = tempMove.solutionJ, tempMove.solutionI

			delta = deltaInterCandidate(solution, tempMove.i, tempMove.j, distanceMatrix)
		} else {
			// if i and j are neighbours, skip
			if tempMove.solutionJ == (tempMove.solutionI+1)%len(solution) || tempMove.solutionI == (tempMove.solutionJ+1)%len(solution) {
				continue
			}
			tempMove.moveType = "twoEdgesExchange"

			delta = deltaTwoEdgesExchange(solution, tempMove.solutionI, tempMove.solutionJ, distanceMatrix)
		}

		if delta < bestDelta {
			bestDelta = delta
			if tempMove.moveType == "interRouteExchange" {
				bestMove = Move{moveType: tempMove.moveType, i: (tempMove.solutionI + 1) % len(solution), j: tempMove.j}
			} else {
				if tempMove.solutionI > tempMove.solutionJ {
					bestMove = Move{moveType: tempMove.moveType, i: tempMove.solutionJ, j: tempMove.solutionI}
				} else {
					bestMove = Move{moveType: tempMove.moveType, i: tempMove.solutionI, j: tempMove.solutionJ}
				}
			}
		}
	}

	if bestDelta < 0 {
		applyMove(solution, bestMove, &unselectedNodes)
		improved = true
		return solution, improved
	}

	return solution, improved
}
