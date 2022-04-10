package blackbox

import (
	"Fault-Tolerant-Agreement/src/galil-mayer/types"
	"time"
)

type Blackbox struct {
	// God process
	Id			int
	NumNodes	int
	MyCh		chan types.Msg
	NodeCh		[]chan types.Msg
	OutputCh	chan types.Msg
}

//func (b *Blackbox) ReinitialiseDiffusionTree (startIndex int) {
//	// called with startIndex=0 at the start of round 1, and then with lowest alive leaf when E2={}
//	n := b.NumNodes - startIndex
//}

func (b *Blackbox) Execute () {
	// Main loop listening to msg on blackbox channel

	time.Sleep(2 * time.Second)
	msg := types.Msg{
				Sender: b.Id,
				Content: "hey main from god",
			}
	b.OutputCh <- msg
}


