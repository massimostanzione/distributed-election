package fredricksonlynch

import (
	. "distributedelection/node/env"
	. "distributedelection/node/tools/monitoring"

	. "distributedelection/node/structure/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
)

func Run() {
	initializeWatchdogs()
	ElectionChannel_fl = make(chan *MsgElectionFL)
	CoordChannel = make(chan *MsgCoordinator)

	go run()
	ListenToIncomingRMI()
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
			SetMonitoringState(MONITORING_HALT)
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
			go startElection()
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
	SetMonitoringState(MONITORING_HALT)
	DirtyNetCache = true
	CurState.Participant = true
	success := sendElection(NewElectionFLMsg(), NextNode)
	// if message was not sent, don't start the watchdog
	if !success {
		SetWatchdog(MSG_ELECTION_FL, true)
	}
}

func vote(inp *MsgElectionFL) *MsgElectionFL {
	var ret *MsgElectionFL
	if !RedudantElectionCheck(CurState.NodeInfo.GetId(), inp) {
		smlog.Trace(LOG_ELECTION, "- voting...")
		ret = inp.AddVoter(CurState.NodeInfo.GetId())
	} else {
		// this case should be already managed in SafeRMI
		// if this else is reached, something is wrong
		smlog.Fatal(LOG_ELECTION, "already voted! something is wrong")
	}
	return ret
}

func elect(candidates []int32) int32 {
	max := candidates[0]
	for _, val := range candidates {
		if val > max {
			max = val
		}
	}
	smlog.Trace(LOG_ELECTION, "Elected node no. %d", max)
	return max
}

func RedudantElectionCheck(voter int32, electionMsg *MsgElectionFL) bool {
	for _, i := range electionMsg.GetVoters() {
		if i == voter {
			return true
		}
	}
	return false
}
