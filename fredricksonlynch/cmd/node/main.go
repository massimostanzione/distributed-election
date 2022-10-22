package main

import (
	"flag"
	. "fredricksonLynch/pkg/node/env"
	. "fredricksonLynch/pkg/node/net"
	. "fredricksonLynch/pkg/node/statemachine"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

/*
func printNodeInfo(node *SMNode) {
	log.Printf("ID:\t\t%d", GetId())
	log.Printf("Address:\t%s", GetFullAddr())
}
*/
func main() {
	loadConfig()
	Sigchan = make(chan os.Signal, 1)
	signal.Notify(Sigchan, syscall.SIGTSTP)
	go func() {
		sig := <-Sigchan
		log.Printf("SEGNALE!", sig)
		Pause = !Pause
		SwitchServerState(Pause)
		log.Printf("Pause = %v", Pause)
	}()
	initializeWaitingMap()
	StartStateMachine()
}
func loadConfig() {
	//nodeaddr := flag.String("a", "localhost:40042", "complete address of the node (host:port), e.g. \"localhost:40043\", 127.0.0.1:40043")
	nodehost := flag.String("h", "localhost", "host of the node, e.g. \"localhost\", 127.0.0.1 or whatever IP address")
	nodeport := flag.Int("p", 40043, "target port")
	servicereghost := flag.String("sh", "localhost", "host of the service registry, e.g. \"localhost\", 127.0.0.1 or whatever IP address")
	serviceregport := flag.Int64("sp", 40042, "target port of the service registry")
	help := flag.Bool("help", false, "show this message")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	/*
		if *nodeaddr {
			split := strings.Split(*host, ":")
			Me.SetHost(split[0])
			Me.SetPort(int32(split[1]))

		}*/

	Me.SetHost(*nodehost)
	Me.SetPort(int32(*nodeport))
	ServRegAddr = *servicereghost + ":" + strconv.FormatInt(*serviceregport, 10)
}
func initializeWaitingMap() {
	WaitingMap = map[MsgType]*WaitingStruct{
		MSG_ELECTION: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(IDLE_WAIT_LIMIT * time.Second),
		},
		MSG_COORDINATOR: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(IDLE_WAIT_LIMIT * time.Second),
		},
	}
	WaitingMap[MSG_ELECTION].Timer.Stop()
	WaitingMap[MSG_COORDINATOR].Timer.Stop()
}
