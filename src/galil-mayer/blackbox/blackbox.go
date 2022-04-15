package blackbox

import (
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
	"time"
	"errors"
	"math"
)

type Blackbox struct {
	// God process
	Id			int
	NumNodes	int
	MyCh		chan types.Msg
	NodeCh		[]chan types.Msg
	OutputCh	chan types.Msg
	CurrentTree	types.DiffusionTree
	DeadNodes	[]int
}

func contains(list []int, elem int) bool {
	for _, v := range list {
		if v == elem {
			return true
		}
	}
	return false
}

func min(a int, b int) {
	if a<=b {
		return a
	} else {
		return b
	}
}

func max(a int, b int) {
	if a>=b {
		return a
	} else {
		return b
	}
}

func (b *Blackbox) ReinitialiseDiffusionTree (startIndex int) error {
	// called with startIndex=0 at the start of round 1, and then with lowest alive leaf when E2={}
	n := b.NumNodes - startIndex
	tree := types.DiffusionTree{}

	if n<3 {
		return errors.New("cannot create DiffusionTree with <3 nodes")
	}

	tree.Root = startIndex
	numCoordinators := math.Ceil(math.Sqrt(n)) - 1
	numLeaves := n - 1 - numCoordinators
	quotient := numLeaves / numCoordinators
	remainder := numLeaves % numCoordinators
	firstLeaf := startIndex + numCoordinators + 1

	for (i:=1; i<=numCoordinators; i++) {
		numChildren := quotient
		if i <= remainder { numChildren++ }
		coordinator := types.Coordinator{
							Id: startIndex+i,
							FirstChild: firstLeaf,
							LastChild: firstLeaf + numChildren - 1
						}
		firstLeaf += numChildren
		tree.Coordinators = append(tree.Coordinators, coordinator)
	}

	fmt.Printf("Reinitialised Tree with startIndex=%d\n", startIndex)

	//return tree
	b.CurrentTree = tree
	return nil
}

func (b *Blackbox) Terminate () {

	for recvMsg, ok := <-b.BlackboxCh {
		if !ok { break }
		if recvMsg.MsgType == types.FAILURE {
			b.DeadNodes = append(b.DeadNodes, recvMsg.Sender)
		} else {
			fmt.Printf("ERROR: Received non-failure msg in Terminate. msg=%v\n", recvMsg)
		}
	}
	msg = types.Message{
				Sender: -1,
				TypeOfMsg: types.MsgType.TERMINATE
			}
	for (i:=0; i<b.NumNodes; i++) {
		b.NodeCh[i] <- msg
	}
}

func (b *Blackbox) RecomputeTree (E2 []int) error {
	// Rearrange the leaves among the coordinators and return the tree
	// assumption - E2 is sorted

	if len(E2)==0 {
		return errors.New("cannot call recompute tree when E2={}")
	}

	tree = types.DiffusionTree{}
	tree.Root = E2[0]

	numDeadL1 := len(b.CurrentTree.Coordinators) - len(E2)
	if b.CurrentTree.Root == E2[0] {
		numDeadL1++
	}
	numAliveL1 := len(E2) - 1

	for _, coordinator := range b.CurrentTree.Coordinators {
		if coordinator.FirstChild==-1 && !contains(E2, coordinator.Id) {
			numDeadL1--
		}
	}
	// numDeadL1: number of dead parents whose children need to hear the value

	numL1 := len(b.CurrentTree.Coordinators)
	firstRecruit := b.CurrentTree.Coordinators[numL1-1].Id + 1
	currRecruit := firstRecruit

	numRecruits := numDeadL1 - numAliveL1
	if numRecruits < 0 {numRecruits = 0}

	indexInE2 := 1
	i := 0
	for i < len(b.CurrentTree.Coordinators) {
		if contains(E2, b.CurrentTree.Coordinators[i].Id) ||
		b.CurrentTree.Coordinators[i].LastChild < firstRecruit + numRecruits ||
		b.CurrentTree.Coordinator[i].LastChild == -1 {
			// either this coordinator does not have children or its children have already heard the value
			i++
			continue
		}
		coordinator := types.Coordinator{
							Id: -1,
							FirstChild: b.CurrentTree.Coordinators[i].FirstChild,
							LastChild: b.CurrentTree.Coordinators[i].LastChild
						}
		coordinator.FirstChild = max(coordinator.FirstChild, firstRecruit + numRecruits)
		if indexInE2 < len(E2) {
			// reassign children of existing coordinator
			coordinator.Id = E2[indexInE2]
			indexInE2++
		} else {
			// Recruit leaves as coordinators
			coordinator.Id = currRecruit
			currRecruit++
		}
		tree.Coordinators = append(tree.Coordinators, coordinator)
		i++
	}
	for indexInE2 < len(E2) {
		// leave extra alive coordinators idle without any children
		coordinator := types.Coordinator{Id:E2[indexInE2], FirstChild:-1, LastChild:-1}
		tree.Coordinators = append(tree.Coordinators, coordinator)
		indexInE2++
	}

	for currRecruit < firstRecruit + numRecruits {
		// recruit extra coordinators without any children
		coordinator := types.Coordinator{Id:E2[indexInE2], FirstChild:-1, LastChild:-1}
		tree.Coordinators = append(tree.Coordinators, coordinator)
		currRecruit++
	}

	b.CurrentTree = tree
	return nil
}

