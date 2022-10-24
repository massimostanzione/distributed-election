package net

import (
	"context"
	. "fredricksonlynch/node/pkg/env"
	. "fredricksonlynch/tools/formatting"

	//"fredricksonlynch/node/pkg"
	. "fredricksonlynch/tools/smlog"
	smlog "fredricksonlynch/tools/smlog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "fredricksonlynch/node/pb"
)

func (s *DGnode) ForwardElection(ctx context.Context, in *pb.Election) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_SENT, ColorBlkBckgrYellow+BoldBlack+"RECEIVED ELECTION %v"+ColorReset, in)
	if !Pause {
		ElectionChannel <- ToSMElectionMsg(in)
	}
	return NONE, status.New(codes.OK, "").Err()
}

func (s *DGnode) ForwardCoordinator(ctx context.Context, in *pb.Coordinator) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"RECEIVED COORDINATOR %v\033[0m"+ColorReset, in)
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
