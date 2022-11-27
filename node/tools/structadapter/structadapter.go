// Conversion tools between gRPC structures and internal (SM) ones.
package structadapter

import (
	. "distributedelection/node/env"
	pb "distributedelection/node/pb"
)

// From internal to gRPC -------------------------------------------------------

func ToNetElectionBullyMsg(sm *MsgElectionBully) *pb.ElectionBully {
	return &pb.ElectionBully{Starter: sm.GetStarter()}
}

func ToNetElectionFLMsg(sm *MsgElectionFL) *pb.ElectionFL {
	return &pb.ElectionFL{Starter: sm.GetStarter(), Ids: sm.GetVoters()}
}

func ToNetCoordinatorMsg(sm *MsgCoordinator) *pb.Coordinator {
	return &pb.Coordinator{Starter: sm.GetStarter(), Coordinator: sm.GetCoordinator()}
}

func ToNetOkMsg(sm *MsgOk) *pb.Ok {
	return &pb.Ok{Starter: sm.GetStarter()}
}

func ToNetHeartbeatMsg(sm *MsgHeartbeat) *pb.Heartbeat {
	return &pb.Heartbeat{Id: sm.GetId()}
}

// From gRPC to internal -------------------------------------------------------

func ToSMElectionBullyMsg(net *pb.ElectionBully) *MsgElectionBully {
	return &MsgElectionBully{Starter: net.GetStarter()}
}

func ToSMElectionFLMsg(net *pb.ElectionFL) *MsgElectionFL {
	return &MsgElectionFL{Starter: net.GetStarter(), Voters: net.GetIds()}
}

func ToSMCoordinatorMsg(net *pb.Coordinator) *MsgCoordinator {
	return &MsgCoordinator{Starter: net.GetStarter(), Coordinator: net.GetCoordinator()}
}

func ToSMOkMsg(net *pb.Ok) *MsgOk {
	return &MsgOk{Starter: net.GetStarter()}
}

func ToSMHeartbeat(net *pb.Heartbeat) *MsgHeartbeat {
	return &MsgHeartbeat{Id: net.GetId()}
}
