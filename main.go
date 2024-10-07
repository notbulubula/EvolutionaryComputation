package main

import (
	"evolutionary_computation/methods"
	"evolutionary_computation/utils"
	"fmt"
	"log"
	"os"
)

// Enum for methods
const (
	Random                 = "random"
	NearestNeighborEndOnly = "nearest_neighbor_end_only"
	NearestNeighborFlexible = "nearest_neighbor_flexible"
	GreedyCycle            = "greedy_cycle"
)

func main() {
	// Check for method and file arguments
	if len(os.Args) < 3 {
		log.Fatalf("Usage: go run main.go  <data_file.csv> <method>")
	}
	
	file := os.Args[1]
	method := os.Args[2]

	// Load node data from CSV
	nodes, err := utils.LoadNodes(file)
	if err != nil {
		log.Fatalf("Error loading nodes from %s: %v", file, err)
	}

	// Calculate the distance matrix
	distanceMatrix := utils.CalculateDistanceMatrix(nodes)

	// Execute the chosen method
	switch method {
	case Random:
		// TODO: This should be a function call with callable methods.
		// Should do the computation x times and return best/worst/average results. 
		// Afterwards it should call python script to plot the results.
		solution := methods.RandomSolution(nodes, distanceMatrix)
		fitness := utils.Fitness(solution, nodes, distanceMatrix)

		fmt.Println("Random solution (node indices):", solution)
		fmt.Println("Fitness:", fitness)
	default:
		log.Fatalf("Unknown method: %s. ", method)
	}
}