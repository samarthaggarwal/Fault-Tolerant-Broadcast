#!/bin/bash

CSVFILE=$1 #"expt.csv"
LOGDIR=$2

if [ $# -ne 2 ]; then
	echo "Run Instructions: $0 <csvfile> <logdir>"
	exit 1
fi

if [ -f $CSVFILE ]; then
	echo "$CSVFILE already exists"
	exit 1
fi

# range
range_n=($(seq 10 50 40))
range_p=($(seq 0 0.9 1))
num_trials=2
logidx=0
logfile=${LOGDIR}/${logidx}.log

# headers
echo "numNodes,maxFaults,failProb,realFaults,honestMsgCount,failedMsgCount,totalMsgCount,multiplier" > $CSVFILE

for n in ${range_n[@]}; do
	for f in $(seq 0 1 $n); do
		for p in ${range_p[@]}; do
			for try in $(seq ${num_trials}); do
				#echo "go run main.go $n $f $p"
				echo "#### n=${n},f=${f},p=${p},try=${try} ####" > $logfile
				go run main.go $n $f $p >> $logfile 2>>$CSVFILE

				# check if log has error
				grep -iq ERROR $logfile
				if [[ $? -eq 0 ]]; then
					let "logidx=logidx+1"
					logfile=${LOGDIR}/${logidx}.log
				fi
			done
		done
	done
done

