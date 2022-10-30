package net

import (
	"context"
	. "distributedelection/node/pkg/env"
	. "distributedelection/tools/formatting"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	pb "distributedelection/node/pb"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DGnode) ForwardElection(ctx context.Context, in *pb.Election) (*empty.Empty, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"RECEIVED ELECTION %v"+ColorReset, in)
	ElectionChannel <- ToSMElectionMsg(in)
	return new(empty.Empty), status.New(codes.OK, "").Err()
}

func (s *DGnode) ForwardOk(ctx context.Context, in *pb.Ok) (*empty.Empty, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"RECEIVED OK %v\033[0m"+ColorReset, in)
	OkChannel <- ToSMOkMsg(in)
	return new(empty.Empty), status.New(codes.OK, "").Err()
}

func (s *DGnode) ForwardCoordinator(ctx context.Context, in *pb.Coordinator) (*empty.Empty, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"RECEIVED COORDINATOR %v\033[0m"+ColorReset, in)
	CoordChannel <- ToSMCoordinatorMsg(in)
	return new(empty.Empty), status.New(codes.OK, "").Err()
}

func (s *DGnode) SendHeartBeat(ctx context.Context, in *pb.Heartbeat) (*empty.Empty, error) {
	smlog.Info(LOG_MSG_RECV, "Received HB from node %d", in.GetId())
	Heartbeat <- ToSMHeartbeat(in)
	return new(empty.Empty), status.New(codes.OK, "").Err()
}
