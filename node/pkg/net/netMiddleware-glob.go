// netMiddleware
package net

import (
	"context"
	pbn "distributedelection/node/pb"
	pbsr "distributedelection/serviceregistry/pb"

	. "distributedelection/node/pkg/env"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"net"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
)

type DGnode struct {
	pbn.UnimplementedDistrElectNodeServer
}
type DGservreg struct {
	pbsr.UnimplementedDistrElectServRegServer
}

var cs pbsr.DistrElectServRegClient

var netList []*SMNode
var w *grpc.Server
var lis net.Listener
var serverConn *grpc.ClientConn

func InitializeNetMW() {
	// il centrale espone il servizio di identificazione dei nodi
	serverConn = ConnectToNode(State.ServRegAddr)

	listener, err := net.Listen("tcp", State.NodeInfo.GetFullAddr())
	lis = listener
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to listen to port %v:\n%v", State.NodeInfo.GetPort(), err)
	}
	// New server instance and service registering
	w = grpc.NewServer()
	pbn.RegisterDistrElectNodeServer(w, &DGnode{})
	// Defining client interface, to be used to invoke the fredricksonlynch service
	cs = pbsr.NewDistrElectServRegClient(serverConn)
	DirtyNetList = true
}

// Returns a connection with the node whose address is specified by <code>addr</code>
func ConnectToNode(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while contacting server on %v:\n %v", addr, err)
	}
	return conn
}

func Listen() {
	smlog.Info(LOG_NETWORK, "Listening on port %v.", State.NodeInfo.GetPort())
	if err := w.Serve(lis); err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to serve request: %v", err)
	}
}

func contactServiceReg() *grpc.ClientConn {
	smlog.Trace(LOG_NETWORK, "Contacting service registry")
	conn := ConnectToNode(State.ServRegAddr)
	defer conn.Close() //chiusura, se porta problemi controllare
	return conn
}
func AskForJoining() *SMNode {
	DirtyNetList = true
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	defer cancel()
	smlog.Info(LOG_SERVREG, "asking for joining the ring...")
	node, err := cs.JoinRing(ctx, &pbsr.NodeAddr{Host: State.NodeInfo.GetHost(), Port: State.NodeInfo.GetPort()})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing fredricksonlynch:\n%v", err)
	}
	return ToSMNode(node)
}

func AskForNodeInfo(i int32) *SMNode {
	smlog.Debug(LOG_SERVREG, "Asking servReg for info about node n. %d", i)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	defer cancel()
	ret, err := cs.GetNode(ctx, &pbsr.NodeId{Id: int32(i)})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing GetNode:\n%v", err)
		return nil
	}
	return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}

}

// For monitoring use only
func AskForAllNodesList() []*SMNode {
	var ret []*SMNode
	smlog.Debug(LOG_SERVREG, "Asking for info about all nodes")
	if DirtyNetList {
		smlog.Debug(LOG_SERVREG, "Election has occurred, so net could have changed. Asking to ServReg...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
		defer cancel()
		allNodesList, err := cs.GetAllNodes(ctx, new(empty.Empty))
		if err != nil {
			smlog.Fatal(LOG_NETWORK, "Error while executing GetAllNodes:\n%v", err)
		}
		for _, node := range allNodesList.List {
			ret = append(ret, ToSMNode(node))
		}
		netList = ret
	} else {
		ret = netList
	}
	DirtyNetList = false
	return ret
}

func SafeHB(hb *pbn.Heartbeat, node *SMNode) {
	connN := ConnectToNode(node.GetFullAddr())
	defer connN.Close()
	// New server instance and service registering
	nodoServer := grpc.NewServer()
	pbn.RegisterDistrElectNodeServer(nodoServer, &DGnode{})
	csN := pbn.NewDistrElectNodeClient(connN)
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

func AskForNodesWithGreaterIds(baseId int32) []*SMNode {
	smlog.Trace(LOG_SERVREG, "Chiedo al centrale informazioni sui nodi con id > %d", baseId)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	//	locCtx = ctx
	defer cancel()

	ret, err := cs.GetAllNodesWithIdGreaterThan(ctx, &pbsr.NodeId{Id: int32(baseId)})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "errore in GETNODOa:\n%v", err)
		return nil
	}
	var array []*SMNode
	for _, node := range ret.GetList() {
		array = append(array, ToSMNode(node))
	}
	return array //&SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}

}
func AskForAllNodes() []*SMNode {
	smlog.Trace(LOG_SERVREG, "Chiedo al centrale informazioni su TUTTI i nodi")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Second)
	//	locCtx = ctx
	defer cancel()

	ret, errr := cs.GetAllNodes(ctx, new(empty.Empty))
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
