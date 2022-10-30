package net

import (
	"context"
	. "distributedelection/serviceregistry/pkg/env"
	. "distributedelection/serviceregistry/pkg/registrymgt"
	. "distributedelection/tools/smlog"
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

func (s *DGserver) ReportAsFailed(ctx context.Context, in *pb.Node) (*empty.Empty, error) {
	smlog.InfoU("Il nodo n. %d, presso %s, mi è stato segnalato come failed", in.GetHost()+":"+fmt.Sprint(in.GetPort()), in.GetHost()+":"+fmt.Sprint(in.GetPort()))
	for i := range Nodes {
		if int32(Nodes[i].Id) == in.GetId() {
			Nodes[i].ReportedAsFailed = true
			PrintRing()
			return new(empty.Empty), status.New(codes.OK, "").Err()
		}
	}
	smlog.Fatal(LOG_UNDEFINED, "Cannot find node %d, to be flagged as FAILED", in.GetId())
	return new(empty.Empty), status.New(codes.OK, "").Err()
}

// TODO documentazione: può succedere che un nodo veda un altro temporaneamente out,
// quindi lo segnala come tale, e al ricevere di un messaggio da tale nodo,
// non risponde in quanto il server gli dice ancora che è out, quindi bisogna aggiornare
func (s *DGserver) ReportAsRunning(ctx context.Context, in *pb.Node) (*empty.Empty, error) {
	smlog.InfoU("Il nodo n. %d, presso %s, mi è stato segnalato come RUNNING", in.GetId(), in.GetHost()+":"+fmt.Sprint(in.GetPort()))
	for i := range Nodes {
		if int32(Nodes[i].Id) == in.GetId() {
			Nodes[i].ReportedAsFailed = false
			PrintRing()
			return new(empty.Empty), status.New(codes.OK, "").Err()
		}
	}
	smlog.Fatal(LOG_UNDEFINED, "Cannot find node %d, to be flagged as RUNNING", in.GetId())
	return new(empty.Empty), status.New(codes.OK, "").Err()
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

func (s *DGserver) GetNextRunningNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	smlog.InfoU("*** REQUEST RECEIVED ***")
	smlog.InfoU("Serve conoscere chi è %d", in.Id)
	node, _ := FetchRecordbyId(int(in.Id), true)
	// TODO controllo su err
	/*if err != nil {
		smlog.Fatal(LOG_UNKNOWN,"gestire")
		return new(empty.Empty), false
	}*/
	smlog.InfoU("Ti ritorno il nodo richiesto, che è %s", node.Host+":"+fmt.Sprint(node.Port))
	PrintRing()
	return &pb.Node{Id: int32(node.Id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}
