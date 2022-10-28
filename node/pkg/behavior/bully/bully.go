package bully

import (
	. "distributedelection/node/pkg/behavior/monitoring"
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

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

func Run() {
	Heartbeat = make(chan *MsgHeartbeat, 1)
	Events = make(chan string, 1)
	ElectionChannel = make(chan *MsgElection, 1)
	OkChannel = make(chan *MsgOk, 1)
	CoordChannel = make(chan *MsgCoordinator, 1)
	smlog.InitLogger(false, Cfg.TERMINAL_SMLOG_LEVEL)
	smlog.Info(LOG_UNDEFINED, "Starting SM...")
	smlog.InfoU("Type CTRL+C to terminate")
	InitializeNetMW()

	setState(STATE_JOINING)
	go run()
	Listen()
}

func run() {
	//	initHbMgt()
	state_joining()
}

func state_joining() {
	Me = AskForJoining()

	IsElectionStarted = true
	for {
		if IsElectionStarted {
			SetMonitoringState(HB_HALT)
			//ELECTION messages alraeady sent in startElection()
			ElectionTimer = time.NewTimer(time.Duration(Cfg.ELECTION_ESPIRY+Cfg.ELECTION_ESPIRY_TOLERANCE) * time.Millisecond)
			go startElection()
			IsElectionStarted = false
		}
		select {
		case <-ElectionTimer.C:
			ElectionTimer.Stop()
			SetMonitoringState(HB_SEND)
			smlog.InfoU("sono il coord, autoproclamato")
			CoordId = Me.GetId()
			for _, dest := range AskForAllNodes() {
				if dest.GetId() != Me.GetId() {
					go sendCoord(NewCoordinatorMsg(Me.GetId(), Me.GetId()), dest)
				}
			}
			IsElectionStarted = false
			break
		case inp := <-ElectionChannel:
			//SetMonitoringState(HB_HALT)
			//electionTimer.Stop()
			smlog.InfoU("arrivato E")
			// other elections are occurring
			if Me.GetId() > inp.GetStarter() {
				go sendOk(NewOkMsg(Me.GetId()), AskForNodeInfo(inp.GetStarter()))
				IsElectionStarted = true
				//ElectionTimer.Reset(ELECTION_ESPIRY + ELECTION_ESPIRY_TOLERANCE)
			} else {

				IsElectionStarted = false
			}
			break
		case <-OkChannel:
			//ElectionTimer.Stop()
			smlog.InfoU("qualcuno è più bully di me")
			//time.Sleep(ELECTION_ESPIRY + ELECTION_ESPIRY_TOLERANCE)
			ElectionTimer.Stop()
			SetMonitoringState(HB_HALT)
			IsElectionStarted = false
			break
		case inp := <-CoordChannel:
			//electionTimer.Stop()
			CoordId = inp.GetCoordinator()
			smlog.InfoU("arrivato C")
			if inp.GetCoordinator() != Me.GetId() {
				smlog.InfoU("nuovo coord è %d", inp.GetCoordinator())
				SetMonitoringState(HB_LISTEN)
			} else {
				smlog.Fatal(LOG_UNDEFINED, "unreachable")
			}
			IsElectionStarted = false
			break

		}
	}
}
