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
