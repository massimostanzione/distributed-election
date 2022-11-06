// net.go
package net

import (
	pb "distributedelection/serviceregistry/pb"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

// Port to be listening to. No need to put it on an "env" package.
var Port int

func Listen(host string, port string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.FormatInt(int64(Port), 10))
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to listen to port %v:\n%v", Port, err)
	}
	smlog.Info(LOG_NETWORK, "Listening on port %v...", port)

	// New server instance and service registering
	server := grpc.NewServer()
	pb.RegisterDistrElectServRegServer(server, &DEAServer{})

	// Serve incoming calls
	if err := server.Serve(listener); err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to serve request: %v", err)
	}
}
