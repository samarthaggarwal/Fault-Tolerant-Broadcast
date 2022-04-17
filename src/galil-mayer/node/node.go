package node

import (
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
)

func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

type Node struct {
	Id         int
	NumNodes   int
	Tree       types.DiffusionTree
	NodeCh     []chan types.Msg
	BlackboxCh chan types.Msg
	OutputCh   chan types.Msg
	Value      int
	MsgCount   int
	//FailureProbability		float
}

func (n *Node) Initialise(id int, numNodes int, nodeCh []chan types.Msg, blackboxCh chan types.Msg, outputCh chan types.Msg) {
	n.Id = id
	n.NumNodes = numNodes
	n.NodeCh = nodeCh
	n.BlackboxCh = blackboxCh
	n.OutputCh = outputCh
	n.Value = -1
	n.MsgCount = 0
}

func (n *Node) Get_id() int {
	return n.Id
}

func (n *Node) Execute() {
	// Main loop listening to msg on the channel of the node
	//msg := types.Msg{
	//			Sender: n.Id,
	//			Content: fmt.Sprintf("hey main from node %d", n.Id),
	//		}
	//n.OutputCh <- msg

	tree := types.DiffusionTree{}
	phase := -1
	level := types.Leaf
	children := make([]int, 0)
	valueSent := false

	for {
		types.DPrintf("Waiting for msg node %d\n", n.Id)
		recvMsg := <-n.NodeCh[n.Id]
		if recvMsg.TypeOfMsg == types.DIFF_TREE {
			types.DPrintf("Node %d Recvied diff tree : %v\n", n.Id, recvMsg.Content.(types.DiffusionTree))
			tree = recvMsg.Content.(types.DiffusionTree)
			if n.Id == tree.Root {
				level = types.Leader
			} else {
				level = types.Coordinator
			}
			for _, node := range tree.Coordinators {
				if level == types.Leader {
					children = append(children, node.Id)
				} else if node.Id == n.Id {
					for j := max(0, node.FirstChild); j <= node.LastChild; j++ {
						children = append(children, j)
					}
					break
				}
			}
		} else if recvMsg.TypeOfMsg == types.START_PHASE {
			types.DPrintf("Node %d Recvied start phase no : %v\n", n.Id, recvMsg.Content.(int))
			phase = recvMsg.Content.(int)
			if phase == 1 {
				if !valueSent && level == types.Leader {
					sendMsg := types.Msg{
						Sender:    n.Id,
						TypeOfMsg: types.VALUE,
						Content:   n.Value,
					}
					for _, child := range children {
						n.NodeCh[child] <- sendMsg
						n.MsgCount++
					}
					valueSent = true
				}
			} else if phase == 2 {
				sendMsg := types.Msg{
					Sender:    n.Id,
					TypeOfMsg: types.CHECKPOINT,
					Content:   n.Value,
				}
				n.BlackboxCh <- sendMsg
			} else if phase == 3 {
				if !valueSent && level == types.Coordinator {
					sendMsg := types.Msg{
						Sender:    n.Id,
						TypeOfMsg: types.VALUE,
						Content:   n.Value,
					}
					for _, child := range children {
						n.NodeCh[child] <- sendMsg
						n.MsgCount++
					}
					valueSent = true
				}
			} else if phase == 4 {
				valueSent = false
				sendMsg := types.Msg{
					Sender:    n.Id,
					TypeOfMsg: types.CHECKPOINT,
					Content:   n.Value,
				}
				n.BlackboxCh <- sendMsg
			}

		} else if recvMsg.TypeOfMsg == types.VALUE {
			types.DPrintf("Node %d Recvied VALUE : %v\n", n.Id, recvMsg.Content.(int))
			n.Value = recvMsg.Content.(int)
			if phase == 1 && !valueSent && level == types.Coordinator {
				sendMsg := types.Msg{
					Sender:    n.Id,
					TypeOfMsg: types.VALUE,
					Content:   n.Value,
				}
				for _, child := range children {
					n.NodeCh[child] <- sendMsg
					n.MsgCount++
				}
				valueSent = true
			} else if phase == 3 && level != types.Leaf {
				types.DPrintf("ERROR: non-leaf received value in phase 3\n")
			}
		} else if recvMsg.TypeOfMsg == types.CHECKPOINT {
			cpMsg, ok := recvMsg.Content.(types.CPmsg)
			if ok {
				types.DPrintf("Node %d Recvied CHECKPOINT : %v\n", n.Id, recvMsg.Content.(types.CPmsg))
				n.Value = cpMsg.Value
				//E1 := cpMsg.E
			} else {
				types.DPrintf("Node %d Recvied CHECKPOINT 2 : %v\n", n.Id, recvMsg.Content.([]int))
				//E2 := msg.Content.([]int)
			}
		} else if recvMsg.TypeOfMsg == types.TERMINATE {
			types.DPrintf("Node %d Recvied TERMINATE : %v\n", n.Id, recvMsg.Content)
			sendMsg := types.Msg{
				Sender:    n.Id,
				TypeOfMsg: types.TERMINATE,
				Content:   n.Value,
			}
			n.OutputCh <- sendMsg
			break
		}
	}
}
