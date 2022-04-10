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
}

func (b *Blackbox) ReinitialiseDiffusionTree (startIndex int) types.DiffusionTree {
	// called with startIndex=0 at the start of round 1, and then with lowest alive leaf when E2={}
	n := b.NumNodes - startIndex
	tree := types.DiffusionTree{}

	if n<3 {
		return tree, errors.New("cannot create DiffusionTree with <3 nodes")
	}

	tree.Root = startIndex
	numCoordinators := math.Ceil(math.Sqrt(n)) - 1

	// TODO - complete this

	return tree
}

func (b *Blackbox) Terminate () {
	// TODO
}

func (b *Blackbox) RecomputeTree () {
	// Rearrange the leaves among the coordinators and return the tree
	// TODO
}

func (b *Blackbox) Execute () {
	// Main loop listening to msg on blackbox channel

	b.CurrentTree = b.ReinitialiseDiffusionTree(0)

	for {
		L0 = b.CurrentTree.Root
		L1 = b.CurrentTree.Coordinators

		// send CurrentTree to all L0+L1

		// send roundStart msg to all at the start of each round

		// p1
		// p2, we have E1
		// p3
		// p4, we have E2

		/* if L1 is_subset_of E2 {
			b.Terminate()
		} else if E2 = {} {
			b.CurrentTree = b.ReinitialiseDiffusionTree(L1[0].FirstChild)
		} else {
			b.RecomputeTree(E2)
		}
		*/
	}

	time.Sleep(2 * time.Second)
	msg := types.Msg{
				Sender: b.Id,
				Content: "hey main from god",
			}
	b.OutputCh <- msg
}


