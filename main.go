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
	"time"
)

var iterations = 200

type MethodFunc func([][]int, int) []int

var methodsMap = map[string]MethodFunc{
	"random":                                           methods.RandomSolution,
	"nearest_neighbour_end_only":                       methods.NearestNeighborEndOnly,
	"nearest_neighbour_flexible":                       methods.NearestNeighborFlexible,
	"greedy_cycle":                                     methods.GreedyCycle,
	"greedy2regret":                                    methods.GreedyTwoRegret,
	"greedy2regret_weights":                            methods.GreedyRegretWeight,
	"LS_random_greedy_intranode":                       local_search.RandomGreedyIntraNode,
	"LS_random_greedy_intraedge":                       local_search.RandomGreedyIntraEdge,
	"LS_random_steepest_intranode":                     local_search.RandomSteepestIntraNode,
	"LS_random_steepest_intraedge":                     local_search.RandomSteepestIntraEdge,
	"LS_nearest_neighbour_flexible_greedy_intranode":   local_search.NearestNeighbourFlexibleGreedyIntraNode,
	"LS_nearest_neighbour_flexible_greedy_intraedge":   local_search.NearestNeighbourFlexibleGreedyIntraEdge,
	"LS_nearest_neighbour_flexible_steepest_intranode": local_search.NearestNeighbourFlexibleSteepestIntraNode,
	"LS_nearest_neighbour_flexible_steepest_intraedge": local_search.NearestNeighbourFlexibleSteepestIntraEdge,
	"LS_candidates":                                    local_search.LS_Candidates,
	"LS_delta":                                         local_search.LS_Delta,
	"LS_multi":                                         local_search.MultiLocalSearch,
	"LS_iterative":                                     local_search.IterativeLocalSearch,
	"large_noLS":                                       local_search.LargeNeighbourhood,
	"large_LS":                                         local_search.LargeNeighbourhoodWithLS,
	"custom":                                           local_search.CustomMethod,
	"hybrid":                                           local_search.HybridEA,
}

type Results struct {
	BestSolution   []int     `json:"best_solution"`
	BestFitness    int       `json:"best_fitness"`
	WorstSolution  []int     `json:"worst_solution"`
	WorstFitness   int       `json:"worst_fitness"`
	AverageFitness float32   `json:"average_fitness"`
	ExecutionTime  []float64 `json:"execution_time"` // in seconds
}

/////////////////////////////////////////////////////////////////////////////////
// Variables for global convexity

// run with respective command depending on data used !!!
// go run main.go data/TSPA.csv global_convexity
var PATH_TO_BEST = "logs/NAMEOFFUNCTION/TSPA/results.json"

// go run main.go data/TSPB.csv global_convexity
// var PATH_TO_BEST = "logs/NAMEOFFUNCTION/TSPB/results.json"

// this runs all experiments for given dataset(change if needed)
var similarity_measures = []string{"common_nodes", "common_edges"}
var similarities_to = []string{"best", "average"}

// ///////////////////////////////////////////////////////////////////////////////
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

		cmd := exec.Command("python", "scripts/log_results.py", inputFile, tempFile.Name(), methodName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error running python script: %v", err)
		}
	} else if methodName == "global_convexity" {
		// Open the JSON file with best solution
		bestSolution, err := utils.LoadBestSolution(PATH_TO_BEST)
		if err != nil {
			log.Fatalf("Error loading best solution from %s: %v", PATH_TO_BEST, err)
		}

		for _, similarity_measure := range similarity_measures {
			for _, similarity_to := range similarities_to {
				fmt.Printf("Running global convexity with similarity measure: %s, similarity to: %s\n", similarity_measure, similarity_to)
				similarities, fitnesses := local_search.GlobalConvexityLS(costMatrix, bestSolution, similarity_measure, similarity_to)

				//save to json file
				results := map[string]interface{}{
					"similarities": similarities,
					"fitnesses":    fitnesses,
				}

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

				cmd := exec.Command("python", "scripts/log_results_glob_conv.py", inputFile, tempFile.Name(), methodName, similarity_measure, similarity_to)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Fatalf("Error running python script: %v", err)
				}
			}
		}
	} else {
		log.Fatalf("Unknown method: %s", methodName)
	}
}

func runMethod(method MethodFunc, costMatrix [][]int) Results {
	var bestFitness, worstFitness, totalFitness int
	var bestSolution, worstSolution []int
	var times []float64

	for i := 0; i < iterations; i++ {
		startNode := i % len(costMatrix)

		timeIt := time.Now()
		solution := method(costMatrix, startNode)
		elapsed := time.Since(timeIt).Seconds()
		times = append(times, elapsed)
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
		ExecutionTime:  times,
	}
}

func parseArgs() (string, string) {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: go run main.go <data_file.csv> <method>| optional <num_iterations>\n")
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
