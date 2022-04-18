package types

import "log"

const debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if debug > 0 {
		log.Printf(format, a...)
	}
	return
}

const (
	SLEEPTIME = 2000
)

type Level int

const (
	Leader Level = iota
	Coordinator
	Leaf
)

type CoordinatorNode struct {
	Id         int
	FirstChild int
	LastChild  int
}

type DiffusionTree struct {
	Root         int
	Coordinators []CoordinatorNode
}

type MsgType int

const (
	DIFF_TREE MsgType = iota
	START_PHASE
	VALUE
	CHECKPOINT
	TERMINATE
	FAILURE
	DONE
)

type Msg struct {
	Sender    int
	TypeOfMsg MsgType
	Content   interface{}
}

type CPmsg struct {
	Value int
	E     []int
}
