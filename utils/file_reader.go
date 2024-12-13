package utils

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strconv"
)

type Node struct {
	ID, X, Y, Cost int
}

// LoadNodes reads the CSV file and returns a slice of Node structs, each with an ID.
func LoadNodes(filename string) ([]Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nodes []Node
	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		x, _ := strconv.Atoi(record[0])
		y, _ := strconv.Atoi(record[1])
		cost, _ := strconv.Atoi(record[2])

		nodes = append(nodes, Node{ID: i, X: x, Y: y, Cost: cost})
	}

	return nodes, nil
}

func LoadBestSolution(filename string) ([]int, error) {
	// Open the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Read the contents of the file
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// Create a struct to represent the JSON data
	var data struct {
		BestSolution []int `json:"best_solution"`
	}
	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data.BestSolution, nil
}
