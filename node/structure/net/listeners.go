package net

import (
	"context"
	. "distributedelection/node/env"
	pb "distributedelection/node/pb"
	sa "distributedelection/node/tools/structadapter"
	. "distributedelection/tools/formatting"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	// following import is replaced with EMPTY_NODE message,
	// ref. https://github.com/massimostanzione/distributed-election/issues/88
	// empty "github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DEANode) ForwardElectionBully(ctx context.Context, in *pb.ElectionBully) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"ELECT %v"+ColorReset, in)
	ElectionChannel_bully <- sa.ToSMElectionBullyMsg(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) ForwardElectionFL(ctx context.Context, in *pb.ElectionFL) (*pb.EMPTY_NODE, error) {
	MsgOrderIn <- MSG_ELECTION_FL
	ElectChIn <- sa.ToSMElectionFLMsg(in)
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"ELECT %v"+ColorReset, in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) ForwardOk(ctx context.Context, in *pb.Ok) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"OK %v\033[0m"+ColorReset, in)
	OkChannel <- sa.ToSMOkMsg(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) ForwardCoordinator(ctx context.Context, in *pb.Coordinator) (*pb.EMPTY_NODE, error) {
	if Cfg.ALGORITHM == DE_ALGORITHM_FREDRICKSONLYNCH {
		MsgOrderIn <- MSG_COORDINATOR
	}
	CoordChIn <- sa.ToSMCoordinatorMsg(in)
	smlog.Info(LOG_MSG_RECV, ColorBlkBckgrYellow+BoldBlack+"COORD %v\033[0m"+ColorReset, in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}

func (s *DEANode) SendHeartBeat(ctx context.Context, in *pb.Heartbeat) (*pb.EMPTY_NODE, error) {
	smlog.Info(LOG_MSG_RECV, "HB from node %d", in.GetId())
	Heartbeat <- sa.ToSMHeartbeat(in)
	return new(pb.EMPTY_NODE), status.New(codes.OK, "").Err()
}
