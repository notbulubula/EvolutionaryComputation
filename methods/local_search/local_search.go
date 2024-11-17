package local_search

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
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
		} else if tempMove.solutionI == -1 {
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

func SteepestDelta(solution []int, visited map[int]bool, unselectedNodes []int, distanceMatrix [][]int, moves []MoveDelta) ([]int, bool) {
	bestDelta := 0
	improved := false
	var bestMove MoveDelta
	var move_i_index int
	var move_j_index int
	var moveType string

	for i, move := range moves {
		delta := move.delta
		move_i_index = findIndex(solution, move.i)
		move_j_index = findIndex(solution, move.j)

		// removing unwanted moves (set delta to 0 in the moves list)
		if move.moveType == "twoEdgesExchange" &&
			((move_i_index == -1 || move_j_index == -1) || math.Abs(float64(move_i_index-move_j_index)) == 1) {
			moves[i].delta = 0
			continue
		}
		if move.moveType == "interRouteExchange" &&
			(move_i_index == -1 || move_j_index != -1) {
			moves[i].delta = 0
			continue
		}

		if delta < bestDelta {
			bestDelta = delta
			bestMove = move
			moveType = move.moveType

			// remove move from moves and break (set delta to 0 in the moves list)
			moves[i].delta = 0
			break
		}
	}

	if bestDelta < 0 {
		if moveType == "twoEdgesExchange" {
			// fmt.Println("Best move: ", bestMove, move_i_index, move_j_index)
			min := int(math.Min(float64(move_i_index), float64(move_j_index)))
			max := int(math.Max(float64(move_i_index), float64(move_j_index)))
			applyMove(solution, Move{moveType: moveType, i: min, j: max}, &unselectedNodes)
		} else {
			// fmt.Println("Best move: ", bestMove, move_i_index, move_j_index)
			applyMove(solution, Move{moveType: moveType, i: move_i_index, j: bestMove.j}, &unselectedNodes)
		}
		improved = true

		// Update the moves list
		updateMovesDelta(bestMove, moves, distanceMatrix, solution, move_i_index, move_j_index)

		return solution, improved
	}

	return solution, improved
}

func updateMovesDelta(bestMove MoveDelta, moves []MoveDelta, distanceMatrix [][]int, solution []int, move_i_index int, move_j_index int) {
	var NodestoCheck []int
	if bestMove.moveType == "twoEdgesExchange" {
		min := (move_i_index - 1 + len(solution)) % len(solution)
		max := (move_j_index + 2) % len(solution)
		min = int(math.Min(float64(max), float64(min)))
		max = int(math.Max(float64(max), float64(min)))
		for i := min; i < max; i++ {
			NodestoCheck = append(NodestoCheck, solution[i])
		}
	} else {
		j_index := findIndex(solution, bestMove.j)
		nextJ := solution[(j_index+1)%len(solution)]
		prevJ := solution[(j_index-1+len(solution))%len(solution)]
		NodestoCheck = []int{prevJ, bestMove.j, nextJ}
	}

	for M, move := range moves {
		if move.moveType == "twoEdgesExchange" &&
			(contains(NodestoCheck, move.i) && contains(solution, move.j)) ||
			(contains(NodestoCheck, move.j) && contains(solution, move.i)) {
			//check if i and j are not neighbours
			i_index := findIndex(solution, move.i)
			j_index := findIndex(solution, move.j)

			min := int(math.Min(float64(i_index), float64(j_index)))
			max := int(math.Max(float64(i_index), float64(j_index)))

			if math.Abs(float64(i_index-j_index)) != 1 {
				delta := deltaTwoEdgesExchange(solution, min, max, distanceMatrix)
				if delta < 0 {
					moves[M].delta = delta
				} else {
					moves[M].delta = 0
				}
			} else {
				moves[M].delta = 0
			}
		}
		if move.moveType == "interRouteExchange" &&
			contains(NodestoCheck, move.i) && !contains(solution, move.j) {

			i_index := findIndex(solution, move.i)
			delta := deltaInterRouteExchange(solution, i_index, move.j, distanceMatrix)
			if delta < 0 {
				moves[M].delta = delta
			} else {
				moves[M].delta = 0
			}
		}

	}

	// Sort moves by delta in ascending order
	sort.Slice(moves, func(a, b int) bool {
		return moves[a].delta < moves[b].delta
	})
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
