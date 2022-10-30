package bully

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/monitoring"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
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
	initializeWaitingMap()
	ElectionChannel = make(chan *MsgElection)
	OkChannel = make(chan *MsgOk)
	CoordChannel = make(chan *MsgCoordinator)
	smlog.Initialize(false, Cfg.TERMINAL_SMLOG_LEVEL)
	smlog.Info(LOG_UNDEFINED, "Starting SM...")
	smlog.InfoU("Type CTRL+C to terminate")

	setState(STATE_JOINING)
	go run()
	Listen()
}
func initializeWaitingMap() {
	WaitingMap = map[MsgType]*WaitingStruct{
		/*MSG_ELECTION: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},*/
		MSG_COORDINATOR: &WaitingStruct{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
	}
	//WaitingMap[MSG_ELECTION].Timer.Stop()
	WaitingMap[MSG_COORDINATOR].Timer.Stop()
}

func run() {
	//	initHbMgt()
	state_joining()
}

func state_joining() {
	State.NodeInfo = AskForJoining()

	IsElectionStarted = true
	for {
		if IsElectionStarted {
			SetMonitoringState(MONITORING_HALT)
			//ELECTION messages alraeady sent in startElection()
			ElectionTimer = time.NewTimer(time.Duration(Cfg.ELECTION_ESPIRY+Cfg.ELECTION_ESPIRY_TOLERANCE) * time.Millisecond)
			go startElection()
			IsElectionStarted = false
		}
		select {
		case <-ElectionTimer.C:
			ElectionTimer.Stop()
			SetMonitoringState(MONITORING_SEND)
			SetWaiting(MSG_COORDINATOR, false)
			State.Participant = false
			smlog.InfoU("sono il coord, autoproclamato")
			State.Coordinator = State.NodeInfo.GetId()
			for _, dest := range AskForAllNodes() {
				if dest.GetId() != State.NodeInfo.GetId() {
					go sendCoord(NewCoordinatorMsg(State.NodeInfo.GetId(), State.NodeInfo.GetId()), dest)
				}
			}
			IsElectionStarted = false
			break
		case inp := <-ElectionChannel:
			State.Participant = true
			SetMonitoringState(MONITORING_HALT)
			//electionTimer.Stop()
			SetWaiting(MSG_COORDINATOR, true)
			smlog.InfoU("arrivato E")
			// other elections are occurring
			if State.NodeInfo.GetId() > inp.GetStarter() {
				go sendOk(NewOkMsg(State.NodeInfo.GetId()), AskForNodeInfo(inp.GetStarter()))
				IsElectionStarted = true
				//ElectionTimer.Reset(ELECTION_ESPIRY + ELECTION_ESPIRY_TOLERANCE)
			} else {

				IsElectionStarted = false
			}
			break
		case <-OkChannel:
			State.Participant = true
			SetMonitoringState(MONITORING_HALT)
			//ElectionTimer.Stop()
			SetWaiting(MSG_COORDINATOR, true)
			smlog.InfoU("qualcuno è più bully di me")
			//time.Sleep(ELECTION_ESPIRY + ELECTION_ESPIRY_TOLERANCE)
			ElectionTimer.Stop()
			SetMonitoringState(MONITORING_HALT)
			IsElectionStarted = false
			break
		case inp := <-CoordChannel:
			SetMonitoringState(MONITORING_HALT)
			SetWaiting(MSG_COORDINATOR, false)
			//electionTimer.Stop()
			State.Coordinator = inp.GetCoordinator()
			smlog.InfoU("arrivato C")
			if inp.GetCoordinator() != State.NodeInfo.GetId() {
				smlog.InfoU("nuovo coord è %d", inp.GetCoordinator())
				SetMonitoringState(MONITORING_LISTEN)
				State.Participant = false
			} else {
				smlog.Fatal(LOG_UNDEFINED, "unreachable")
			}
			IsElectionStarted = false
			break
		case <-MonitoringChannel:
			SetMonitoringState(MONITORING_HALT)
			startElection()
			break
		case <-WaitingMap[MSG_COORDINATOR].Timer.C:
			SetMonitoringState(MONITORING_HALT)
			smlog.Error(LOG_NETWORK, "COORDINATOR message not returned back within time limit. Sending it again...")
			SetWaiting(MSG_COORDINATOR, false)
			startElection()
			break
		}
	}
}
