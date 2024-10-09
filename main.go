package main

import (
	"evolutionary_computation/methods"
	"evolutionary_computation/utils"
	"fmt"
	"os"
)

const iterations = 200

// TODO: Change to distance matrix only
type MethodFunc func([]utils.Node, [][]int) []int

var methodsMap = map[string]MethodFunc{
	"random":                   methods.RandomSolution,
	"nearest_neighbor_end_only": methods.NearestNeighborEndOnly,
	"nearest_neighbor_flexible": methods.NearestNeighborFlexible,
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run main.go  <data_file.csv> <method>")
	}

	file := os.Args[1]
	method := os.Args[2]

	nodes, err := utils.LoadNodes(file)
	if err != nil {
		fmt.Printf("Error loading nodes from %s: %v", file, err)
	}

	distanceMatrix := utils.CalculateDistanceMatrix(nodes)

	// Check if the provided method exists in the map
	if methodFunc, ok := methodsMap[method]; ok {
		runMethod(methodFunc, nodes, distanceMatrix)
	} else {
		fmt.Printf("Unknown method: %s. ", method)
	}
}

// runMethod handles running the method multiple times and calculating the stats
func runMethod(method MethodFunc, nodes []utils.Node, distanceMatrix [][]int) {
	var bestFitness, worstFitness, totalFitness int
	var bestSolution, worstSolution []int

	for i := 0; i < iterations; i++ {
		solution := method(nodes, distanceMatrix)
		fitness := utils.Fitness(solution, nodes, distanceMatrix)

		// Update best/worst fitness
		if i == 0 || fitness < bestFitness {
			bestFitness = fitness
			bestSolution = solution
		}
		if i == 0 || fitness > worstFitness {
			worstFitness = fitness
			worstSolution = solution
		}
		totalFitness += fitness
	}

	averageFitness := float32(totalFitness) / float32(iterations)

	fmt.Printf("Best solution (node indices): %v\nBest fitness: %v\n", bestSolution, bestFitness)
	fmt.Printf("Worst solution (node indices): %v\nWorst fitness: %v\n", worstSolution, worstFitness)
	fmt.Printf("Average fitness: %f\n", averageFitness)

	// TODO: Call Python script for plotting results
}