package net

import (
	"context"
	. "distributedelection/node/pkg/env"
	. "distributedelection/tools/formatting"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "distributedelection/node/pb"
)

func (s *DGnode) ForwardElection(ctx context.Context, in *pb.Election) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"RECEIVED ELECTION %v"+ColorReset, in)
	ElectionChannel <- ToSMElectionMsg(in)
	return NONE, status.New(codes.OK, "").Err()
}

func (s *DGnode) ForwardOk(ctx context.Context, in *pb.Ok) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"RECEIVED OK %v\033[0m"+ColorReset, in)
	OkChannel <- ToSMOkMsg(in)
	return NONE, status.New(codes.OK, "").Err()
}

func (s *DGnode) ForwardCoordinator(ctx context.Context, in *pb.Coordinator) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"RECEIVED COORDINATOR %v\033[0m"+ColorReset, in)
	CoordChannel <- ToSMCoordinatorMsg(in)
	return NONE, status.New(codes.OK, "").Err()
}

func (s *DGnode) SendHeartBeat(ctx context.Context, in *pb.Heartbeat) (*pb.NONE, error) {
	smlog.Info(LOG_MSG_RECV, "Received HB from node %d", in.GetId())
	Heartbeat <- ToSMHeartbeat(in)
	return NONE, status.New(codes.OK, "").Err()
}
