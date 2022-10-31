package net

import (
	"context"
	. "distributedelection/serviceregistry/pkg/registrymgt"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	pb "distributedelection/serviceregistry/pb"
	. "distributedelection/tools/api"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DGserver struct {
	pb.UnimplementedDistrElectServRegServer
}

func (s *DGserver) JoinNetwork(ctx context.Context, in *pb.NodeAddr) (*pb.Node, error) {
	smlog.Info(LOG_MSG_RECV, "JoinNetwork [%v]", in)
	node := ManageJoining(in.GetHost(), in.GetPort())
	return ToNetNode(node), status.New(codes.OK, "").Err()
}

func (s *DGserver) GetNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodesWithIdGreaterThan [%v]", in)
	node := FetchRecordbyId(int(in.Id))
	return ToNetNode(node), status.New(codes.OK, "").Err()
}

func (s *DGserver) GetAllNodes(ctx context.Context, in *empty.Empty) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodes")
	return GetAllNodesExecutive(0), status.New(codes.OK, "").Err()
}

func (s *DGserver) GetAllNodesWithIdGreaterThan(ctx context.Context, in *pb.NodeId) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "GetAllNodesWithIdGreaterThan [%v]", in)
	return GetAllNodesExecutive(in.GetId()), status.New(codes.OK, "").Err()
}
