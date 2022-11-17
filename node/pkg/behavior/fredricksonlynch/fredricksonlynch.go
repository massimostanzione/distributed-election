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
	ElectionChannel_fl = make(chan *MsgElectionFL)
	CoordChannel = make(chan *MsgCoordinator)

	go run()
	Listen()
}

func initializeWatchdogs() {
	Watchdogs = map[MsgType]*Watchdog{
		MSG_ELECTION_FL: &Watchdog{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
		MSG_COORDINATOR: &Watchdog{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
	}
	Watchdogs[MSG_ELECTION_FL].Timer.Stop()
	Watchdogs[MSG_COORDINATOR].Timer.Stop()
}

func run() {
	CurState.NodeInfo = AskForJoining()
	startElection()
	for {
		select {
		case in := <-ElectionChannel_fl:
			smlog.Debug(LOG_ELECTION, "Handling ELECTION message")
			CurState.Participant = true
			if in.GetStarter() == CurState.NodeInfo.GetId() {
				SetWatchdog(MSG_ELECTION_FL, false)
				coord := elect(in.GetVoters())
				CurState.Coordinator = coord
				go sendCoord(NewCoordinatorMsg(CurState.NodeInfo.GetId(), CurState.Coordinator), NextNode)
				SetWatchdog(MSG_COORDINATOR, true)
			} else {
				voted := vote(in)
				go sendElection(voted, NextNode)
			}
			break
		case in := <-CoordChannel:
			smlog.Debug(LOG_STATEMACHINE, "Handling COORDINATOR message")
			SetMonitoringState(MONITORING_HALT)
			SetWatchdog(MSG_ELECTION_FL, false)
			if in.GetStarter() == CurState.NodeInfo.GetId() {
				SetWatchdog(MSG_COORDINATOR, false)
			} else {
				go sendCoord(in, NextNode)
			}
			CurState.Coordinator = in.GetCoordinator()
			CurState.Participant = false
			if CurState.Coordinator == CurState.NodeInfo.GetId() {
				smlog.Info(LOG_ELECTION, "*** I am the new coordinator ***")
				SetMonitoringState(MONITORING_SEND)
			} else {
				smlog.Trace(LOG_ELECTION, "I am NOT the new coordinator")
				SetMonitoringState(MONITORING_LISTEN)
			}
			break
		case <-MonitoringChannel:
			smlog.Critical(LOG_ELECTION, "Coordinator failed!")
			SetMonitoringState(MONITORING_HALT)
			go startElection()
			break
		case <-Watchdogs[MSG_ELECTION_FL].Timer.C:
			smlog.Error(LOG_NETWORK, "ELECTION message not returned back within time limit. Starting new election...")
			SetWatchdog(MSG_ELECTION_FL, false)
			startElection()
			break
		case <-Watchdogs[MSG_COORDINATOR].Timer.C:
			smlog.Error(LOG_NETWORK, "COORDINATOR message not returned back within time limit. Sending it again...")
			SetWatchdog(MSG_COORDINATOR, false)
			go sendCoord(NewCoordinatorMsg(CurState.NodeInfo.GetId(), CurState.Coordinator), NextNode)
			break

		}
	}
}
func startElection() {

	DirtyNetList = true
	CurState.Participant = true
	err := sendElection(NewElectionFLMsg(), NextNode)
	if !err {
		SetWatchdog(MSG_ELECTION_FL, true)
	}
}
