package statemachine

import (
	. "bully/pkg/node/env"
	. "bully/pkg/node/net"
	. "bully/tools/smlog"
	smlog "bully/tools/smlog"

	//	smlog "bully/smlog"

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

	//TODO per documentazione: STATE_JOINING è obbligatorio, non tanto per la SM quanto per la fase di registrazione
	STATE_JOINING
	STATE_ELECTION
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
	case STATE_ELECTION:
		return "ELECTION"
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
	//	Heartbeat = make(chan *MsgHeartbeat, 1)
	Events = make(chan string, 1)
	ElectionChannel = make(chan *MsgElection, 1)
	OkChannel = make(chan *MsgOk, 1)
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
		}
		smlog.Info(LOG_STATEMACHINE, "Running state cycle")
		//smlog.Println("*** RUNNING STATE: ", (msgType)(currentState))
		switch currentState {
		case STATE_JOINING: // 1
			state_joining()
			break
		case STATE_ELECTION: //2
			state_election()
			break
		case STATE_COORDINATOR: // 3
			state_coordinator()
			break
		case STATE_NON_COORDINATOR: // 4
			state_nonCoordinator()
			break
		default: //TODO
		}
	}
}

func state_joining() {
	for Pause {
	}
	Me = AskForJoining()
	setState(STATE_ELECTION)
	//startElection()
}

func state_election() {
	isElectionStarted := true
	//late_hb_received := 0
	//ELECTION messages alraeady sent in startElection()
	electionTimer := time.NewTimer(ELECTION_ESPIRY + ELECTION_ESPIRY_TOLERANCE)
	for {
		for Pause {
		}
		//TODO check questo approccio su FL
		if isElectionStarted {
			startElection()
		}
		select {
		case <-electionTimer.C:
			setState(STATE_COORDINATOR)
			break
		case <-OkChannel:
			setState(STATE_NON_COORDINATOR)
			break
		// --------------------------------
		// fault tolerance below
		case inp := <-ElectionChannel:
			// other elections are occurring
			if Me.GetId() >= inp.GetStarter() {

				sendOk(NewOkMsg(Me.GetId()), AskForNodeInfo(inp.GetStarter(), false))
				// TODO verificare:
				// no need to start new election, I am already
				// doing mine
			}
			break
		case <-CoordChannel:
			// ignore it: I am doing election anyway
			break
		}
		break
		isElectionStarted = false
	}
	electionTimer.Stop()
}

func state_coordinator() {
	//late_hb_received := 0
	//setState(STATE_ELECTION_VOTER)
	//confirmedCoord := false
	isNewCoord := true
	for {
		if isNewCoord {
			//invia COORD a tutti
			//TODO in FL è implementato diversamente
			for _, dest := range AskForAllNodes() {
				if dest.GetId() != Me.GetId() {
					sendCoord(NewCoordinatorMsg(Me.GetId()), dest)
				}
			}
		}
		for Pause {
		}
		/*if !confirmedCoord {
			go HBroutine()
		}*/
		select {
		case inp := <-ElectionChannel:
			if Me.GetId() > inp.GetStarter() {
				sendOk(NewOkMsg(Me.GetId()), AskForNodeInfo(inp.GetStarter(), false))
				setState(STATE_ELECTION)
				//startElection()

			} else {
				// I am sure that I am no more the coordinator
				setState(STATE_NON_COORDINATOR)
			}
			break

		// --------------------------------
		// fault tolerance below
		case <-OkChannel:
			// ignore
			break
		case inp := <-CoordChannel:
			if inp.GetCoordinator() != Me.GetId() {
				setState(STATE_NON_COORDINATOR)
			}
			break

		}
		//if !confirmedCoord {
		break
		//}
		isNewCoord = false
	}
}

func state_nonCoordinator() {
	//	late_hb_received := 0
	//	nonCoordTimer = time.NewTimer(HB_TIMEOUT + HB_TOLERANCE)
	//go listenHB()
	for {
		if Pause {
			//		nonCoordTimer.Stop()
			for Pause {
			}
			//		nonCoordTimer.Reset(HB_TIMEOUT + HB_TOLERANCE)
		}
		select {
		case inp := <-ElectionChannel:
			if Me.GetId() > inp.GetStarter() {
				sendOk(NewOkMsg(Me.GetId()), AskForNodeInfo(inp.GetStarter(), false))
				setState(STATE_ELECTION)
				//startElection()
			} // else not needed, simply ignore

			break

		// --------------------------------
		// fault tolerance below
		case <-OkChannel:
			// ignore
			break
		case inp := <-CoordChannel:
			if inp.GetCoordinator() == Me.GetId() {
				setState(STATE_COORDINATOR)
			}
			break

		}
		break
	}
}
