package utils

import "math"

// CalculateDistanceMatrix returns a matrix with rounded Euclidean distances.
func CalculateDistanceMatrix(nodes []Node) [][]int {
	numNodes := len(nodes)
	matrix := make([][]int, numNodes)

	for i := 0; i < numNodes; i++ {
		matrix[i] = make([]int, numNodes)
		for j := 0; j < numNodes; j++ {
			matrix[i][j] = CalculateDistance(nodes[i], nodes[j])
		}
	}
	return matrix
}

// CalculateDistance computes the Euclidean distance between two nodes, rounded to the nearest integer.
func CalculateDistance(a, b Node) int {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return int(math.Round(math.Sqrt(dx*dx + dy*dy)))
}
