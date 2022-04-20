# Fault-Tolerant-Agreement
Implementation of the [**Galil-Mayer-Yung**](https://ieeexplore.ieee.org/abstract/document/492674) algorithm for Crash Fault Tolerant Broadcast with linear message complexity.

We also compare the performance (message complexity) of the GMY algorithm against [Dolev-Strong Algorithm](https://www.semanticscholar.org/paper/Authenticated-Algorithms-for-Byzantine-Agreement-Dolev-Strong/38c830bf6192d9e83cf6793d01c54032b63bb8f8).

## Details
This implementation is done as part of a course project for the graduate-level course 'Fault Tolerant Distributed Algorithms (CS598FTD)' taught by [**Prof. Ling Ren**](https://sites.google.com/view/renling). \
The final project report can be found [here](http://link_to_report).

## Description of Modules
### Main (Central server)
+ Initialisation of nodes and blackbox
+ During initialisation, it choose a random subset of atmost f nodes and sets a non-zero probability of failure for them.
+ Sanity Checks

### Blackbox
+ Maintains lockstep-synchrony by sending a `START_PHASE` message to each alive node at the start of next phase
+ [Checkpointing](https://dl.acm.org/doi/10.1145/197917.198082) (Phase 2 and Phase 4)
+ Recruitment of Coordinators (Phase 5)
+ Termination

### Node
+ Carries out the functions of a node as per the algorithm
+ If failureProbability>0, it fails at a random time during the run of algorithm. It also sends a message to blackbox before failing.

## Usage
```bash
go run src/main.go n f p
```
> `n` : (int) number of nodes \
> `f` : (int) maximum number of failures \
> `p` : (float) probability of failure for each of the `f` nodes

Example: ```go run src/main.go 500 400 0.9```

## Main Files
+ `src/main.go` - Top level module that initialises the nodes and blackbox. It triggers the algorithm and conducts sanity checks once it       terminates.
+ `src/node/node.go` - Node module
+ `src/blackbox/blackbox.go` - Blackbox module
+ `src/types/types.go` - This file contains all the structs and enums used in the code.
+ `scripts/expt.sh` - Bash script to run large scale experiments to generate empirical data. \
    Running Instructions: 
    ```./scripts/expt.sh csvfile logdir``` \
    `csvfile` : Output CSV file for recording the experiment data. Note that the `csvfile` should not already exist. This is done to prevent an existing csv from being accidentally overwritten. \
    `logdir` : Directory in which error logs from the experiment are stored. Note that only the logs of the experiments indicating any error are saved. The logs of the successful experiments are overwritten to optimise space. 
+ `scripts/analyse.py` - Python script to compare the performance of Galil-Mayer-Yung algorithm with Dolev-Strong algorithm. It plots the      graphs for multiple settings and saves them to the `plots` directory.

## Message Complexity
### Galil-Mayer-Yung Algorithm
We empirically count the number of messages exchanged during the execution of the algorithm.
### Dolev-Strong Algorithm
+ messages sent by Honest(Alive) Nodes = (n-f)*(n-1)
+ messages sent by Faulty Nodes = f*(n-1)
+ total messages sent = n*(n-1)

## Contributors
+ [Samarth Aggarwal](https://github.com/samarthaggarwal)
+ [Shivram Gowtham](https://github.com/ShivramIITG)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)
