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
	State.NodeInfo = AskForJoining()
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

func state_joining() {
	// WaitingMap già inizializzata prima
	for {
		smlog.Debug(LOG_STATEMACHINE, "Waiting for messages...")
		select {
		case in := <-ElectionChannel:
			State.Participant = true
			smlog.Debug(LOG_STATEMACHINE, "Handling ELECTION message")
			if in.GetStarter() == State.NodeInfo.GetId() {
				SetWaiting(MSG_ELECTION, false)
				coord := elect(in.GetVoters())
				State.Coordinator = coord
				go sendCoord(NewCoordinatorMsg(State.NodeInfo.GetId(), State.Coordinator), NextNode)
				SetWaiting(MSG_COORDINATOR, true)
			} else {
				voted := vote(in)
				go sendElection(voted, NextNode)
			}
			break
		case in := <-CoordChannel:
			smlog.Debug(LOG_STATEMACHINE, "Handling COORDINATOR message")
			SetMonitoringState(MONITORING_HALT)
			if in.GetStarter() == State.NodeInfo.GetId() {
				SetWaiting(MSG_COORDINATOR, false)
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
		case <-WaitingMap[MSG_ELECTION].Timer.C:
			smlog.Error(LOG_NETWORK, "ELECTION message not returned back within time limit. Starting new election...")
			SetWaiting(MSG_ELECTION, false)
			startElection()
			break
		case <-WaitingMap[MSG_COORDINATOR].Timer.C:
			smlog.Error(LOG_NETWORK, "COORDINATOR message not returned back within time limit. Sending it again...")
			SetWaiting(MSG_COORDINATOR, false)
			sendCoord(NewCoordinatorMsg(State.NodeInfo.GetId(), State.Coordinator), NextNode)
			break

		}
	}
}

func startElection() {
	/*
		l'inoltro al successivo è già gestito in safeRMI,
		difatti il successivo parametro success non serve
		for {
			i := State.NodeInfo.GetId() + 1
			next := AskForNodeInfo(i, false)
			if next.GetId() == State.NodeInfo.GetId() {
				break
			}
			//TODO mettere starter implicito
			success := sendElection(NewElectionMsg(State.NodeInfo.GetId()), next)
			if !success {
				break
			}
			i++
		}*/
	State.Participant = true
	err := sendElection(NewElectionMsg(), NextNode)
	if !err {
		SetWaiting(MSG_ELECTION, true)
	}
}
