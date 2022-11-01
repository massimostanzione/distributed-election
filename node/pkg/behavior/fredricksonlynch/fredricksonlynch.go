package fredricksonlynch

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/monitoring"

	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
)

func Run() {
	initializeWatchdogs()
	ElectionChannel = make(chan *MsgElection)
	CoordChannel = make(chan *MsgCoordinator)
	smlog.InfoU("Starting SM...")
	smlog.InfoU("Type CTRL+C to terminate")

	go run()
	Listen()
}

func initializeWatchdogs() {
	Watchdogs = map[MsgType]*Watchdog{
		MSG_ELECTION: &Watchdog{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
		MSG_COORDINATOR: &Watchdog{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
	}
	Watchdogs[MSG_ELECTION].Timer.Stop()
	Watchdogs[MSG_COORDINATOR].Timer.Stop()
}
func run() {
	State.NodeInfo = AskForJoining()
	go startElection()
	for {
		smlog.Debug(LOG_STATEMACHINE, "Waiting for messages...")
		select {
		case in := <-ElectionChannel:
			State.Participant = true
			smlog.Debug(LOG_STATEMACHINE, "Handling ELECTION message")
			if in.GetStarter() == State.NodeInfo.GetId() {
				SetWatchdog(MSG_ELECTION, false)
				coord := elect(in.GetVoters())
				State.Coordinator = coord
				go sendCoord(NewCoordinatorMsg(State.NodeInfo.GetId(), State.Coordinator), NextNode)
				SetWatchdog(MSG_COORDINATOR, true)
			} else {
				voted := vote(in)
				go sendElection(voted, NextNode)
			}
			break
		case in := <-CoordChannel:
			smlog.Debug(LOG_STATEMACHINE, "Handling COORDINATOR message")
			SetMonitoringState(MONITORING_HALT)
			if in.GetStarter() == State.NodeInfo.GetId() {
				SetWatchdog(MSG_COORDINATOR, false)
			} else {
				go sendCoord(in, NextNode)
			}
			State.Coordinator = in.GetCoordinator()
			State.Participant = false
			if State.Coordinator == State.NodeInfo.GetId() {
				smlog.Info(LOG_ELECTION, "*** I am the new coordinator ***")
				SetMonitoringState(MONITORING_SEND)
			} else {
				smlog.Trace(LOG_ELECTION, "I am NOT the new coordinator")
				SetMonitoringState(MONITORING_LISTEN)
			}
			break
		case <-MonitoringChannel:
			SetMonitoringState(MONITORING_HALT)
			smlog.Critical(LOG_ELECTION, "non sento più, ricomincio elez")
			go startElection()
			break
		case <-Watchdogs[MSG_ELECTION].Timer.C:
			smlog.Error(LOG_NETWORK, "ELECTION message not returned back within time limit. Starting new election...")
			SetWatchdog(MSG_ELECTION, false)
			startElection()
			break
		case <-Watchdogs[MSG_COORDINATOR].Timer.C:
			smlog.Error(LOG_NETWORK, "COORDINATOR message not returned back within time limit. Sending it again...")
			SetWatchdog(MSG_COORDINATOR, false)
			sendCoord(NewCoordinatorMsg(State.NodeInfo.GetId(), State.Coordinator), NextNode)
			break

		}
	}
}
func startElection() {

	DirtyNetList = true
	State.Participant = true
	err := sendElection(NewElectionMsg(), NextNode)
	if !err {
		SetWatchdog(MSG_ELECTION, true)
	}
}