func (b *Blackbox) SendPhaseStart (phase_id int, sendOnlyToCoordinators bool) {
	msg = types.Msg{
				Sender: -1,
				TypeOfMsg: types.MsgType.START_PHASE,
				Content: phase_id
			}
	lastNode := b.NumNodes-1
	if sendOnlyToCoordinators {
		lastNode = b.Coordinators[len(b.Coordinators)-1].Id
	}
	for (i:=0, i<=lastNode; i++) {
		b.NodeCh[i] <- msg
	}
}

func (b *Blackbox) Execute () {
	// Main loop listening to msg on blackbox channel

	b.ReinitialiseDiffusionTree(0)
	b.DeadNodes := make([]int)

	for {
		L0 = b.CurrentTree.Root
		L1 = b.CurrentTree.Coordinators

		// send CurrentTree to all L0+L1
		L0UL1 := make([]int)
		L0UL1 = append(L0UL1, L0)
		for _, node := range L1 {
			L0UL1 = append(L0UL1, node.Id)
		}
		msg := types.Msg{
					Sender: -1,
					TypeOfMsg: types.MsgType.DIFF_TREE,
					Content: b.CurrentTree
				}
		for _, node := range L0UL1 {
			b.NodeCh[node] <- msg
		}

		// send roundStart msg to all at the start of each round

		// p1
		b.SendPhaseStart(1, false)
		time.Sleep(types.SLEEPTIME * time.Millisecond)

		// p2, we have E1
		b.SendPhaseStart(2, true)
		time.Sleep(types.SLEEPTIME * time.Millisecond)
		value := -1
		E1 := make([]int)
		for recvMsg, ok := <-b.BlackboxCh {
			if !ok { break }
			if recvMsg.MsgType == types.FAILURE {
				b.DeadNodes = append(b.DeadNodes, recvMsg.Sender)
			} else if recvMsg.MsgType == types.CHECKPOINT {
				E1 = append(E1, recvMsg.Sender)
				recv_value := recvMsg.Content.(int)
				value = max(value, recv_value)
			}
		}
		msg = types.Message{
					Sender: -1,
					TypeOfMsg: types.MsgType.CHECKPOINT,
					Content: types.CPmsg{value: value, E: E1}
				}
		for _, node := range L0UL1 {
			b.NodeCh[node] <- msg
		}

		// p3
		b.SendPhaseStart(3, false)
		time.Sleep(types.SLEEPTIME * time.Millisecond)

		// p4, we have E2 (assumption: E2 is sorted)
		b.SendPhaseStart(4, true)
		time.Sleep(types.SLEEPTIME * time.Millisecond)
		E2 := make([]int)
		for recvMsg, ok := <-b.BlackboxCh {
			if !ok { break }
			if recvMsg.MsgType == types.FAILURE {
				b.DeadNodes = append(b.DeadNodes, recvMsg.Sender)
			} else if recvMsg.MsgType == types.CHECKPOINT {
				E2 = append(E2, recvMsg.Sender)
			}
		}
		msg = types.Message{
					Sender: -1,
					TypeOfMsg: types.MsgType.CHECKPOINT,
					Content: E2
				}
		for _, node := range L0UL1 {
			b.NodeCh[node] <- msg
		}

		// check for termination, else start p5
		terminate := true
		for i, coordinator := range L1 {
			if coordinator.FirstChild!=-1 && !contains(E2, coordinator.Id) {
				terminate = false
				break
			}
		}

		if terminate {
			// wait for nodes failures if any
			time.Sleep(types.SLEEPTIME * time.Millisecond)
			b.Terminate()
			break
		} else if len(E2)==0 {
			b.ReinitialiseDiffusionTree(L1[len(L1)-1] + 1)
		} else {
			b.RecomputeTree(E2)
		}
	}

	time.Sleep(2 * time.Second)
	msg := types.Msg{
				Sender: b.Id,
				TypeOfMsg: types.MsgType.TERMINATE,
			}
	b.OutputCh <- msg
}


