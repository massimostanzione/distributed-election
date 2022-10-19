// centrale, mantiene id e indirizzi
package main

import (
	. "bully/pkg/serviceRegistry/net"
	smlog "bully/tools/smlog"
	"flag"
	"fmt"

	//	"net"
	"os"
	//	"google.golang.org/grpc"
)

func main() {
	fmt.Println("*** SERVICE REGISTRY ***")
	fmt.Println("Loading...")
	smlog.InitLogger(true)
	// Parsing input arguments
	//port := flag.String("p", "40042", "port to listen for distgrep requests")
	//workers := flag.String("w", "localhost:40043;localhost:40044;localhost:40045", "addresses and ports of the workers to be bound with, in the following format:\nADDRESS_1:PORT_1;ADDRESS2:PORT_2;...;ADDRESS_N:PORT_N\nMust be between 1 and 15")
	help := flag.Bool("help", false, "show this message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	//StartServiceRegistry()

	// Start listening for incoming calls
	InitializeNetMW()
	// TODO prendere indirizzo da flags
	Listen("localhost", "40042")
}
