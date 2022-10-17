package net

import (
	"context"
	"fmt"
	. "fredricksonLynch/pkg/node/env"
	. "fredricksonLynch/tools/formatting"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "fredricksonLynch/pb/node"
)

type DGserver struct {
	pb.UnimplementedDistGrepServer
}

func (s *DGserver) JoinRing(ctx context.Context, in *pb.NodeAddr) (*pb.Node, error) {
	log.Printf("*** REQUEST RECEIVED ***")
	log.Printf("Un nodo vuole entrare, all'indirizzo %s", in.GetHost()+":"+fmt.Sprint(in.GetPort()))

	node, existent := fetchRecordbyAddr(in.GetHost(), in.GetPort())
	if existent {
		log.Printf("Il nodo all'indirizzo %s risulta registrato, con id=%d", in.GetHost()+":"+fmt.Sprint(in.GetPort()), node.id)
	} else {
		nodes = append(nodes, node)
		log.Printf("Il nodo mi è nuovo, gli assegno id=%d", node.id)
	}
	PrintRing()
	//compreso map della struct locale in quella grpc-compatibile:
	return &pb.Node{Id: int32(node.id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}

func (s *DGserver) ReportAsFailed(ctx context.Context, in *pb.Node) (*pb.NONE, error) {
	log.Printf("Il nodo n. %d, presso %s, mi è stato segnalato come failed", in.GetHost()+":"+fmt.Sprint(in.GetPort()), in.GetHost()+":"+fmt.Sprint(in.GetPort()))
	for i := range nodes {
		if int32(nodes[i].id) == in.GetId() {
			nodes[i].reportedAsFailed = true
			PrintRing()
			return NONE, status.New(codes.OK, "").Err()
		}
	}
	log.Fatalf("Cannot find node %d, to be flagged as FAILED", in.GetId())
	return NONE, status.New(codes.OK, "").Err()
}

// TODO documentazione: può succedere che un nodo veda un altro temporaneamente out,
// quindi lo segnala come tale, e al ricevere di un messaggio da tale nodo,
// non risponde in quanto il server gli dice ancora che è out, quindi bisogna aggiornare
func (s *DGserver) ReportAsRunning(ctx context.Context, in *pb.Node) (*pb.NONE, error) {
	log.Printf("Il nodo n. %d, presso %s, mi è stato segnalato come RUNNING", in.GetId(), in.GetHost()+":"+fmt.Sprint(in.GetPort()))
	for i := range nodes {
		if int32(nodes[i].id) == in.GetId() {
			nodes[i].reportedAsFailed = false
			printRing()
			return NONE, status.New(codes.OK, "").Err()
		}
	}
	log.Fatalf("Cannot find node %d, to be flagged as RUNNING", in.GetId())
	return NONE, status.New(codes.OK, "").Err()
}

func (s *DGserver) GetAllNodes(ctx context.Context, in *pb.NONE) (*pb.NodeList, error) {
	return getAllNodesExecutive(false), status.New(codes.OK, "").Err()
}
func (s *DGserver) GetAllRunningNodes(ctx context.Context, in *pb.NONE) (*pb.NodeList, error) {
	return getAllNodesExecutive(true), status.New(codes.OK, "").Err()
}

func (s *DGserver) GetNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	log.Printf("*** REQUEST RECEIVED ***")
	log.Printf("Serve conoscere chi è %d", in.Id)
	node, _ := fetchRecordbyId(int(in.Id), false)
	log.Printf("Ti ritorno il nodo richiesto, che è %s", node.Host+":"+fmt.Sprint(node.Port))
	printRing()
	return &pb.Node{Id: int32(node.id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}

func (s *DGserver) GetNextRunningNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	log.Printf("*** REQUEST RECEIVED ***")
	log.Printf("Serve conoscere chi è %d", in.Id)
	node, _ := fetchRecordbyId(int(in.Id), true)
	// TODO controllo su err
	/*if err != nil {
		log.Fatalf("gestire")
		return NONE, false
	}*/
	log.Printf("Ti ritorno il nodo richiesto, che è %s", node.Host+":"+fmt.Sprint(node.Port))
	printRing()
	return &pb.Node{Id: int32(node.id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}
