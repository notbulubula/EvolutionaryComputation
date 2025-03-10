package local_search

import (
	"evolutionary_computation/methods"
	"evolutionary_computation/utils"
	"math/rand"
	"time"
)

// HybridEA implements the hybrid evolutionary algorithm
func HybridEA(costMatrix [][]int, pointless_value int) []int {
	var bestFitness int
	var bestSolution []int

	EliteSize := 20
	MaxGenerations := 200

	startTime := time.Now()

	// Initialize elite population
	elitePopulation := initializePopulation(costMatrix, EliteSize)

	for time.Since(startTime) < 24*time.Second {
		for gen := 0; gen < MaxGenerations; gen++ {
			// Select parents
			parent1, parent2 := selectParents(elitePopulation)

			// Apply recombination
			offspring := recombine(parent1.Path, parent2.Path, costMatrix)
			// Perform local search
			offspring.Path = NearestNeighbourFlexibleSteepestIntraEdgeFromSolution(costMatrix, offspring.Path)

			offspring.Fitness = utils.Fitness(offspring.Path, costMatrix)

			// Check diversity and update elite population
			if !isDuplicate(elitePopulation, offspring) {
				elitePopulation = replaceWorst(elitePopulation, offspring)
			}
		}
		// Find the best solution in the elite population
		for _, solution := range elitePopulation {
			if bestSolution == nil || solution.Fitness < bestFitness {
				bestSolution = solution.Path
				bestFitness = solution.Fitness
			}
		}
	}

	return bestSolution
}

type HybridSolution struct {
	Path    []int
	Fitness int
}

func initializePopulation(costMatrix [][]int, size int) []HybridSolution {
	population := make([]HybridSolution, size)
	for i := 0; i < size; i++ {
		path := RandomSteepestIntraEdge(costMatrix, rand.Intn(len(costMatrix)))
		fitness := utils.Fitness(path, costMatrix)
		population[i] = HybridSolution{Path: path, Fitness: fitness}
	}
	return population
}

func selectParents(population []HybridSolution) (HybridSolution, HybridSolution) {
	parent1 := population[rand.Intn(len(population))]
	parent2 := population[rand.Intn(len(population))]

	// Ensure parents are different
	for parent1.Fitness == parent2.Fitness {
		parent2 = population[rand.Intn(len(population))]
	}

	return parent1, parent2
}

func recombine(parent1, parent2 []int, costMatrix [][]int) HybridSolution {
	if rand.Float64() < 0.6 {
		return recombineOperator1(parent1, parent2)
	}
	return recombineOperator2(parent1, parent2, costMatrix)
}

func recombineOperator1(parent1, parent2 []int) HybridSolution {
	child := make([]int, len(parent1))
	inChild := make(map[int]bool)

	// Add common nodes
	for i := range child {
		if utils.Contains(parent2, parent1[i]) {
			child[i] = parent1[i]
			inChild[parent1[i]] = true
		} else {
			child[i] = -1
		}
	}

	// Fill remaining nodes randomly
	for i := range child {
		if child[i] == -1 {
			for {
				//add random node that is not in solution from random parent solution
				node := parent1[rand.Intn(len(parent2))]
				if rand.Float64() < 0.5 {
					node = parent2[rand.Intn(len(parent2))]
				}

				if !inChild[node] {
					child[i] = node
					inChild[node] = true
					break
				}
			}
		}
	}

	return HybridSolution{Path: child}
}

func recombineOperator2(parent1, parent2 []int, costMatrix [][]int) HybridSolution {
	commonNodes := []int{}
	for _, node := range parent1 {
		if utils.Contains(parent2, node) {
			commonNodes = append(commonNodes, node)
		}
	}

	// Perform nearest neighbor repair starting with the common nodes
	repairedChild := methods.NearestNeighborFlexibleFromSolution(costMatrix, commonNodes)

	return HybridSolution{Path: repairedChild}
}

func isDuplicate(population []HybridSolution, offspring HybridSolution) bool {
	for _, s := range population {
		if s.Fitness == offspring.Fitness {
			return true
		}
	}
	return false
}

func replaceWorst(population []HybridSolution, offspring HybridSolution) []HybridSolution {
	worstIdx := 0
	for i, solution := range population {
		if solution.Fitness > population[worstIdx].Fitness {
			worstIdx = i
		}
	}
	population[worstIdx] = offspring

	return population
}
