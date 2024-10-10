import json
import os
import datetime as dt
import sys
import pandas as pd
import matplotlib.pyplot as plt

def plot_solution(nodes, solution, title, save_path):
    """
    nodes: DataFrame with columns x, y, cost
    solution: list of indexes to form a hamiltonian cycle
    title: title of the plot
    save_path: path to save the plot
    """
    plt.figure(figsize=(10, 10))
    plt.scatter(nodes.x, nodes.y, c="#5b7ebd")
    for i in range(len(solution) - 1):
        plt.plot([nodes.x[solution[i]], nodes.x[solution[i + 1]]], [nodes.y[solution[i]], nodes.y[solution[i + 1]]], c="black")
    plt.plot([nodes.x[solution[-1]], nodes.x[solution[0]]], [nodes.y[solution[-1]], nodes.y[solution[0]]], c="black")
    plt.title(title)
    plt.savefig(save_path)
    plt.close()


if len(sys.argv) != 4:
    print("Usage: python log_results.py <data.csv> <results.json> <method>")
    sys.exit(1)

data_file = sys.argv[1]
results_file = sys.argv[2]
method = sys.argv[3]

# Data loading
nodes = pd.read_csv(data_file, header=None, sep=";")
nodes.columns = ["x", "y", "cost"]

results = json.load(open(results_file, "r"))

# folder management
timestamp = dt.datetime.now().strftime("%m%d_%H%M")
current_folder = f"logs/{method}_{timestamp}"
os.makedirs(current_folder, exist_ok=True)

# Save results to results.json
results["method"] = method
results["timestamp"] = timestamp
results["best_solution"] = results.pop("best_solution")
results["worst_solution"] = results.pop("worst_solution")

with open(f"{current_folder}/results.json", "w") as f:
    json.dump(results, f, indent=4)

plot_solution(nodes, results["best_solution"], f"{method.upper()}: best solution", f"{current_folder}/best_solution.png")
plot_solution(nodes, results["worst_solution"], f"{method.upper()}: worst solution", f"{current_folder}/worst_solution.png")