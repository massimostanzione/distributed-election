// net.go
package net

import (
	pb "distributedelection/serviceregistry/pb"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

// Port to be listening to. No need to put it on an "env" package.
var Port int

func Listen(host string, port string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.FormatInt(int64(Port), 10))
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to listen to port %v:\n%v", Port, err)
	}
	smlog.Info(LOG_NETWORK, "Listening at %s:%v", GetOutboundIP(), Port)

	// New server instance and service registering
	server := grpc.NewServer()
	pb.RegisterDistrElectServRegServer(server, &DEAServer{})

	// Serve incoming calls
	if err := server.Serve(listener); err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to serve request: %v", err)
	}
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Printf("Could not retrieve IP address: %s", err)
		os.Exit(1)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	split := strings.Split(localAddr.String(), ":")
	return split[0]
}
