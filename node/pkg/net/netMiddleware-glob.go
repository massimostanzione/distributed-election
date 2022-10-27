// netMiddleware
package net

//var cs pb.DistGrepClient
import (
	"context"
	pb "distributedelection/node/pb"

	. "distributedelection/node/pkg/env"
	//"distributedelection/node/pkg/statemachine"
	//	. "distributedelection/tools/formatting"
	//"distributedelection/node/pkg"

	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"net"
	"time"

	"google.golang.org/grpc"
)

type DGnode struct {
	pb.UnimplementedDistGrepServer
}

var NONE = &pb.NONE{}
var cs pb.DistGrepClient

var w *grpc.Server
var lis net.Listener
var serverConn *grpc.ClientConn //server

func InitializeNetMW() {
	// il centrale espone il servizio di identificazione dei nodi
	conn := ConnectToNode(ServRegAddr)
	serverConn = conn

	liss, err := net.Listen("tcp", Me.GetFullAddr())
	lis = liss
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to listen to port %v:\n%v", Me.GetPort(), err)
	}
	// New server instance and service registering
	w = grpc.NewServer()
	pb.RegisterDistGrepServer(w, &DGnode{})
	// Serve incoming calls

	// Defining client interface, to be used to invoke the fredricksonlynch service
	cs = pb.NewDistGrepClient(serverConn)
}

func ConnectToNode(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while contacting server on %v:\n %v", addr, err)
	}
	return conn
}

func Listen() {
	smlog.Info(LOG_NETWORK, "Listening on port %v.", Me.GetPort())
	if err := w.Serve(lis); err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to serve request: %v", err)
	}
	//	for pause {
	//	}
}

func contactServiceReg() *grpc.ClientConn {
	smlog.Trace(LOG_NETWORK, "Contacting service registry")
	conn := ConnectToNode(ServRegAddr)
	defer conn.Close() //chiusura, se porta problemi controllare
	return conn
}
func AskForJoining() *SMNode {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	smlog.Info(LOG_SERVREG, "asking for joining the ring...")
	node, err := cs.JoinRing(ctx, &pb.NodeAddr{Host: Me.GetHost(), Port: Me.GetPort()})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing fredricksonlynch:\n%v", err)
	}
	return ToSMNode(node)
}

func AskForNodeInfo(i int32) *SMNode {
	smlog.Debug(LOG_SERVREG, "Asking servReg for info about node n. %d", i)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ret, err := cs.GetNode(ctx, &pb.NodeId{Id: int32(i)})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing GetNode:\n%v", err)
		return nil
	}
	return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}

}

// For monitoring use only
func AskForAllNodesList() []*SMNode {
	smlog.Debug(LOG_SERVREG, "Asking servReg for info about all nodes")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	allNodesList, err := cs.GetAllNodes(ctx, NONE)
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing GetAllNodes:\n%v", err)
	}
	var ret []*SMNode
	for _, node := range allNodesList.List {
		ret = append(ret, ToSMNode(node))
	}
	return ret
}
func SafeHB(hb *pb.Heartbeat, node *SMNode) {
	connN := ConnectToNode(node.GetFullAddr())
	//defer connN.Close()
	// New server instance and service registering
	nodoServer := grpc.NewServer()
	pb.RegisterDistGrepServer(nodoServer, &DGnode{})
	csN := pb.NewDistGrepClient(connN)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)
	//	locCtx = ctx
	defer cancel()
	//netMsg := ToNetHeartbeatMsg(hb)
	_, errq := csN.SendHeartBeat(ctx, hb)
	//_, errq := csN.InoltraElezione(ctx, elezione)
	if errq != nil {
		smlog.Error(LOG_NETWORK, "error while contacting %v", node.GetFullAddr())
		smlog.Debug(LOG_NETWORK, "(%s)", errq)
	}
}

func AskForNodesWithGreaterIds(baseId int32, forceRunningNode bool) []*SMNode {
	smlog.Trace(LOG_SERVREG, "Chiedo al centrale informazioni sui nodi con id > %d", baseId)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()
	if forceRunningNode {
		//de facto not implemented
		_, errr := cs.GetNextRunningNode(ctx, &pb.NodeId{Id: int32(baseId)})
		if errr != nil {
			smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
			return nil
		}
		return nil //&SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	} else {
		ret, errr := cs.GetAllNodesWithIdGreaterThan(ctx, &pb.NodeId{Id: int32(baseId)})
		if errr != nil {
			smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
			return nil
		}
		var array []*SMNode
		for _, node := range ret.GetList() {
			array = append(array, ToSMNode(node))
		}
		return array //&SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	}

}
func AskForAllNodes() []*SMNode {
	smlog.Trace(LOG_SERVREG, "Chiedo al centrale informazioni su TUTTI i nodi")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()

	ret, errr := cs.GetAllNodes(ctx, NONE)
	if errr != nil {
		smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
		return nil
	}
	//conversion
	//TODO implementare anche nelle altre chiamate simili
	var array []*SMNode
	for _, node := range ret.GetList() {
		array = append(array, ToSMNode(node))
	}
	return array //&SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
}
