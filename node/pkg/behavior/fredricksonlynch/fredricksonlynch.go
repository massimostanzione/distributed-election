package fredricksonlynch

import (
	. "distributedelection/node/pkg/behavior/monitoring"
	. "distributedelection/node/pkg/env"

	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
)

type nodeState uint8 //TODO spostare altrove?

var currentState = STATE_UNDEFINED

const (
	STATE_UNDEFINED nodeState = iota
	STATE_JOINING
	STATE_ELECTION_STARTER
	STATE_ELECTION_VOTER
	STATE_COORDINATOR
	STATE_NON_COORDINATOR
)

func (state nodeState) Short() string {
	switch state {
	case STATE_UNDEFINED:
		return "N/A  "
	case STATE_JOINING:
		return "JOINI"
	case STATE_ELECTION_STARTER:
		return "START"
	case STATE_ELECTION_VOTER:
		return "VOTER"
	case STATE_COORDINATOR:
		return "COORD"
	case STATE_NON_COORDINATOR:
		return "NONCO"
	}
	return "err"
}

func setState(state nodeState) {
	currentState = state
	smlog.SetStateSMLogger(state.Short())
	smlog.Info(LOG_STATEMACHINE, "new state: %s", currentState.Short())
}

func Run() {
	initializeWaitingMap()
	ElectionChannel = make(chan *MsgElection)
	CoordChannel = make(chan *MsgCoordinator)
	smlog.InfoU("Starting SM...")
	smlog.InfoU("Type CTRL+C to terminate")

	setState(STATE_JOINING)

	go run()
	Listen()
}

func initializeWaitingMap() {
	WaitingMap = map[MsgType]*WaitingStruct{
		MSG_ELECTION: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
		MSG_COORDINATOR: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
	}
	WaitingMap[MSG_ELECTION].Timer.Stop()
	WaitingMap[MSG_COORDINATOR].Timer.Stop()
}
func run() {
	Me = AskForJoining()
	go startElection()
	//	for {
	smlog.Info(LOG_STATEMACHINE, "Running state cycle")
	//smlog.Println("*** RUNNING STATE: ", (msgType)(currentState))
	//		switch currentState {
	//		case STATE_JOINING: // 1
	state_joining()
	//			break
	/*case STATE_ELECTION_STARTER: // 2
		state_election_starter()
		break
	case STATE_ELECTION_VOTER: // 3
		state_election_voter()
		break
	case STATE_COORDINATOR: // 4
		state_coordinator()
		break
	case STATE_NON_COORDINATOR: // 5
		state_nonCoordinator()
		break*/
	//		default:
	//	break
	//		}
	//	}
}
func setWaiting(msgType MsgType, active bool) {
	WaitingMap[msgType].Waiting = active
	if active {
		WaitingMap[msgType].Timer.Reset(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Millisecond)
	} else {
		WaitingMap[msgType].Timer.Stop()
	}
}

func state_joining() {
	// WaitingMap già inizializzata prima
	for {
		smlog.Debug(LOG_STATEMACHINE, "Waiting for messages...")
		select {
		case in := <-ElectionChannel:
			smlog.Debug(LOG_STATEMACHINE, "Handling ELECTION message")
			smlog.Debug(LOG_STATEMACHINE, "setmonitoringstate to HALT")
			if in.GetStarter() == Me.GetId() {
				setWaiting(MSG_ELECTION, false)
				coord := elect(in.GetVoters())
				CoordId = coord
				go sendCoord(NewCoordinatorMsg(Me.GetId(), CoordId), NextNode)
				setWaiting(MSG_COORDINATOR, true)
			} else {
				smlog.Debug(LOG_STATEMACHINE, "voting")
				voted := vote(in)
				smlog.Debug(LOG_STATEMACHINE, "voted")
				go sendElection(voted, NextNode)
				smlog.Debug(LOG_STATEMACHINE, "goroutine sent started")
			}
			break
		case in := <-CoordChannel:
			smlog.Debug(LOG_STATEMACHINE, "Handling COORDINATOR message")
			SetMonitoringState(MONITORING_HALT)
			if in.GetStarter() == Me.GetId() {
				setWaiting(MSG_COORDINATOR, false)
			} else {
				go sendCoord(in, NextNode)
			}
			CoordId = in.GetCoordinator()
			if CoordId == Me.GetId() {
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
		case <-WaitingMap[MSG_ELECTION].Timer.C:
			smlog.Error(LOG_NETWORK, "ELECTION message not returned back within time limit. Starting new election...")
			setWaiting(MSG_ELECTION, false)
			startElection()
			break
		case <-WaitingMap[MSG_COORDINATOR].Timer.C:
			smlog.Error(LOG_NETWORK, "COORDINATOR message not returned back within time limit. Sending it again...")
			setWaiting(MSG_COORDINATOR, false)
			sendCoord(NewCoordinatorMsg(Me.GetId(), CoordId), NextNode)
			break

		}
	}
}

func startElection() {
	/*
		l'inoltro al successivo è già gestito in safeRMI,
		difatti il successivo parametro success non serve
		for {
			i := Me.GetId() + 1
			next := AskForNodeInfo(i, false)
			if next.GetId() == Me.GetId() {
				break
			}
			//TODO mettere starter implicito
			success := sendElection(NewElectionMsg(Me.GetId()), next)
			if !success {
				break
			}
			i++
		}*/
	err := sendElection(NewElectionMsg(), NextNode)
	if !err {
		setWaiting(MSG_ELECTION, true)
	}
}
