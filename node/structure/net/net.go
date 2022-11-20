// net.go
package net

import (
	. "distributedelection/node/env"
	pbn "distributedelection/node/pb"
	pbsr "distributedelection/serviceregistry/pb"
	ip "distributedelection/tools/ip"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

var cs pbsr.DistrElectServRegClient

var nodeServerIntf *grpc.Server
var lis net.Listener
var serverConn *grpc.ClientConn

func ListenToIncomingRMI() {
	smlog.Info(LOG_NETWORK, "Listening at %s:%v", ip.GetOutboundIP(), CurState.NodeInfo.GetPort())
	if err := nodeServerIntf.Serve(lis); err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to serve request: %v", err)
	}
}

// Returns a connection with the node whose address is specified by <code>addr</code>
func ConnectToNode(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while contacting server on %v:\n %v", addr, err)
	}
	return conn
}

func InitializeNetMW() {
	// Build up listener
	// (listening is started in 'behavior')
	listener, err := net.Listen("tcp", CurState.NodeInfo.GetFullAddr())
	lis = listener
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to listen to port %v:\n%v", CurState.NodeInfo.GetPort(), err)
	}

	// Register Node server interface
	nodeServerIntf = grpc.NewServer()
	pbn.RegisterDistrElectNodeServer(nodeServerIntf, &DEANode{})

	// Establish connection with the service registry
	serverConn = ConnectToNode(Cfg.SERVREG_HOST + ":" + strconv.FormatInt(Cfg.SERVREG_PORT, 10))
	cs = pbsr.NewDistrElectServRegClient(serverConn)

	// Initialize NetCache
	DirtyNetCache = true
}
