package main

import (
	. "bully/pkg/node/env"
	. "bully/pkg/node/net"
	. "bully/pkg/node/statemachine"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	StartStateMachine()
}
