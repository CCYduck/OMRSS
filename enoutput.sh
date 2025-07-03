#!/bin/bash
for i in $(seq 1 1); do
    go run main.go --topology_name typical_complex --plan_name osro --test_case 50 --important_can 20 --unimportant_can 100 
    go run main.go --topology_name typical_complex --plan_name osro --test_case 50 --important_can 30 --unimportant_can 150 
    go run main.go --topology_name typical_complex --plan_name osro --test_case 50 --important_can 40 --unimportant_can 200 
    go run main.go --topology_name typical_complex --plan_name osro --test_case 50 --important_can 50 --unimportant_can 250 
    go run main.go --topology_name typical_complex --plan_name osro --test_case 50 --important_can 60 --unimportant_can 300 
    go run main.go --topology_name typical_complex --plan_name osro --test_case 50 --important_can 70 --unimportant_can 350 

   
done
