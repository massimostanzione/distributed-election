// Messages
package env

type msgType uint8

//TODO applicare a safeRMI
const (
	MSG_UNDEFINED msgType = iota
	MSG_ELECTION
	MSG_COORDINATOR
	MSG_HEARTBEAT
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
func (msg *MsgElection) AddVoter(newVoter int32) {
	msg.Voters = append(msg.Voters, newVoter)
}
func NewElectionMsg(Starter int32) *MsgElection {
	return &MsgElection{Starter: Starter, Voters: []int32{Me.GetId()}}
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

type MsgHeartbeat struct {
	Id int32
}

func (msg *MsgHeartbeat) GetId() int32 {
	return msg.Id
}
