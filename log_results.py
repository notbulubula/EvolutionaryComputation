import os
import sys
import json

import numpy as np
import pandas as pd
import datetime as dt
import matplotlib.pyplot as plt
import matplotlib.colors as mcolors

def plot_solution(nodes, solution, title, save_path):
    """
    nodes: DataFrame with columns x, y, cost
    solution: list of indexes to form a hamiltonian cycle
    title: title of the plot
    save_path: path to save the plot
    """
    plt.figure(figsize=(10, 8))

    # Set default opacity for nodes not in the solution
    unused_opacity = 0.2

    # Add opacity to nodes that are not part of the solution
    all_node_indexes = set(range(len(nodes)))
    solution_node_indexes = set(solution)
    unused_nodes = list(all_node_indexes - solution_node_indexes)
    
    # Scatter all nodes, reducing opacity for unused nodes
    plt.scatter(nodes.x[unused_nodes], nodes.y[unused_nodes], c="palevioletred", alpha=unused_opacity)
    plt.scatter(nodes.x[solution], nodes.y[solution], c="palevioletred", alpha=1.0)  # Full opacity for used nodes

    # Pre-calculate the total costs (Euclidean distance + node costs)
    costs = []
    for i in range(len(solution)):
        node_a = solution[i]
        node_b = solution[(i + 1) % len(solution)]  # Connects last node to the first

        # Calculate Euclidean distance between the two nodes
        dist = np.sqrt((nodes.x[node_a] - nodes.x[node_b]) ** 2 + (nodes.y[node_a] - nodes.y[node_b]) ** 2)

        # Total cost = Euclidean distance + cost of the two nodes
        total_cost = dist + nodes.cost[node_a] + nodes.cost[node_b]
        costs.append(total_cost)

    # Normalize costs for the color mapping
    min_cost = min(costs)
    max_cost = max(costs)

    # Create a colormap (light to red)
    norm = mcolors.Normalize(vmin=min_cost, vmax=max_cost)
    cmap = plt.get_cmap('summer_r')  # 'Reds' goes from light to red

    # Plot the solution's paths with colors based on the pre-calculated costs
    for i in range(len(solution)):
        node_a = solution[i]
        node_b = solution[(i + 1) % len(solution)]

        total_cost = costs[i]

        # Map the cost to a color in the heatmap
        color = cmap(norm(total_cost))

        # Plot the edge with the heatmap color
        plt.plot([nodes.x[node_a], nodes.x[node_b]], [nodes.y[node_a], nodes.y[node_b]], c=color)

    plt.title(title)
    ax = plt.gca()
    sm = plt.cm.ScalarMappable(norm=norm, cmap=cmap)
    sm.set_array([])
    plt.colorbar(sm, ax=ax, label="Cost (Euclidean + Node Costs)")
    plt.tight_layout()

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
timestamp = dt.datetime.now().strftime("%m%d")
file = data_file.replace(".csv", "").replace("data/", "")
current_folder = f"logs/{method}_{timestamp}/{file}"
os.makedirs(current_folder, exist_ok=True)

# Save results to results.json
results["method"] = method
results["timestamp"] = timestamp
results["best_solution"] = results.pop("best_solution")
results["worst_solution"] = results.pop("worst_solution")

with open(f"{current_folder}/results.json", "w") as f:
    json.dump(results, f, indent=4)

plot_solution(
    nodes, 
    results["best_solution"], 
    f"{method.upper()}: best solution({results['best_fitness']})", 
    f"{current_folder}/best_solution.png",
    )

plot_solution(
    nodes, 
    results["worst_solution"], 
    f"{method.upper()}: worst solution({results['worst_fitness']})", 
    f"{current_folder}/worst_solution.png",
    )