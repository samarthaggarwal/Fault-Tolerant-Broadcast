package main

import (
	"Fault-Tolerant-Agreement/src/galil-mayer/blackbox"
	"Fault-Tolerant-Agreement/src/galil-mayer/node"
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
	"fmt"
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
	//fmt.Println("hello world!")
	//n := node.Node{Id:1234}
	//fmt.Printf("Id = %d\n", n.Get_id())

	// Parameter Selection
	numNodes := 8
	//faults := 0
	bufferSize := 10000 // so that sending to channel is unblocking
	godId := -1

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
		nodes[i].Initialise(i, numNodes, nodeCh, blackboxCh, outputCh)
		fmt.Printf("Initialised node with id=%d \n", nodes[i].Get_id())
	}
	secret := 1234 // TODO - change to random
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
		values[i] = -1
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

	// Sanity Checks - TODO
	//success := true
	//for (i:=0; i<numNodes; i++) {
	//	if contains(b.DeadNodes, i) ||
	//		values[i] ==
	//}
	fmt.Printf("deadNodes:%v\n values:%v lenvalues: %d\n", god.DeadNodes, values, len(values))

	fmt.Println("Exiting main")
}
