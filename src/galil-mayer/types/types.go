package types

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
	type1 = iota
	// TODO
)

type Msg struct {
	Sender			int
	TypeOfMsg		MsgType
	Content			string
}


