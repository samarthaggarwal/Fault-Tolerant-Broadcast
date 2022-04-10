package main

import (
	"fmt"
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
	"Fault-Tolerant-Agreement/src/galil-mayer/node"
	"Fault-Tolerant-Agreement/src/galil-mayer/blackbox"
)

func main() {
	//fmt.Println("hello world!")
	//n := node.Node{Id:1234}
	//fmt.Printf("Id = %d\n", n.Get_id())

	// Parameter Selection
	numNodes := 8
	//faults := 0
	bufferSize := 256 // so that sending to channel is unblocking
	godId := -1

	// Initialise all channels
	nodeCh := make([]chan types.Msg, numNodes)
	for i:=0; i<numNodes; i++ {
		nodeCh[i] = make(chan types.Msg, bufferSize)
	}
	blackboxCh := make(chan types.Msg, bufferSize)
	outputCh := make(chan types.Msg, bufferSize)

	// Initialise all nodes
	nodes := make([]node.Node, numNodes)
	for i:=0; i<numNodes; i++ {
		nodes[i].Initialise(i, numNodes, nodeCh, blackboxCh, outputCh)
		fmt.Printf("Initialised node with id=%d \n", nodes[i].Get_id())
	}
	nodes[0].Value = 1234 // TODO - change to random

	// Initialise blackbox
	god := blackbox.Blackbox{
					Id: godId,
					NumNodes: numNodes,
					MyCh: blackboxCh,
					NodeCh: nodeCh,
					OutputCh: outputCh,
				}
	fmt.Printf("God is watching over %d nodes \n", god.NumNodes)

	for i:=0; i<numNodes; i++ {
		go nodes[i].Execute()
	}
	go god.Execute()

	// Wait for outputs
	for {
		msg := <-outputCh
		if msg.Sender == godId {
			fmt.Printf("msg.Sender = %d, msg.Content = %s \n", msg.Sender, msg.Content)
			break
		} else {
			fmt.Printf("msg.Sender = %d, msg.Content = %s \n", msg.Sender, msg.Content)
		}
	}

	fmt.Println("Exiting main")
}




