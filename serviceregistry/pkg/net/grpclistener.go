package net

import (
	"context"
	. "distributedelection/serviceregistry/pkg/env"
	. "distributedelection/serviceregistry/pkg/registrymgt"
	smlog "distributedelection/tools/smlog"
	"fmt"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "distributedelection/serviceregistry/pb"
)

type DGserver struct {
	pb.UnimplementedDistrElectServRegServer
}

func (s *DGserver) JoinRing(ctx context.Context, in *pb.NodeAddr) (*pb.Node, error) {
	smlog.InfoU("*** REQUEST RECEIVED ***")
	smlog.InfoU("Un nodo vuole entrare, all'indirizzo %s", in.GetHost()+":"+fmt.Sprint(in.GetPort()))

	node, existent := FetchRecordbyAddr(in.GetHost(), in.GetPort())
	if existent {
		smlog.InfoU("Il nodo all'indirizzo %s risulta registrato, con id=%d", in.GetHost()+":"+fmt.Sprint(in.GetPort()), node.Id)
	} else {
		Nodes = append(Nodes, node)
		smlog.InfoU("Il nodo mi è nuovo, gli assegno id=%d", node.Id)
	}
	PrintRing()
	//compreso map della struct locale in quella grpc-compatibile:
	return &pb.Node{Id: int32(node.Id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}

func (s *DGserver) GetAllNodes(ctx context.Context, in *empty.Empty) (*pb.NodeList, error) {
	return GetAllNodesExecutive(0), status.New(codes.OK, "").Err()
}

func (s *DGserver) GetAllNodesWithIdGreaterThan(ctx context.Context, in *pb.NodeId) (*pb.NodeList, error) {
	return GetAllNodesExecutive(in.GetId()), status.New(codes.OK, "").Err()
}
func (s *DGserver) GetNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	smlog.InfoU("*** REQUEST RECEIVED ***")
	smlog.InfoU("Serve conoscere chi è %d", in.Id)
	node, _ := FetchRecordbyId(int(in.Id), false)
	smlog.InfoU("Ti ritorno il nodo richiesto, che è %s", node.Host+":"+fmt.Sprint(node.Port))
	PrintRing()
	return &pb.Node{Id: int32(node.Id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}
