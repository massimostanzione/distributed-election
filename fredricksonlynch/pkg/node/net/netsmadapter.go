package net

import (
	pb "fredricksonLynch/pb/node"
	. "fredricksonLynch/pkg/node/env"
)

func ToNetElectionMsg(sm *MsgElection) *pb.Election {
	return &pb.Election{Starter: sm.GetStarter(), Ids: sm.GetVoters()}
}

func ToNetCoordinatorMsg(sm *MsgCoordinator) *pb.Coordinator {
	return &pb.Coordinator{Starter: sm.GetStarter(), Coordinator: sm.GetCoordinator()}
}
func ToNetHeartbeatMsg(sm *MsgHeartbeat) *pb.Heartbeat {
	return &pb.Heartbeat{Id: sm.GetId()}
}

func ToSMElectionMsg(net *pb.Election) *MsgElection {
	return &MsgElection{Starter: net.GetStarter(), Voters: net.GetIds()}
}

func ToSMCoordinatorMsg(net *pb.Coordinator) *MsgCoordinator {
	return &MsgCoordinator{Starter: net.GetStarter(), Coordinator: net.GetCoordinator()}
}

func ToSMHeartbeat(net *pb.Heartbeat) *MsgHeartbeat {
	return &MsgHeartbeat{Id: net.GetId()}
}

func ToSMNode(net *pb.Node) *SMNode {
	return &SMNode{Id: net.GetId(), Host: net.GetHost(), Port: net.GetPort()}
}

func ToNetNode(sm SMNode) *pb.Node {
	return &pb.Node{Id: sm.GetId(), Host: sm.GetHost(), Port: int32(sm.GetPort())}
}
