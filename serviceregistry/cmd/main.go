package main

import (
	. "distributedelection/serviceregistry/pkg/net"
	smlog "distributedelection/tools/smlog"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("+------------------------------------------------+")
	fmt.Println("|         DISTRIBUTED ELECTION ALGORITHMS        |")
	fmt.Println("|github.com/massimostanzione/distributed-election|")
	fmt.Println("+------------------------------------------------+")
	fmt.Println("|                Service Registry                |")
	fmt.Println("+------------------------------------------------+")
	fmt.Println("")
	fmt.Println("Loading configuration environment...")
	loadConfig()
	fmt.Println("... done. Starting...")

	smlog.Initialize(true, "INFO")

	//InitializeNetMW()
	Listen("localhost", strconv.FormatInt(int64(Port), 10))
}

func loadConfig() {
	nodeport := flag.Int("p", 40049, "listening port")
	help := flag.Bool("help", false, "show this message")
	flag.Parse()
	Port = *nodeport
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
