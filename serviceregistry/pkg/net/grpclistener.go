package net

import (
	"context"
	//	. "distributedelection/serviceregistry/pkg/env"
	. "distributedelection/serviceregistry/pkg/registrymgt"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"fmt"

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
	smlog.Info(LOG_MSG_RECV, "Request for joining, from address %s", in.GetHost()+":"+fmt.Sprint(in.GetPort()))
	node := ManageJoining(in.GetHost(), in.GetPort())
	return ToNetNode(node), status.New(codes.OK, "").Err()
	//return &pb.Node{Id: int32(nodeRecord.Id), Host: nodeRecord.Host, Port: nodeRecord.Port}, status.New(codes.OK, "").Err()
}

func (s *DGserver) GetAllNodes(ctx context.Context, in *empty.Empty) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "Request to get all nodes")
	return GetAllNodesExecutive(0), status.New(codes.OK, "").Err()
}

func (s *DGserver) GetAllNodesWithIdGreaterThan(ctx context.Context, in *pb.NodeId) (*pb.NodeList, error) {
	smlog.Info(LOG_MSG_RECV, "Request to get all nodes with ID greater than %d", in.GetId())
	return GetAllNodesExecutive(in.GetId()), status.New(codes.OK, "").Err()
}
func (s *DGserver) GetNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	smlog.Info(LOG_MSG_RECV, "Request to get node with id = %d", in.GetId())
	node := FetchRecordbyId(int(in.Id))
	return &pb.Node{Id: int32(node.Id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}
