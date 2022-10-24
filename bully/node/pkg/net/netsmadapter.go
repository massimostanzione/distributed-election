package net

import (
	pb "bully/node/pb"
	. "bully/node/pkg/env"
)

func ToNetElectionMsg(sm *MsgElection) *pb.Election {
	return &pb.Election{Starter: sm.GetStarter()}
}
func ToNetOkMsg(sm *MsgOk) *pb.Ok {
	return &pb.Ok{Starter: sm.GetStarter()}
}

func ToNetCoordinatorMsg(sm *MsgCoordinator) *pb.Coordinator {
	return &pb.Coordinator{Coordinator: sm.GetCoordinator()}
}

func ToSMElectionMsg(net *pb.Election) *MsgElection {
	return &MsgElection{Starter: net.GetStarter()}
}

func ToSMOkMsg(net *pb.Ok) *MsgOk {
	return &MsgOk{Starter: net.GetStarter()}
}
func ToSMCoordinatorMsg(net *pb.Coordinator) *MsgCoordinator {
	return &MsgCoordinator{Coordinator: net.GetCoordinator()}
}

/*
func ToSMHeartbeat(net *pb.Heartbeat) *MsgHeartbeat {
	return &MsgHeartbeat{Id: net.GetId()}
}
*/
func ToSMNode(net *pb.Node) *SMNode {
	return &SMNode{Id: net.GetId(), Host: net.GetHost(), Port: net.GetPort()}
}

func ToNetNode(sm SMNode) *pb.Node {
	return &pb.Node{Id: sm.GetId(), Host: sm.GetHost(), Port: int32(sm.GetPort())}
}

func ToNetHeartbeatMsg(sm *MsgHeartbeat) *pb.Heartbeat {
	return &pb.Heartbeat{Id: sm.GetId()}
}
func ToSMHeartbeat(net *pb.Heartbeat) *MsgHeartbeat {
	return &MsgHeartbeat{Id: net.GetId()}
}
