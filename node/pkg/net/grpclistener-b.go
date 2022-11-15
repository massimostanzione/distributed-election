package net

import (
	"context"
	. "distributedelection/node/pkg/env"
	. "distributedelection/tools/formatting"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	pb "distributedelection/node/pb"

	// following import is replaced with EMPTY_NODE message,
	// ref. https://github.com/massimostanzione/distributed-election/issues/88
	// empty "github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DEANode) ForwardElectionBully(ctx context.Context, in *pb.ElectionBully) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"ELECTION %v"+ColorReset, in)
	ElectionChannel_bully <- ToSMElectionBullyMsg(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) ForwardElectionFL(ctx context.Context, in *pb.ElectionFL) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"ELECTION %v"+ColorReset, in)
	ElectionChannel_fl <- ToSMElectionFLMsg(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) ForwardOk(ctx context.Context, in *pb.Ok) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"OK %v\033[0m"+ColorReset, in)
	OkChannel <- ToSMOkMsg(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) ForwardCoordinator(ctx context.Context, in *pb.Coordinator) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"COORDINATOR %v\033[0m"+ColorReset, in)
	CoordChannel <- ToSMCoordinatorMsg(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) SendHeartBeat(ctx context.Context, in *pb.Heartbeat) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, "HB from node %d", in.GetId())
	Heartbeat <- ToSMHeartbeat(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}
