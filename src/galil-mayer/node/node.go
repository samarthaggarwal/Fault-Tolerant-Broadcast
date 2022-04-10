package node

import (
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
	"fmt"
)

type Node struct {
	Id				int
	NumNodes		int
	Tree			types.DiffusionTree
	NodeCh			[]chan types.Msg
	BlackboxCh		chan types.Msg
	OutputCh		chan types.Msg
	Value			int
	//FailureProbability		float
}

func (n *Node) Initialise(id int, numNodes int, nodeCh []chan types.Msg, blackboxCh chan types.Msg, outputCh chan types.Msg) {
	n.Id = id
	n.NumNodes = numNodes
	n.NodeCh = nodeCh
	n.BlackboxCh = blackboxCh
	n.OutputCh = outputCh
	n.Value = -1
}

func (n *Node) Get_id() int {
	return n.Id
}

func (n *Node) Execute () {
	// Main loop listening to msg on the channel of the node
	msg := types.Msg{
				Sender: n.Id,
				Content: fmt.Sprintf("hey main from node %d", n.Id),
			}
	n.OutputCh <- msg
}



