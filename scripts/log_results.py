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
    plt.figure(figsize=(10, 6))

    # Set default opacity for nodes not in the solution
    unused_opacity = 0.2

    # Normalize the sizes based on node costs
    min_cost = nodes.cost.min()
    max_cost = nodes.cost.max()

    # Calculate sizes based on normalized costs, using a size range (e.g., 50 to 500)
    size_range = (50, 500)  # Minimum and maximum size for the dots
    sizes = size_range[0] + (size_range[1] - size_range[0]) * (nodes.cost - min_cost) / (max_cost - min_cost)

    # Set sizes for all nodes
    # Scatter all unused nodes, reducing opacity
    unused_nodes = list(set(range(len(nodes))) - set(solution))
    plt.scatter(nodes.x[unused_nodes], nodes.y[unused_nodes], c="palevioletred", s=sizes[unused_nodes], alpha=unused_opacity)

    # Full opacity for used nodes with size based on their cost
    plt.scatter(nodes.x[solution], nodes.y[solution], c="palevioletred", s=sizes[solution], alpha=1.0)


    # Pre-calculate the vertices costs (Euclidean distance)
    costs = []
    for i in range(len(solution)):
        node_a = solution[i]
        node_b = solution[(i + 1) % len(solution)]  # Connects last node to the first

        # Calculate Euclidean distance between the two nodes
        dist = np.sqrt((nodes.x[node_a] - nodes.x[node_b]) ** 2 + (nodes.y[node_a] - nodes.y[node_b]) ** 2)

        costs.append(dist)

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
    ax.set_aspect('equal', adjustable='box')
    sm = plt.cm.ScalarMappable(norm=norm, cmap=cmap)
    sm.set_array([])
    plt.colorbar(sm, ax=ax, label="Cost (Euclidean distance)")
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