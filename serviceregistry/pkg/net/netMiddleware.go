// netMiddleware
package net

//var cs pb.DistGrepClient
import (
	pb "distributedelection/serviceregistry/pb"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"flag"
	"net"

	"google.golang.org/grpc"
)

type DGnode struct {
	pb.UnimplementedDistrElectServRegServer
}

var lis net.Listener
var serverConn *grpc.ClientConn

func InitializeNetMW() {
	liss, err := net.Listen("tcp", serverAddr)
	lis = liss
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while trying to listen to port %v:\n%v", port, err)
	}

	// New server instance and service registering
	w = grpc.NewServer()
	pb.RegisterDistrElectServRegServer(w, &DGnode{})
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
	pb.RegisterDistrElectServRegServer(s, &DGserver{})
	// Serve incoming calls
	if err := s.Serve(lis); err != nil {
		smlog.Fatal(LOG_UNDEFINED, "Error while trying to serve request: %v", err)
	}
}
