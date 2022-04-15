package types

const (
	SLEEPTIME = 100
)

type Level int
const (
	Leader = iota
	Coordinator
	Leaf
)

type CoordinatorNode struct {
	Id			int
	FirstChild	int
	LastChild	int
}

type DiffusionTree struct {
	Root			int
	Coordinators	[]CoordinatorNode
}

type MsgType int
const (
	DIFF_TREE = iota
	START_PHASE
	VALUE
	CHECKPOINT
	TERMINATE
	FAILURE
)

type Msg struct {
	Sender			int
	TypeOfMsg		MsgType
	Content			interface{}
}

type CPmsg struct {
	Value	int
	E		[]int
}

