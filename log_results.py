import json
import os
import datetime as dt
import sys
import pandas as pd

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