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

var nonCoordTimer *time.Timer

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
	Heartbeat = make(chan *MsgHeartbeat, 1)
	Events = make(chan string, 1)
	ElectionChannel = make(chan *MsgElection, 1)
	CoordChannel = make(chan *MsgCoordinator, 1)
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
		default: //TODO
		}
	}
}
func setWaiting(msgType MsgType, val bool) {
	WaitingMap[msgType].Waiting = val
	if val {
		WaitingMap[msgType].Timer.Reset(5 * time.Second)
	} else {
		WaitingMap[msgType].Timer.Stop()
	}
}

func state_joining() {
	for Pause {
	}
	Me = AskForJoining()
	// WaitingMap già inizializzata prima
	startElection()
	for {
		select {
		case in := <-ElectionChannel:
			if in.GetStarter() == Me.GetId() {
				smlog.InfoU("tornato election partito da me")
				setWaiting(MSG_ELECTION, false)
				coord := elect(in.GetVoters())
				smlog.InfoU("eletto")
				CoordId = coord
				// simmetria nella gestione
				sendCoord(NewCoordinatorMsg(Me.GetId(), CoordId), AskForNodeInfo(Me.GetId()+1, false))
				setWaiting(MSG_COORDINATOR, true)
			} else {
				smlog.InfoU("arrivato election non mio")
				vote(in)
				smlog.InfoU("ho votato")
				Events <- "STOP"
				sendElection(in, AskForNodeInfo(Me.GetId()+1, false))
			}
			break
		case in := <-CoordChannel:
			smlog.InfoU("coordch")
			if in.GetStarter() == Me.GetId() {
				setWaiting(MSG_COORDINATOR, false)
			}
			CoordId = in.GetCoordinator() // duplicato?
			if CoordId == Me.GetId() {
				smlog.InfoU("sono il nuovo coord!")
				SendingHB = true
				go InviaHB()
			} else {
				smlog.InfoU("NON sono il nuovo coord!")
				ListeningtoHb = true
				go ListenToHb()
			}
			break
			/*
				case <-WaitingMap[MSG_ELECTION].Timer.C:
					smlog.Trace(LOG_UNDEFINED, "scaduto timer E")
					setWaiting(MSG_ELECTION, false)
					startElection()
					break
				case <-WaitingMap[MSG_COORDINATOR].Timer.C:
					smlog.Trace(LOG_UNDEFINED, "scaduto timer C")
					// AskForNextNode
					setWaiting(MSG_COORDINATOR, false)
					sendCoord(NewCoordinatorMsg(Me.GetId(), CoordId), AskForNodeInfo(Me.GetId()+1, false))
					break
			*/
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
	sendElection(NewElectionMsg(Me.GetId()), AskForNodeInfo(Me.GetId()+1, false))
	setWaiting(MSG_ELECTION, true)
}

// da separare nel comportamento
func ListenToHb() {
	interrupt := false
	noncoordTimer := time.NewTicker(HB_TIMEOUT + HB_TOLERANCE)
	for {
		select {
		case <-Heartbeat:
			noncoordTimer.Reset(HB_TIMEOUT + HB_TOLERANCE)
			break
		case <-noncoordTimer.C:
			go startElection()
			Events <- "STOP"
			interrupt = true
			break
		}
		if interrupt {
			break
		}
	}
	smlog.Info(LOG_HB, "Esco dalla routine di invio HB...")
	SendingHB = false
}
