// Messages
package env

type MsgType uint8

const (
	MSG_UNDEFINED MsgType = iota
	MSG_ELECTION
	MSG_OK
	MSG_COORDINATOR
	MSG_HEARTBEAT
)

type MsgElection struct {
	Starter int32
}

func (msg *MsgElection) GetStarter() int32 {
	return msg.Starter
}
func NewElectionMsg(Starter int32) *MsgElection {
	return &MsgElection{Starter: Starter}
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

type MsgCoordinator struct {
	Coordinator int32
}

func (msg *MsgCoordinator) GetCoordinator() int32 {
	return msg.Coordinator
}

func NewCoordinatorMsg(elected int32) *MsgCoordinator {
	return &MsgCoordinator{Coordinator: elected}
}

type MsgHeartbeat struct {
	Id int32
}

func (msg *MsgHeartbeat) GetId() int32 {
	return msg.Id
}
