//TODO ev. in nuovo file che è net ma specializzato grpc?
package net

import (
	"context"
	. "fredricksonLynch/pkg/node/env"

	//"fredricksonLynch/pkg/node"
	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "fredricksonLynch/pb/node"
)

func (s *DGnode) ForwardElection(ctx context.Context, in *pb.Election) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_SENT, "\033[43m\033[1;30mRECEIVED ELECTION %v\033[0m", in)
	if !Pause {
		ElectionChannel <- ToSMElectionMsg(in)
	}
	return NONE, status.New(codes.OK, "").Err()
}

func (s *DGnode) ForwardCoordinator(ctx context.Context, in *pb.Coordinator) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_RECV, "\033[43m\033[1;30mRECEIVED COORDINATOR %v\033[0m", in)
	if !Pause {
		CoordChannel <- ToSMCoordinatorMsg(in)
	}
	return NONE, status.New(codes.OK, "").Err()
}
func (s *DGnode) SendHeartBeat(ctx context.Context, in *pb.Heartbeat) (*pb.NONE, error) {

	smlog.Info(LOG_HB, "Ricevo HEARTBEAT dal nodo %d", in.GetId())
	if !Pause {
		Heartbeat <- ToSMHeartbeat(in)
	}
	return NONE, status.New(codes.OK, "").Err()
}

func ToSMElectionMsg(net *pb.Election) *MsgElection {
	return &MsgElection{Starter: net.GetStarter(), Voters: net.GetIds()}
}

func ToNetElectionMsg(sm *MsgElection) *pb.Election {
	return &pb.Election{Starter: sm.GetStarter(), Ids: sm.GetVoters()}
}

func ToNetCoordinatorMsg(sm *MsgCoordinator) *pb.Coordinator {
	return &pb.Coordinator{Starter: sm.GetStarter(), Coordinator: sm.GetCoordinator()}
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
