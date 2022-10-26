// netMiddleware
package net

//var cs pb.DistGrepClient
import (
	"context"
	pb "fredricksonlynch/node/pb"

	. "fredricksonlynch/node/pkg/env"
	//"fredricksonlynch/node/pkg/statemachine"
	//	. "fredricksonlynch/tools/formatting"
	//"fredricksonlynch/node/pkg"

	. "fredricksonlynch/tools/smlog"
	smlog "fredricksonlynch/tools/smlog"
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

const RMI_RETRY_TOLERANCE = 1
const LATE_HB_TOLERANCE = 3

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

func StopServer() {
	smlog.Info(LOG_NETWORK, "Stopping server")
	w.Stop()
	//	for !pause {
	//	}
}
func SwitchServerState(run bool) {
	if !run {
		Listen()
	} else {
		StopServer()
	}
}
func contactServiceReg() *grpc.ClientConn {
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

//func askForNodeInfo(i int32, forceRunningNode bool) (int32, string) {
func AskForNodeInfo(i int32) *SMNode {
	smlog.Info(LOG_SERVREG, "Chiedo al centrale informazioni sul nodo %d", i)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()

	ret, errr := cs.GetNode(ctx, &pb.NodeId{Id: int32(i)})
	if errr != nil {
		smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
		return nil
	}
	return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}

}
func AskForAllNodesList() []*SMNode {
	smlog.Debug(LOG_SERVREG, "vado a chiedere tutti i nodi")
	var errl error = error(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	allNodesList, errl := cs.GetAllNodes(ctx, NONE)
	if errl != nil {
		smlog.Fatal(LOG_SERVREG, "problema in GetAllNodes: %v", errl)
	}
	var ret []*SMNode
	for _, node := range allNodesList.List {
		ret = append(ret, ToSMNode(node))

	}
	return ret
}