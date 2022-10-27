package env

import (
	pb "distributedelection/node/pb"
	//. "distributedelection/node/pkg/env"
)

func ToNetElectionMsg(sm *MsgElection) *pb.Election {
	return &pb.Election{Starter: sm.GetStarter(), Ids: sm.GetVoters()}
}

func ToNetCoordinatorMsg(sm *MsgCoordinator) *pb.Coordinator {
	return &pb.Coordinator{Starter: sm.GetStarter(), Coordinator: sm.GetCoordinator()}
}

func ToSMElectionMsg(net *pb.Election) *MsgElection {
	return &MsgElection{Starter: net.GetStarter(), Voters: net.GetIds()}
}

func ToSMCoordinatorMsg(net *pb.Coordinator) *MsgCoordinator {
	return &MsgCoordinator{Starter: net.GetStarter(), Coordinator: net.GetCoordinator()}
}

func ToNetOkMsg(sm *MsgOk) *pb.Ok {
	return &pb.Ok{Starter: sm.GetStarter()}
}
func ToSMOkMsg(net *pb.Ok) *MsgOk {
	return &MsgOk{Starter: net.GetStarter()}
}
func ToNetHeartbeatMsg(sm *MsgHeartbeat) *pb.Heartbeat {
	return &pb.Heartbeat{Id: sm.GetId()}
}

func ToSMHeartbeat(net *pb.Heartbeat) *MsgHeartbeat {
	return &MsgHeartbeat{Id: net.GetId()}
}
