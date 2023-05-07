# MST-ACOvsOther
Simulate TSN and AVB data streams to find the best path in the topology.

## Installation
Clone this repo by `git clone https://github.com/helgesander02/MST-ACOvsOther`<br />

## Running
Quickstart: `go run main.go`<br />
More options:
| Option | Description |
| -------- | ---- | 
| --test_case | Conducting n experiments |
| --tsn | Number of TSN flows |
| --avb | Number of AVB flows |
| --HyperPeriod | Greatest Common Divisor of Simulated Time LCM |
| --topology_print | Display all topology information |
| --flows_print | Display all Flows information |


## Reference
[Bang Ye Wu, Kun-Mao Chao, "Steiner Minimal Trees"](https://www.csie.ntu.edu.tw/~kmchao/tree10spr/Steiner.pdf)<br />
[geeksforgeeks, "Steiner Tree Problem"](https://www.geeksforgeeks.org/steiner-tree/)<br />
[Sharon Peng, "Kruskal’s Algorithm"](https://mycollegenotebook.medium.com/kruskal-algorithm-deb0d64ce271)<br />
