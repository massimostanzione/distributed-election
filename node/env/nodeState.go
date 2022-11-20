// state.go
package env

import . "distributedelection/tools/api"

type NodeState struct {
	NodeInfo    *SMNode
	Coordinator int32
	Participant bool

	// limit servReg requests if network is not changed,
	// i.e. if no election has occurred
	DirtyNetCache bool
}

var CurState *NodeState = &NodeState{}

// Node knowledge, for ring-based algorithm(s)
var NextNode *SMNode = &SMNode{}

var NetCache []*SMNode
