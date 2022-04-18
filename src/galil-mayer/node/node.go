package node

import (
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
	"math/rand"
	"fmt"
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
	FP			float64
}

func (n *Node) Initialise(id int, numNodes int, nodeCh []chan types.Msg, blackboxCh chan types.Msg, outputCh chan types.Msg, fp float64) {
	n.Id = id
	n.NumNodes = numNodes
	n.NodeCh = nodeCh
	n.BlackboxCh = blackboxCh
	n.OutputCh = outputCh
	n.Value = -1
	n.MsgCount = 0
	n.FP = fp
}

func (n *Node) Get_id() int {
	return n.Id
}

func (n *Node) try_failure(location string, content interface{}) bool {
	if rand.Float64() < n.FP {
		sendMsg := types.Msg{
			Sender:    n.Id,
			TypeOfMsg: types.FAILURE,
			Content: content,
		}
		n.BlackboxCh <- sendMsg
		types.DPrintf("Node:%d, failed at %s\n", n.Id, location)
		return true
	}
	return false
}

func (n *Node) Execute() {

	tree := types.DiffusionTree{}
	phase := -1
	level := types.Leaf
	children := make([]int, 0)
	valueSent := false
	killed := false

	for !killed {
		types.DPrintf("Node:%d, waiting for msg\n", n.Id)
		recvMsg := <-n.NodeCh[n.Id]
		if recvMsg.TypeOfMsg == types.DIFF_TREE {
			types.DPrintf("Node:%d, received diff tree: %v\n", n.Id, recvMsg.Content.(types.DiffusionTree))
			tree = recvMsg.Content.(types.DiffusionTree)
			if n.Id == tree.Root {
				level = types.Leader
			} else {
				level = types.Coordinator
			}
			children = make([]int, 0)
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
			doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
			n.BlackboxCh <- doneMsg
		} else if recvMsg.TypeOfMsg == types.START_PHASE {
			types.DPrintf("Node:%d, received start phase no: %v\n", n.Id, recvMsg.Content.(int))
			phase = recvMsg.Content.(int)
			if phase == 1 {
				if valueSent { panic("valueSent=true at phase1") }
				if !valueSent && level == types.Leader {
					sendMsg := types.Msg{
						Sender:    n.Id,
						TypeOfMsg: types.VALUE,
						Content:   n.Value,
					}
					for i, child := range children {
						if i==len(children)/2 {
							children_sent := make([]int, 0)
							for j:=0; j<i; j++ {
								children_sent = append(children_sent, children[j])
							}
							killed = n.try_failure("Leader/phase=1", children_sent)
							if killed { break }
						}
						n.NodeCh[child] <- sendMsg
						n.MsgCount++
						types.DPrintf("Node:%d, sent value to node:%d, phase1/leader", n.Id, child)
					}
					if killed { break }
					valueSent = true
					doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
					n.BlackboxCh <- doneMsg
				}
			} else if phase == 2 {
				sendMsg := types.Msg{
					Sender:    n.Id,
					TypeOfMsg: types.CHECKPOINT,
					Content:   n.Value,
				}
				killed = n.try_failure("Coordinator/phase=2", make([]int, 0))
				if killed { break }
				n.BlackboxCh <- sendMsg
				doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
				n.BlackboxCh <- doneMsg
			} else if phase == 3 {
				if !valueSent && level == types.Coordinator {
					sendMsg := types.Msg{
						Sender:    n.Id,
						TypeOfMsg: types.VALUE,
						Content:   n.Value,
					}
					for i, child := range children {
						if i==len(children)/2 {
							killed = n.try_failure("Coordinator/phase=3", make([]int, 0))
							if killed { break }
						}
						n.NodeCh[child] <- sendMsg
						n.MsgCount++
						types.DPrintf("Node:%d, sent value to node:%d, phase3/coordinator", n.Id, child)
					}
					if killed { break }
					valueSent = true
				}
				if level == types.Coordinator {
					doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
					n.BlackboxCh <- doneMsg
				}
			} else if phase == 4 {
				valueSent = false
				sendMsg := types.Msg{
					Sender:    n.Id,
					TypeOfMsg: types.CHECKPOINT,
					Content:   n.Value,
				}
				killed = n.try_failure("Coordinator/phase=4", make([]int, 0))
				if killed { break }
				n.BlackboxCh <- sendMsg
				doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
				n.BlackboxCh <- doneMsg
			}

		} else if recvMsg.TypeOfMsg == types.VALUE {
			types.DPrintf("Node:%d, received VALUE : %v\n", n.Id, recvMsg.Content.(int))
			n.Value = recvMsg.Content.(int)
			if phase == 1 && !valueSent && level == types.Coordinator {
				sendMsg := types.Msg{
					Sender:    n.Id,
					TypeOfMsg: types.VALUE,
					Content:   n.Value,
				}
				for i, child := range children {
					if i==len(children)/2 {
						killed = n.try_failure("Coordinator/phase=1", make([]int, 0))
						if killed { break }
					}
					n.NodeCh[child] <- sendMsg
					n.MsgCount++
					types.DPrintf("Node:%d, sent value to node:%d, phase1/coordinator", n.Id, child)
				}
				if killed { break }
				valueSent = true
				doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
				n.BlackboxCh <- doneMsg
			} else if phase == 3 && level != types.Leaf {
				fmt.Printf("ERROR: non-leaf %d/%v received value from %d in phase 3\n", n.Id, level, recvMsg.Sender)
			} else if level == types.Leaf {
				killed = n.try_failure(fmt.Sprintf("Leaf/phase%d", phase), make([]int, 0))
				if killed { break }
			} else if phase == 1 && valueSent && level == types.Coordinator {
				panic("phase1, level=Coordinator, valueSent=True")
			}
		} else if recvMsg.TypeOfMsg == types.CHECKPOINT {
			cpMsg, ok := recvMsg.Content.(types.CPmsg)
			if ok {
				types.DPrintf("Node:%d, received CHECKPOINT : %v\n", n.Id, recvMsg.Content.(types.CPmsg))
				n.Value = cpMsg.Value
				//E1 := cpMsg.E
			} else {
				types.DPrintf("Node:%d, received CHECKPOINT 2 : %v\n", n.Id, recvMsg.Content.([]int))
				//E2 := msg.Content.([]int)
			}
			doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
			n.BlackboxCh <- doneMsg
		} else if recvMsg.TypeOfMsg == types.TERMINATE {
			types.DPrintf("Node:%d, received TERMINATE : %v\n", n.Id, recvMsg.Content)
			sendMsg := types.Msg{
				Sender:    n.Id,
				TypeOfMsg: types.TERMINATE,
				Content:   n.Value,
			}
			n.OutputCh <- sendMsg
			doneMsg := types.Msg{Sender: n.Id, TypeOfMsg: types.DONE}
			n.BlackboxCh <- doneMsg
			break
		}
	}
}
