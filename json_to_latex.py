import os
import json

root_dir = 'logs'

def process_results_json(json_path, tex_path):
    with open(json_path, 'r') as json_file:
        data = json.load(json_file)

    best_fitness = data.get("best_fitness", "N/A")
    worst_fitness = data.get("worst_fitness", "N/A")
    average_fitness = data.get("average_fitness", "N/A")
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

\\textbf{{Average Fitness: {average_fitness}}}
"""

    with open(tex_path, 'w') as tex_file:
        tex_file.write(latex_content)

    print(f'Generated LaTeX file: {tex_path}')

for dirpath, dirnames, filenames in os.walk(root_dir):
    for filename in filenames:
        if filename == 'results.json':
            json_file_path = os.path.join(dirpath, filename)
            tex_file_path = os.path.join(dirpath, 'results.tex')
            process_results_json(json_file_path, tex_file_path)

print("Finished generating LaTeX files.")
