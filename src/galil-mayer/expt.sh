#!/bin/bash

FILE=$1 #"expt.csv"

if [ -f $FILE ]; then
	echo "$FILE already exists"
	exit 1
fi

# headers
echo "numNodes,max_faults,failProb,real_faults,honestMsgCount,failedMsgCount,totalMsgCount,multiplier" > $FILE

for n in $(seq 10 10 100); do
	for f in $(seq 0 1 $n); do
		for p in $(seq 0 0.3 1); do
			for try in {0..0}; do
				echo "go run main.go $n $f $p"
				go run main.go $n $f $p > /dev/null 2>>$FILE
			done
		done
	done
done

