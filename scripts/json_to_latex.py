import os
import json

# Set the logs directory path
root_dir = r'C:\Users\Uni\Documents\Uni\sem7\EvolutionaryComputation\EvolutionaryComputation\logs'


def process_results_json(json_path, tex_path):
    try:
        with open(json_path, 'r') as json_file:
            data = json.load(json_file)

        best_fitness = data.get("best_fitness", "N/A")
        worst_fitness = data.get("worst_fitness", "N/A")
        average_fitness = data.get("average_fitness", "N/A")
        execution_time = data.get("execution_time", "N/A")  # New line to get execution time
        best_solution = data.get("best_solution", [])
        worst_solution = data.get("worst_solution", [])
        method = data.get("method", "unknown").replace('_', ' ').title()

        latex_content = f"""
\\subsubsection{{Results for {method}}}
\\textbf{{Best Fitness: {best_fitness}}}

\\begin{{lstlisting}}
{', '.join(map(str, best_solution))}
\\end{{lstlisting}}

\\textbf{{Worst Fitness: {worst_fitness}}}

\\begin{{lstlisting}}
{', '.join(map(str, worst_solution))}
\\end{{lstlisting}}

\\textbf{{Average Fitness: {average_fitness}}} \\\\
\\textbf{{Execution Time: {execution_time} seconds}}
"""

        with open(tex_path, 'w') as tex_file:
            tex_file.write(latex_content)

        print(f'Generated LaTeX file: {tex_path}')

    except Exception as e:
        print(f"Error processing file {json_path}: {e}")

# Print root directory and check existence
print(f"Generating LaTeX files in directory: {root_dir}")
print("Absolute path of root_dir:", os.path.abspath(root_dir))
if not os.path.isdir(root_dir):
    print(f"Directory does not exist: {root_dir}")
else:
    print("Contents of the logs directory:", os.listdir(root_dir))

    # Traverse directories and look for results.json files
    for dirpath, dirnames, filenames in os.walk(root_dir):
        print(f'Inspecting directory: {dirpath}')
        print(f'Found files: {filenames}')
        for filename in filenames:
            if filename == 'results.json':
                print(f'Processing {filename} in {dirpath}...')
                json_file_path = os.path.join(dirpath, filename)
                tex_file_path = os.path.join(dirpath, 'results.tex')
                process_results_json(json_file_path, tex_file_path)

print("Finished generating LaTeX files.")
