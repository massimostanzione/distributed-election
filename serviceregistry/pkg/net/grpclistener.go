// Handlers for gRPC service calls
// Processing is performed in netMiddleware.go
package net

import (
	"context"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	pb "distributedelection/serviceregistry/pb"

	//empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DEAServer struct {
	pb.UnimplementedDistrElectServRegServer
}

func (s *DEAServer) JoinNetwork(ctx context.Context, in *pb.NodeAddr) (*pb.Node, error) {
	smlog.Info(LOG_MSG_RECV, "JoinNetwork [%v]", in)
	return ManageJoining(in.GetHost(), in.GetPort()), status.New(codes.OK, "").Err()
}

func (s *DEAServer) GetNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodesWithIdGreaterThan [%v]", in)
	return FetchRecordById(in.GetId()), status.New(codes.OK, "").Err()
}

func (s *DEAServer) GetAllNodes(ctx context.Context, in *pb.EMPTY_SR) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodes")
	return GetAllNodesExecutive(0), status.New(codes.OK, "").Err()
}

func (s *DEAServer) GetAllNodesWithIdGreaterThan(ctx context.Context, in *pb.NodeId) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodesWithIdGreaterThan [%v]", in)
	return GetAllNodesExecutive(in.GetId()), status.New(codes.OK, "").Err()
}
