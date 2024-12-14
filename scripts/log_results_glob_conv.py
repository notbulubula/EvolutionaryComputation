import os
import sys
import json
import numpy as np
import datetime as dt
import matplotlib.pyplot as plt

def plot_similarity_vs_fitness(fitnesses, similarities, title, save_path):
    """
    Generates a plot of the average similarity vs fitness values.

    Parameters:
    fitnesses (list of floats): Objective function values (x-axis).
    similarities (list of floats): Average similarity values (y-axis).
    title (str): Title of the chart.
    save_path (str): Path to save the chart image.
    """
    plt.figure(figsize=(8, 6))
    
    # Scatter plot for fitness vs similarity
    plt.scatter(fitnesses, similarities, color="palevioletred", alpha=0.7, label="Data Points")

    plt.title(title)
    plt.xlabel("Objective Function Value (Fitness)")
    plt.ylabel("Similarity")
    plt.grid(True, linestyle="--", alpha=0.6)

    plt.tight_layout()

    plt.savefig(save_path)
    plt.close()

if len(sys.argv) != 6:
    print("Usage: python log_results_glob_conv.py <data_file> <results.json> <method> <similarity_measure> <similarity_to>")
    sys.exit(1)

data_file = sys.argv[1]
results_file = sys.argv[2]
method = sys.argv[3]
similarity_measure = sys.argv[4]
similarity_to = sys.argv[5]

results = json.load(open(results_file))

# folder management
timestamp = dt.datetime.now().strftime("%m%d")
file = data_file.replace(".csv", "").replace("data/", "")
current_folder = f"logs/{method}_{timestamp}/{file}/{similarity_measure}_{similarity_to}"
os.makedirs(current_folder, exist_ok=True)

# Save results to results.json
results["method"] = method
results["timestamp"] = timestamp
results["fitnesses"] = results.pop("fitnesses")
results["similarities"] = results.pop("similarities")

with open(f"{current_folder}/results.json", "w") as f:
    json.dump(results, f, indent=4)

plot_similarity_vs_fitness(
    results["fitnesses"], 
    results["similarities"], 
    f"{method.upper()}: {similarity_measure} to {similarity_to}", 
    f"{current_folder}/similarity_vs_fitness.png",
    )

