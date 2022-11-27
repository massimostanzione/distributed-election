// Handlers for gRPC remote method invocations (RMI).
package net

import (
	"context"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	pb "distributedelection/serviceregistry/pb"

	reg "distributedelection/serviceregistry/structure/behavior"
	api "distributedelection/tools/api"

	// following import is replaced with EMPTY_NODE message,
	// ref. https://github.com/massimostanzione/distributed-election/issues/88
	// empty "github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DEAServer struct {
	pb.UnimplementedDistrElectServRegServer
}

func (s *DEAServer) JoinNetwork(ctx context.Context, in *pb.NodeAddr) (*pb.Node, error) {
	smlog.Info(LOG_MSG_RECV, "JoinNetwork [%v]", in)
	return api.ToNetNode(reg.RegisterNewNode(in.GetHost(), in.GetPort())), status.New(codes.OK, "").Err()
}

func (s *DEAServer) GetNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodesWithIdGreaterThan [%v]", in)
	return api.ToNetNode(reg.FetchRecordById(int(in.GetId()))), status.New(codes.OK, "").Err()
}

func (s *DEAServer) GetAllNodes(ctx context.Context, in *pb.EMPTY_SR) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodes")
	return api.ToNetNodeList(reg.GetNodesWithBaseId(0)), status.New(codes.OK, "").Err()
}

func (s *DEAServer) GetAllNodesWithIdGreaterThan(ctx context.Context, in *pb.NodeId) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodesWithIdGreaterThan [%v]", in)
	return api.ToNetNodeList(reg.GetNodesWithBaseId(int32(in.GetId()))), status.New(codes.OK, "").Err()
}
