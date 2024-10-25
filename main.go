package main

import (
	"encoding/json"
	"evolutionary_computation/methods"
	"evolutionary_computation/methods/local_search"
	"evolutionary_computation/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var iterations = 200

type MethodFunc func([][]int, int) []int

var methodsMap = map[string]MethodFunc{
	"random":                     methods.RandomSolution,
	"nearest_neighbour_end_only": methods.NearestNeighborEndOnly,
	"nearest_neighbour_flexible": methods.NearestNeighborFlexible,
	"greedy_cycle":               methods.GreedyCycle,
	"greedy2regret":              methods.GreedyTwoRegret,
	"greedy2regret_weights":      methods.GreedyRegretWeight,
	"LS_random_greedy_intranode": local_search.RandomGreedyIntraNode,
	"LS_random_greedy_intraedge": local_search.RandomGreedyIntraEdge,
	// "random_steepest_intranode": local_search.RandomSteepestIntraNode,
	// "random_steepest_intraedge": local_search.RandomSteepestIntraEdge,
	"LS_nearest_neighbour_flexible_greedy_intranode": local_search.NearestNeighbourFlexibleGreedyIntraNode,
	"LS_nearest_neighbour_flexible_greedy_intraedge": local_search.NearestNeighbourFlexibleGreedyIntraEdge,
	// "nearest_neighbour_flexible_steepest_intranode": local_search.NearestNeighbourFlexibleSteepestIntraNode,
	// "nearest_neighbour_flexible_steepest_intraedge": local_search.NearestNeighbourFlexibleSteepestIntraEdge,
}

type Results struct {
	BestSolution   []int   `json:"best_solution"`
	BestFitness    int     `json:"best_fitness"`
	WorstSolution  []int   `json:"worst_solution"`
	WorstFitness   int     `json:"worst_fitness"`
	AverageFitness float32 `json:"average_fitness"`
}

func main() {
	inputFile, methodName := parseArgs()

	nodes, err := utils.LoadNodes(inputFile)
	if err != nil {
		log.Fatalf("Error loading nodes from %s: %v", inputFile, err)
	}

	costMatrix := utils.CalculateCostMatrix(nodes)

	if methodFunc, ok := methodsMap[methodName]; ok {
		results := runMethod(methodFunc, costMatrix)

		jsonResults, err := json.Marshal(results)
		if err != nil {
			log.Fatalf("Error marshalling results: %v", err)
		}

		tempFile, err := os.CreateTemp("", "results.json")
		if err != nil {
			log.Fatalf("Error creating temp file: %v", err)
		}
		defer tempFile.Close()

		if _, err := tempFile.Write(jsonResults); err != nil {
			log.Fatalf("Error writing to temp file: %v", err)
		}

		pythonCmd := detectPython()

		cmd := exec.Command(pythonCmd, "scripts/log_results.py", inputFile, tempFile.Name(), methodName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error running python script: %v", err)
		}
	} else {
		log.Fatalf("Unknown method: %s", methodName)
	}
}

func runMethod(method MethodFunc, costMatrix [][]int) Results {
	var bestFitness, worstFitness, totalFitness int
	var bestSolution, worstSolution []int

	for i := 0; i < iterations; i++ {
		startNode := i % len(costMatrix)

		solution := method(costMatrix, startNode)
		fitness := utils.Fitness(solution, costMatrix)

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

func detectPython() string {
	pythonCmd := "python3"
	if _, err := exec.LookPath("python3"); err != nil {
		if _, err := exec.LookPath("python"); err == nil {
			pythonCmd = "python"
		} else {
			log.Fatal("Error: Neither python3 nor python is installed or available in PATH.")
		}
	}
	return pythonCmd
}

func parseArgs() (string, string) {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: go run main.go  <data_file.csv> <method>| optional <num_iterations>\n")
	}

	file := os.Args[1]
	method := os.Args[2]

	if len(os.Args) == 4 {
		i, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("Couldn't convert num iterations to int: %v", err)
		}
		iterations = i
	}

	return file, method
}
