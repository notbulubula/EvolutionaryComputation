setlocal
@echo off

set instances=data/tspA.csv data/tspB.csv
set methods=random nearest_neighbor_end_only nearest_neighbor_flexible, greedy_cycle
set ls_methods=LS_random_greedy_intranode LS_random_greedy_intraedge LS_random_steepest_intranode LS_random_steepest_intraedge LS_nearest_neighbour_flexible_greedy_intranode LS_nearest_neighbour_flexible_greedy_intraedge LS_nearest_neighbour_flexible_steepest_intranode LS_nearest_neighbour_flexible_steepest_intraedge


for %%i in (%instances%) do (
    for %%m in (%ls_methods%) do (
        echo Running for instance: %%i with method: %%m
        go run main.go %%i %%m
    )
)

endlocal