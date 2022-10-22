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
	//TODO processamento parametri va nel main
	host := flag.String("h", "localhost", "host of the node, e.g. \"localhost\", 127.0.0.1 or whatever IP address")
	portParam := flag.String("p", "40043", "target port")
	flag.Parse()
	/*	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}*/
	//port = *portParam
	//addr = "localHost:" + port
	port, _ := strconv.ParseInt(*portParam, 10, 32)
	//smlog.Critical(LOG_UNDEFINED, "%d", int32(port))
	Me.SetPort(int32(port))
	Me.SetHost(*host)
	//addr = "localHost:" + port
	// Parsing input arguments
	/*	filepath := flag.String("f", "../../ILIAD_1STBOOK_IT_ALTERED", "source file to be \"fredricksonLynchp-ed\"")
		substr := flag.String("substr", "Achille", "substr to be searched into the source file")
		serverAddr := flag.String("s", "localHost:40042", "server address and port, in the format ADDRESS:PORT")
		highlight := flag.String("hl", "classic", "[classic/asterisks/none] set substr highlighting in the output\nNOTICE: \"classic\" option may be not available on all systems.")
		help := flag.Bool("help", false, "show this message")

		flag.Parse()
		_, exists := HighlightType[*highlight]
		if !exists {
			fmt.Println("\"-hl\" flag not correctly set.\nSee 'fredricksonLynch -help' for allowed values.")
			os.Exit(-1)
		}
		if *help {
			flag.PrintDefaults()
			os.Exit(0)
		}*/

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
