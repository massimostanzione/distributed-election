package statemachine

import (
	. "fredricksonLynch/pkg/node/env"
	. "fredricksonLynch/pkg/node/net"
	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"

	//	smlog "fredricksonLynch/smlog"

	//	"log"
	//	"os"

	//"strconv"
	"time"
	//	loggo "github.com/juju/loggo"
)

type nodeState uint8 //TODO spostare altrove?

var currentState = STATE_UNDEFINED

//var nonCoordTimer *time.Timer

const (
	STATE_UNDEFINED nodeState = iota
	STATE_JOINING
	STATE_ELECTION_STARTER
	STATE_ELECTION_VOTER
	STATE_COORDINATOR
	STATE_NON_COORDINATOR
)

//var smlog = log.New(os.Stderr, "[SM] "+time.Now().Format("15:04:05")+" "+currentState.Short()+" "+starterToLogout(starter)+" ", 0)

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
	//	smlog.SetState(stateToLogout(currentState))
	//logger := loggo.GetLogger("")
	//loggo.ConfigureLoggers("eeeeee")
	//logger.SetLogLevel(loggo.TRACE)
	//logger.Infof("new state: ", stateToLogout(currentState))
	smlog.SetStateSMLogger(state.Short())
	smlog.Info(LOG_STATEMACHINE, "new state: %s", currentState.Short())
}

func StartStateMachine() {
	Heartbeat = make(chan *MsgHeartbeat)
	EventsSend = make(chan string)
	EventsList = make(chan string)
	//	Events = make(chan string, 1)
	ElectionChannel = make(chan *MsgElection)
	CoordChannel = make(chan *MsgCoordinator)
	smlog.InitLogger(false)
	smlog.Info(LOG_UNDEFINED, "Starting SM...")
	smlog.InfoU("Type CTRL+C to terminate")
	smlog.InfoU("Type CTRL+Z to Pause/resume")
	InitializeNetMW()

	setState(STATE_JOINING)

	go run()
	Listen()
}

func run() {
	Me = AskForJoining()
	startElection()
	for {
		for Pause {
			// test
		}
		smlog.Info(LOG_STATEMACHINE, "Running state cycle")
		//smlog.Println("*** RUNNING STATE: ", (msgType)(currentState))
		switch currentState {
		case STATE_JOINING: // 1
			state_joining()
			break
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
		default:
		}
	}
}
func setWaiting(msgType MsgType, active bool) {
	WaitingMap[msgType].Waiting = active
	if active {
		WaitingMap[msgType].Timer.Reset(IDLE_WAIT_LIMIT * time.Second)
	} else {
		WaitingMap[msgType].Timer.Stop()
	}
}

func state_joining() {
	for Pause {
	}
	// WaitingMap già inizializzata prima
	for {
		smlog.InfoU("attendo messaggi...")
		select {
		case in := <-ElectionChannel:
			//setHbMgt(HB_HALT)
			if in.GetStarter() == Me.GetId() {
				smlog.InfoU("tornato election partito da me")
				setWaiting(MSG_ELECTION, false)
				coord := elect(in.GetVoters())
				smlog.InfoU("eletto")
				CoordId = coord
				// TODO simmetria nella gestione + cache locale sul successivo
				go sendCoord(NewCoordinatorMsg(Me.GetId(), CoordId), NextNode)
				setHbMgt(HB_SEND)
				setWaiting(MSG_COORDINATOR, true)
			} else {
				smlog.InfoU("arrivato election non mio")
				vote(in)
				smlog.InfoU("ho votato")
				//è incluso in vote, da portare fuori
				//sendElection(in, NextNode)
			}
			break
		case in := <-CoordChannel:
			//setHbMgt(HB_HALT)
			smlog.InfoU("coordch")
			if in.GetStarter() == Me.GetId() {
				setWaiting(MSG_COORDINATOR, false)
			} else {
				go sendCoord(in, NextNode)
			}
			CoordId = in.GetCoordinator() // duplicato?
			if CoordId == Me.GetId() {
				smlog.InfoU("sono il nuovo coord!")
				setHbMgt(HB_SEND)
			} else {
				smlog.InfoU("NON sono il nuovo coord!")
				setHbMgt(HB_LISTEN)
			}
			break

		case <-WaitingMap[MSG_ELECTION].Timer.C:
			smlog.Trace(LOG_UNDEFINED, "scaduto timer E")
			setWaiting(MSG_ELECTION, false)
			startElection()
			break
		case <-WaitingMap[MSG_COORDINATOR].Timer.C:
			smlog.Trace(LOG_UNDEFINED, "scaduto timer C")
			setWaiting(MSG_COORDINATOR, false)
			// AskForNextNode
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
	sendElection(NewElectionMsg(Me.GetId()), NextNode)
	setWaiting(MSG_ELECTION, true)
}
