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
mkdir -p $LOGDIR

# range
range_n=($(seq 5 5 45) $(seq 50 50 450) $(seq 500 100 900) $(seq 1000 1000 10000))
range_p=($(seq 0 0.1 0.99))
num_trials=20
logidx=0
logfile=${LOGDIR}/${logidx}.log

date
echo "range_n=${range_n[@]}"
echo "range_p=${range_p[@]}"
echo "num_trials=${num_trials}"

# headers
echo "numNodes,maxFaults,failProb,realFaults,honestMsgCount,failedMsgCount,totalMsgCount,multiplier" > $CSVFILE

for n in ${range_n[@]}; do
	#range_f=$(seq 0 1 $n)
	range_f=(0 1 2 `let "y=n/4"; echo $y` `let "y=n/2"; echo $y` `let "y=3*n/4"; echo $y` `let "y=n-1"; echo $y`)
	range_f=($(echo ${range_f[@]} | tr ' ' '\n' | sort | uniq | tr '\n' ' '))
	for f in ${range_f[@]}; do
		for p in ${range_p[@]}; do
			for try in $(seq ${num_trials}); do
				#echo "go run src/main.go $n $f $p"
				echo "n=${n},f=${f},p=${p},try=${try}"
				echo "#### n=${n},f=${f},p=${p},try=${try} ####" > $logfile
				go run src/main.go $n $f $p >> $logfile 2>>$CSVFILE

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

date

