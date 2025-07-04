#!/bin/bash
for i in $(seq 1 1); do
    # go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 3 --input_avb 7
    # go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 6 --input_avb 14
    # go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 9 --input_avb 21
    # go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 12 --input_avb 28
    # go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 15 --input_avb 35
    # go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 18 --input_avb 42

    go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 21 --input_avb 49
    go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 24 --input_avb 56
    go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 27 --input_avb 63
    go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 30 --input_avb 70
    go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 33 --input_avb 77
    go run main.go --topology_name typical_complex --plan_name osro --test_case 100 --important_can 20 --unimportant_can 100 --input_tsn 36 --input_avb 84
done
