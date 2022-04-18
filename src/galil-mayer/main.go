package main

import (
	"os"
	"Fault-Tolerant-Agreement/src/galil-mayer/blackbox"
	"Fault-Tolerant-Agreement/src/galil-mayer/node"
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
	"fmt"
	"math"
	"time"
	"math/rand"
	"sort"
	"strconv"
)

func contains(list []int, elem int) bool {
	for _, v := range list {
		if v == elem {
			return true
		}
	}
	return false
}

func main() {

	args := os.Args[1:]

	// Parameter Selection
	numNodes,_ := strconv.Atoi(args[0]) //500
	faults,_ := strconv.Atoi(args[1]) //numNodes - 1 //int(float64(numNodes) * 0.3)
	failProb,_ := strconv.ParseFloat(args[2], 64) //0.1
	bufferSize := 20000 // so that sending to channel is unblocking
	godId := -1

	// Failure nodes
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(numNodes)
	if rand.Float64() < 0.5 {
		sort.Ints(perm)
	}
	failureProbability := make([]float64, numNodes)
	for i:=0; i<numNodes; i++ {
		if i<faults {
			failureProbability[perm[i]] = failProb
		} else {
			failureProbability[perm[i]] = -0.1
		}
	}

	// Initialise all channels
	nodeCh := make([]chan types.Msg, numNodes)
	for i := 0; i < numNodes; i++ {
		nodeCh[i] = make(chan types.Msg, bufferSize)
	}
	blackboxCh := make(chan types.Msg, bufferSize)
	outputCh := make(chan types.Msg, bufferSize)

	// Initialise all nodes
	nodes := make([]node.Node, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i].Initialise(i, numNodes, nodeCh, blackboxCh, outputCh, failureProbability[i])
		types.DPrintf("Initialised node with id=%d \n", nodes[i].Get_id())
	}
	secret := 1234
	nodes[0].Value = secret

	// Initialise blackbox
	god := blackbox.Blackbox{
		Id:       godId,
		NumNodes: numNodes,
		MyCh:     blackboxCh,
		NodeCh:   nodeCh,
		OutputCh: outputCh,
	}
	fmt.Printf("God is watching over %d nodes \n", god.NumNodes)

	for i := 0; i < numNodes; i++ {
		go nodes[i].Execute()
	}
	go god.Execute()

	// Wait for outputs
	values := make([]int, numNodes)
	for i := range values {
		values[i] = -2 // -2: unset, -1: default value
	}
	for {
		msg := <-outputCh
		if msg.Sender == godId {
			//fmt.Printf("msg.Sender = %d, msg.Content = %s \n", msg.Sender, msg.Content)
			break
		} else {
			values[msg.Sender] = msg.Content.(int)
			//fmt.Printf("msg.Sender = %d, msg.Content = %s \n", msg.Sender, msg.Content)
		}
	}

	honestMsgCount := 0
	failedMsgCount := 0
	for i:=0; i<numNodes; i++ {
		if !contains(god.DeadNodes, i) {
			honestMsgCount += nodes[i].MsgCount
		} else {
			failedMsgCount += nodes[i].MsgCount
		}
	}
	totalMsgCount := honestMsgCount + failedMsgCount
	multiplier := float64(totalMsgCount) / (float64(numNodes) + float64(len(god.DeadNodes))*math.Sqrt(float64(numNodes)))

	// Sanity Checks
	// Liveness
	liveness := true
	for i:=0; i<numNodes; i++ {
		if !contains(god.DeadNodes, i) &&
			values[i] == -2 {
			liveness = false
			break
		}
	}
	// Safety
	safety := true
	value := -2
	for i:=0; i<numNodes; i++ {
		if contains(god.DeadNodes, i) { continue }
		if value==-2 {
			value = values[i]
		}
		if values[i] != value {
			safety = false
			break
		}
	}
	// Validity
	validity := contains(god.DeadNodes, 0) || value==secret

	sort.Ints(god.DeadNodes)

	msgCount := make([]int, numNodes)
	for i:=0; i<numNodes; i++ {
		msgCount[i] = nodes[i].MsgCount
	}

	//fmt.Printf("deadNodes:%v\n values:%v lenvalues: %d\n", god.DeadNodes, values, len(values))
	//fmt.Printf("deadNodes:%v\nlen(deadNodes):%d\nmsgCount:%v\n", god.DeadNodes, len(god.DeadNodes), msgCount)
	fmt.Printf("len(deadNodes):%d\n", len(god.DeadNodes))
	fmt.Printf("=== Sanity:%v, Liveness:%v, Safety:%v, Validity:%v, totalMsg:%d, multiplier:%f, value=%d ===\n", liveness&&safety&&validity, liveness, safety, validity, totalMsgCount, multiplier, value)

	if !liveness || !safety || !validity {
		//panic("sanity FAILED")
		fmt.Printf("ERROR: sanity failed\n")
	}

	fmt.Fprintf(os.Stderr, "%d,%d,%f,%d,%d,%d,%d,%f\n", numNodes, faults, failProb, len(god.DeadNodes), honestMsgCount, failedMsgCount, totalMsgCount, multiplier)
}
