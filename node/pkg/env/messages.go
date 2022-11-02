// Messages
package env

type MsgType uint8

const (
	MSG_UNDEFINED MsgType = iota
	MSG_ELECTION
	MSG_COORDINATOR
	MSG_HEARTBEAT
	MSG_OK
)

type MsgElection struct {
	Starter int32
	Voters  []int32
}

func (msg *MsgElection) GetStarter() int32 {
	return msg.Starter
}
func (msg *MsgElection) GetVoters() []int32 {
	return msg.Voters
}
func (msg *MsgElection) AddVoter(newVoter int32) *MsgElection {
	msg.Voters = append(msg.Voters, newVoter)
	return msg
}
func NewElectionMsg() *MsgElection {
	return &MsgElection{Starter: CurState.NodeInfo.GetId(), Voters: []int32{CurState.NodeInfo.GetId()}}
}

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

type MsgOk struct {
	Starter int32
}

func (msg *MsgOk) GetStarter() int32 {
	return msg.Starter
}
func NewOkMsg(Starter int32) *MsgOk {
	return &MsgOk{Starter: Starter}
}

type MsgHeartbeat struct {
	Id int32
}

func (msg *MsgHeartbeat) GetId() int32 {
	return msg.Id
}
