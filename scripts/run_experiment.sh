#!/bin/bash

# Define instances and methods
instances=("data/tspA.csv" "data/tspB.csv")
methods=("random" "nearest_neighbor_end_only" "nearest_neighbor_flexible" "greedy_cycle")

# Loop through instances and methods
for instance in "${instances[@]}"; do
    for method in "${methods[@]}"; do
        echo "Running for instance: $instance with method: $method"
        go run main.go "$instance" "$method"
    done
done
