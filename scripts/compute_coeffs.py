import os
import json
from scipy.stats import pearsonr

# This script assumes a directory structure similar to:
# logs/global_convexity_1213/DATA/METHOD/results.json
# Where DATA could be TSPA or TSPB, and METHOD could be one of:
#   common_edges_best, common_edges_average, common_nodes_best, common_nodes_average
# You can adjust the BASE_DIR, data_list, and method_list as needed.

BASE_DIR = "logs/global_convexity_1213"
data_list = ["TSPA", "TSPB"]
method_list = [
    "common_edges_best",
    "common_edges_average",
    "common_nodes_best",
    "common_nodes_average"
]

def compute_correlation(fitnesses, similarities, remove_best=False):
    """Compute the Pearson correlation between fitnesses and similarities.
    If remove_best is True, remove the best solution (lowest fitness) from both arrays
    before computing the correlation."""
    if remove_best:
        # Find the index of the best (lowest) fitness
        best_idx = min(range(len(fitnesses)), key=lambda i: fitnesses[i])
        fitnesses = [f for i, f in enumerate(fitnesses) if i != best_idx]
        similarities = [s for i, s in enumerate(similarities) if i != best_idx]

    # Compute Pearson correlation
    r, p_value = pearsonr(fitnesses, similarities)
    return r, p_value

def main():
    for data in data_list:
        for method in method_list:
            file_path = os.path.join(BASE_DIR, data, method, "results.json")

            if not os.path.exists(file_path):
                print(f"File not found: {file_path}")
                continue

            with open(file_path, 'r') as f:
                results = json.load(f)

            fitnesses = results.get("fitnesses", [])
            similarities = results.get("similarities", [])
            
            if not fitnesses or not similarities or len(fitnesses) != len(similarities):
                print(f"Data issue in {file_path}. Make sure fitnesses and similarities are valid.")
                continue

            # Determine if we need to remove the best solution 
            # For methods ending with '_best', we should remove the best solution
            remove_best = method.endswith("_best")

            r, p_value = compute_correlation(fitnesses, similarities, remove_best=remove_best)
            print(f"Data: {data}, Method: {method}")
            print(f"Pearson correlation: {r:.4f}, p-value: {p_value:.4e}\n")

if __name__ == "__main__":
    main()
