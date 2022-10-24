// netMiddleware
package net

//var cs pb.DistGrepClient
import (
	pb "bully/serviceregistry/pb"
	. "bully/serviceregistry/pkg/env"
	. "bully/tools/smlog"
	smlog "bully/tools/smlog"
	"net"

	"strconv"

	"google.golang.org/grpc"
)

type DGnode struct {
	pb.UnimplementedDistGrepServer
}

//var NONE = &pb.NONE{}
//var cs pb.DistGrepClient

var w *grpc.Server
var lis net.Listener
var serverConn *grpc.ClientConn //server

const RMI_RETRY_TOLERANCE = 3
const LATE_HB_TOLERANCE = 3

func InitializeNetMW() {

	serverAddr := "localHost:" + strconv.FormatInt(int64(Port), 10)
	conn := ConnectToNode(serverAddr)
	serverConn = conn

	liss, err := net.Listen("tcp", serverAddr)
	lis = liss
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while trying to listen to port %v:\n%v", Port, err)
	}

	// New server instance and service registering
	w = grpc.NewServer()
	pb.RegisterDistGrepServer(w, &DGnode{})
	// Serve incoming calls

	// Defining client interface, to be used to invoke the bully service
	//cs = pb.NewDistGrepClient(serverConn)
}

func ConnectToNode(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while contacting server on %v:\n %v", addr, err)
	}
	return conn
}
func Listen(host string, port string) {

	smlog.InfoU("Listening on port %v...", port)
	// New server instance and service registering
	s := grpc.NewServer()
	pb.RegisterDistGrepServer(s, &DGserver{})
	// Serve incoming calls
	if err := s.Serve(lis); err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while trying to serve request: %v", err)
	}

	//	for pause {
	//	}
}
