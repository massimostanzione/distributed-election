// Internal message structures and related functions.
// Conversion from/to gRPC structures are done in netsmadapter.
package env

type MsgType uint8

const (
	MSG_UNDEFINED MsgType = iota
	MSG_ELECTION_BULLY
	MSG_ELECTION_FL
	MSG_COORDINATOR
	MSG_HEARTBEAT
	MSG_OK
)

//-------------------------------------------------------------
// ELECTION (bully-specific)

type MsgElectionBully struct {
	Starter int32
	Voters  []int32
}

func (msg *MsgElectionBully) GetStarter() int32 {
	return msg.Starter
}

func NewElectionBullyMsg() *MsgElectionBully {
	return &MsgElectionBully{Starter: CurState.NodeInfo.GetId()}
}

//-------------------------------------------------------------
// ELECTION (FL-specific)

type MsgElectionFL struct {
	Starter int32
	Voters  []int32
}

func (msg *MsgElectionFL) GetStarter() int32 {
	return msg.Starter
}

func (msg *MsgElectionFL) GetVoters() []int32 {
	return msg.Voters
}

func (msg *MsgElectionFL) AddVoter(newVoter int32) *MsgElectionFL {
	msg.Voters = append(msg.Voters, newVoter)
	return msg
}

func NewElectionFLMsg() *MsgElectionFL {
	return &MsgElectionFL{Starter: CurState.NodeInfo.GetId(), Voters: []int32{CurState.NodeInfo.GetId()}}
}

//-------------------------------------------------------------
// COORDINATOR (common to bully and FL)

type MsgCoordinator struct {
	Starter     int32
	Coordinator int32
}

func (msg *MsgCoordinator) GetStarter() int32 {
	return msg.Starter
}
func (msg *MsgCoordinator) GetCoordinator() int32 {
	return msg.Coordinator
}

func NewCoordinatorMsg(Starter int32, elected int32) *MsgCoordinator {
	return &MsgCoordinator{Starter: Starter, Coordinator: elected}
}

//-------------------------------------------------------------
// OK (bully-specific)

type MsgOk struct {
	Starter int32
}

func (msg *MsgOk) GetStarter() int32 {
	return msg.Starter
}
func NewOkMsg(Starter int32) *MsgOk {
	return &MsgOk{Starter: Starter}
}

//-------------------------------------------------------------
// HEARTBEAT (for monitoing purpose only)

type MsgHeartbeat struct {
	Id int32
}

func (msg *MsgHeartbeat) GetId() int32 {
	return msg.Id
}
