# Fault-Tolerant-Agreement
Course Project for CS598 - Fault Tolerant Distributed Algorithms

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

## Dolev-Strong

### Message Complexity
	Honest = (n-f)*(n-1)
	Faulty = f*(n-1)
	Total = n*(n-1)

## TODO
1. Generate comparison graphs

2. Message complexity analytical formula

3. Report

4. PPT

## Graphs
1. For f = [0, 10, n/2, ->n, n-1], plot of ratio of msg count of galil mayer vs dolev strong.

2. For n = [10, 50, 100, 500, 1000], plot of msg count of galil mayer
	T = total msg count
	






