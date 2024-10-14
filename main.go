package main

import (
	"encoding/json"
	"evolutionary_computation/methods"
	"evolutionary_computation/utils"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var iterations = 200

// TODO: Change to distance matrix only
type MethodFunc func([][]int, int) []int

var methodsMap = map[string]MethodFunc{
	"random":                    methods.RandomSolution,
	"nearest_neighbor_end_only": methods.NearestNeighborEndOnly,
	"nearest_neighbor_flexible": methods.NearestNeighborFlexible,
	"greedy_cycle":              methods.GreedyCycle,
	"greedy2regret":             methods.GreedyTwoRegret,
}

type Results struct {
	BestSolution   []int   `json:"best_solution"`
	BestFitness    int     `json:"best_fitness"`
	WorstSolution  []int   `json:"worst_solution"`
	WorstFitness   int     `json:"worst_fitness"`
	AverageFitness float32 `json:"average_fitness"`
}

func main() {
	var file, method string
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run main.go  <data_file.csv> <method>| optional <num_iterations>")
	} else if len(os.Args) == 4 {
		i, err := strconv.Atoi(os.Args[3])

		if err != nil {
			fmt.Printf("Couldn't convert num iterations to int %v", err)
		}
		iterations = i
	}
	file = os.Args[1]
	method = os.Args[2]

	nodes, err := utils.LoadNodes(file)
	if err != nil {
		fmt.Printf("Error loading nodes from %s: %v", file, err)
	}

	costMatrix := utils.CalculateCostMatrix(nodes)

	// If method exists, run it
	if methodFunc, ok := methodsMap[method]; ok {
		results := runMethod(methodFunc, costMatrix)

		// Parse to JSON for Python to handle
		jsonResults, err := json.Marshal(results)
		if err != nil {
			fmt.Printf("Error marshalling results: %v", err)
			return
		}

		tempFile, err := os.CreateTemp("", "results.json")
		if err != nil {
			fmt.Printf("Error creating temp file: %v", err)
			return
		}
		defer tempFile.Close()

		if _, err := tempFile.Write(jsonResults); err != nil {
			fmt.Printf("Error writing to temp file: %v", err)
			return
		}

		cmd := exec.Command("python3", "scripts/log_results.py", file, tempFile.Name(), method)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running python script: %v", err)
			return
		}
	} else {
		fmt.Printf("Unknown method: %s. ", method)
	}
}

// runMethod handles running the method multiple times and calculating the stats
func runMethod(method MethodFunc, costMatrix [][]int) Results {
	var bestFitness, worstFitness, totalFitness int
	var bestSolution, worstSolution []int

	for i := 0; i < iterations; i++ {
		startNode := i % len(costMatrix)

		solution := method(costMatrix, startNode)
		fitness := utils.Fitness(solution, costMatrix)

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

	return Results{
		BestSolution:   bestSolution,
		BestFitness:    bestFitness,
		WorstSolution:  worstSolution,
		WorstFitness:   worstFitness,
		AverageFitness: averageFitness,
	}
}
