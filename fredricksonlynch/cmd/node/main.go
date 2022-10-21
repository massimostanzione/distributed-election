package main

import (
	. "fredricksonLynch/pkg/node/env"
	. "fredricksonLynch/pkg/node/net"
	. "fredricksonLynch/pkg/node/statemachine"
	"log"
	"os"
	"os/signal"
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
	Sigchan = make(chan os.Signal, 1)
	signal.Notify(Sigchan, syscall.SIGTSTP)
	go func() {
		sig := <-Sigchan
		log.Printf("SEGNALE!", sig)
		Pause = !Pause
		SwitchServerState(Pause)
		log.Printf("Pause = %v", Pause)
	}()

	WaitingMap = map[MsgType]*WaitingStruct{
		MSG_ELECTION: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(5 * time.Second),
		},
		MSG_COORDINATOR: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(5 * time.Second),
		},
	}
	WaitingMap[MSG_ELECTION].Timer.Stop()
	WaitingMap[MSG_COORDINATOR].Timer.Stop()
	StartStateMachine()
}
