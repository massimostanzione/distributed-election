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

	STATE_JOINING
	/*STATE_ELECTION
	STATE_COORDINATOR
	STATE_NON_COORDINATOR*/
)

func (state nodeState) Short() string {
	switch state {
	case STATE_UNDEFINED:
		return "N/A  "
	case STATE_JOINING:
		return "JOINI"
		/*case STATE_ELECTION:
			return "ELECTION"
		case STATE_COORDINATOR:
			return "COORD"
		case STATE_NON_COORDINATOR:
			return "NONCO"*/
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
	initHbMgt()
	state_joining()
}

func state_joining() {
	for Pause {
	}
	Me = AskForJoining()

	isElectionStarted := true
	//ELECTION messages alraeady sent in startElection()
	electionTimer := time.NewTimer(ELECTION_ESPIRY + ELECTION_ESPIRY_TOLERANCE)
	for {
		for Pause {
		}
		if isElectionStarted {
			startElection()
		}
		select {
		case <-electionTimer.C:
			electionTimer.Stop()
			smlog.InfoU("sono il coord, autoproclamato")
			for _, dest := range AskForAllNodes() {
				if dest.GetId() != Me.GetId() {
					go sendCoord(NewCoordinatorMsg(Me.GetId()), dest)
				}
			}
			setHbMgt(HB_SEND)
			break
		case <-OkChannel:
			electionTimer.Stop()
			smlog.InfoU("qualcuno è più bully di me")
			setHbMgt(HB_LISTEN)
			break
		case inp := <-ElectionChannel:
			setHbMgt(HB_HALT)
			electionTimer.Stop()
			smlog.InfoU("arrivato E")
			// other elections are occurring
			if Me.GetId() > inp.GetStarter() {
				sendOk(NewOkMsg(Me.GetId()), AskForNodeInfo(inp.GetStarter()))
				isElectionStarted = true
				//startElection()
				electionTimer.Reset(ELECTION_ESPIRY + ELECTION_ESPIRY_TOLERANCE)
			}
			break
		case inp := <-CoordChannel:
			electionTimer.Stop()
			smlog.InfoU("arrivato C")
			if inp.GetCoordinator() != Me.GetId() {
				smlog.InfoU("nuovo coord è %d", inp.GetCoordinator())
			} else {
				smlog.Fatal(LOG_UNDEFINED, "unreachable")
			}
			setHbMgt(HB_LISTEN)
			break
		}
		isElectionStarted = false
	}
}
