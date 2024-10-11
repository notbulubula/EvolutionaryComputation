setlocal
@echo off

set instances=data/tspA.csv data/tspB.csv
set methods=random nearest_neighbor_end_only nearest_neighbor_flexible


for %%i in (%instances%) do (
    for %%m in (%methods%) do (
        echo Running for instance: %%i with method: %%m
        go run main.go %%i %%m
    )
)

endlocal