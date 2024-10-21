import os
import json
import pandas as pd
import xlwings as xw

excel_checker_path = 'Solution checker.xlsx'
root_dir = 'logs'

def check_solution_in_excel(excel_path, json_data, instance_name, method_name):
    app = xw.App(visible=False)
    wb = xw.Book(excel_path)
    sheet_name = 'TSPA' if 'tspA' in instance_name else 'TSPB'
    sheet = wb.sheets[sheet_name]
    
    # Start inserting the nodes in the 'List of nodes' column (column F), starting at row 3
    for i, node_id in enumerate(json_data['best_solution']):
        sheet.range(f'F{i + 3}').value = node_id  # Insert each node in column F, starting from row 3
    
    # Force Excel to recalculate the formulas
    wb.app.calculate()
    excel_fitness = sheet.range('K2').value
    wb.close()
    app.quit()

    if excel_fitness is None:
        print(f"Warning: Could not retrieve 'Objective function' from {excel_path} for {method_name} on {instance_name}")
    
    return {
        'method': method_name,
        'instance': instance_name,
        'calculated fitness': json_data['best_fitness'],
        'excel fitness': excel_fitness,
        'match': json_data['best_fitness'] == excel_fitness
    }

results = []

for dirpath, dirnames, filenames in os.walk(root_dir):
    for filename in filenames:
        if filename == 'results.json':
            json_file_path = os.path.join(dirpath, filename)
            
            with open(json_file_path, 'r') as json_file:
                json_data = json.load(json_file)
            
            instance_name = os.path.basename(dirpath)
            method_name = json_data.get("method", "unknown").replace('_', ' ').title()
            
            result = check_solution_in_excel(excel_checker_path, json_data, instance_name, method_name)
            results.append(result)

results_df = pd.DataFrame(results)
comparison_csv_path = os.path.join(root_dir, 'fitness_comparison_table.csv')
results_df.to_csv(comparison_csv_path, index=False)
latex_output = results_df.to_latex(index=False, column_format='|l|l|r|r|c|', header=True, float_format="%.2f")
latex_output_path = os.path.join(root_dir, 'fitness_comparison_table.tex')

with open(latex_output_path, 'w') as latex_file:
    latex_file.write(latex_output)

print(results_df)
print(f"LaTeX table saved to: {latex_output_path}")
