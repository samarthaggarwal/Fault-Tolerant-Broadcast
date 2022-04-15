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
	MsgCount		int
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

func (n *Node) Execute () {
	// Main loop listening to msg on the channel of the node
	//msg := types.Msg{
	//			Sender: n.Id,
	//			Content: fmt.Sprintf("hey main from node %d", n.Id),
	//		}
	//n.OutputCh <- msg

	tree := types.DiffusionTree{}
	phase := -1
	level := types.Level.Leaf
	children := make([]int)
	valueSent := false

	for {
		msg := <-n.NodeCh[n.Id]
		if msg.TypeOfMsg == DIFF_TREE {
			tree = msg.Content.(types.DiffusionTree)
			if n.Id == tree.Root {
				level = types.Level.Leader
			} else {
				level = types.Level.Coordinator
			}
			for _, node := range tree.Coordinators {
				if level == types.Level.Leader {
					children = append(children, node.Id)
				} else if node.Id == n.Id {
					for (j:=max(0,node.FirstChild); j<=node.LastChild; j++) {
						children = append(children, j)
					}
					break
				}
			}
		} else if msg.TypeOfMsg == START_PHASE {
			phase = msg.Content.(int)
			if phase == 1 {
				if !valueSent && level == types.Level.Leader {
					sendMsg := types.Msg{
									Sender: n.Id,
									TypeOfMsg: types.MsgType.VALUE,
									Content: n.Value
								}
					for _, child := range children {
						n.NodeCh[child] <- msg
						n.MsgCount++
					}
					valueSent = true
				}
			} else if phase == 2 {
				msg = types.Msg{
							Sender: n.Id,
							TypeOfMsg: types.MsgType.CHECKPOINT,
							Content: n.Value
						}
				n.BlackboxCh <- msg
			} else if phase == 3 {
				if !valueSent && level == types.Level.Coordinator {
					sendMsg := types.Msg{
									Sender: n.Id,
									TypeOfMsg: types.MsgType.VALUE,
									Content: n.Value
								}
					for _, child := range children {
						n.NodeCh[child] <- msg
						n.MsgCount++
					}
					valueSent = true
				}
			} else if phase == 4 {
				valueSent = false
				msg = types.Msg{
							Sender: n.Id,
							TypeOfMsg: types.MsgType.CHECKPOINT,
							Content: n.Value
						}
				n.BlackboxCh <- msg
			}

		} else if msg.TypeOfMsg == VALUE {
			n.Value = msg.Content.(int)
			if phase==1 && !valueSent && level == types.Level.Coordinator {
				sendMsg := types.Msg{
								Sender: n.Id,
								TypeOfMsg: types.MsgType.VALUE,
								Content: n.Value
							}
				for _, child := range children {
					n.NodeCh[child] <- msg
					n.MsgCount++
				}
				valueSent = true
			} else if phase == 3 && level != types.Level.Leaf {
				fmt.Printf("ERROR: non-leaf received value in phase 3\n")
			}
		} else if msg.TypeOfMsg == CHECKPOINT {
			cpMsg, ok := msg.Content.(types.CPmsg)
			if ok {
				n.Value = cpMsg.Value
				//E1 := cpMsg.E
			} else {
				//E2 := msg.Content.([]int)
			}
		} else if msg.TypeOfMsg == TERMINATE {
			sendMsg := types.Msg{
							Sender: n.Id,
							TypeOfMsg: types.MsgType.TERMINATE,
							Content: n.Value
						}
			n.OutputCh <- sendMsg
			break
		}
	}
}



