// centrale, mantiene id e indirizzi
package main

import (
	. "bully/serviceregistry/pkg/env"
	. "bully/serviceregistry/pkg/net"
	smlog "bully/tools/smlog"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	loadConfig()
	fmt.Println("*** SERVICE REGISTRY ***")
	fmt.Println("Loading...")
	smlog.InitLogger(true, "INFO")
	// Parsing input arguments
	//port := flag.String("p", "40042", "port to listen for distgrep requests")
	//workers := flag.String("w", "localhost:40043;localhost:40044;localhost:40045", "addresses and ports of the workers to be bound with, in the following format:\nADDRESS_1:PORT_1;ADDRESS2:PORT_2;...;ADDRESS_N:PORT_N\nMust be between 1 and 15")

	//StartServiceRegistry()

	// Start listening for incoming calls
	InitializeNetMW()
	Listen("localhost", strconv.FormatInt(int64(Port), 10))
}
func loadConfig() {
	nodeport := flag.Int("p", 40042, "listening port")
	help := flag.Bool("help", false, "show this message")
	Port = *nodeport
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
