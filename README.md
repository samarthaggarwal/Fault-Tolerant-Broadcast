# Fault-Tolerant-Agreement
Course Project for CS598 - Fault Tolerant Distributed Algorithms

## Dolev-Strong


## Galil-Mayer

### Central server
- while initialising processes, choose a random subset of atmost f nodes and set fail_flag=True for them

### Blackbox
Justified Uses:
- start each round by sending a message to each alive node
- checkpointing algorithm

Other Uses:
- Recruitment
- Termination

### Individual Nodes
- if fail_flag==True, fail at a random time during the run of algorithm, send a message to blackbox before failing
	- Alternative = every round runs for a (conservatively) preset time

