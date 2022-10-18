// centrale, mantiene id e indirizzi
package main

import (
	"flag"
	"fredricksonLynch/pkg/serviceRegistry/net"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	log.Println("*** DISTGREP SERVER ***")
	// Parsing input arguments
	port := flag.String("p", "40042", "port to listen for distgrep requests")
	//workers := flag.String("w", "localhost:40043;localhost:40044;localhost:40045", "addresses and ports of the workers to be bound with, in the following format:\nADDRESS_1:PORT_1;ADDRESS2:PORT_2;...;ADDRESS_N:PORT_N\nMust be between 1 and 15")
	help := flag.Bool("help", false, "show this message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	/*
		workersArray = strings.Split(*workers, ";")
		if len(workersArray) < MIN_WORKERSNO || len(workersArray) > MAX_WORKERSNO {
			fmt.Printf("Workers to be bound with must be between %v and %v\n", MIN_WORKERSNO, MAX_WORKERSNO)
			os.Exit(-1)
		}
		log.Println("Will bind (on-demand) to the following workers:")
		for i := range workersArray {
			log.Printf("- %v", workersArray[i])
		}*/
	// Start listening for incoming calls

	// TODO sistemare queste chiamate
	lis, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatalf("Error while trying to listen to port %v:\n%v", *port, err)
	}
	log.Printf("--------------------------")

	log.Printf("Listening on port %v...", *port)
	// New server instance and service registering
	s := grpc.NewServer()
	pb.RegisterDistGrepServer(s, &DGserver{})
	// Serve incoming calls
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while trying to serve request: %v", err)
	}
}
