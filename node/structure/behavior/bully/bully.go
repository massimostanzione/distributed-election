// Behavior specification for the bully algorithm.
package bully

import (
	. "distributedelection/node/env"
	. "distributedelection/node/structure/net"
	. "distributedelection/node/tools/monitoring"
	. "distributedelection/tools/misc"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
)

var ElectionTimer *time.Timer

// exposed to runner.go
func Run() {
	initWatchdogs()
	ElectionTimer = TimerFictitiousInit(ElectionTimer)
	ElectionChannel_bully = make(chan *MsgElectionBully)
	OkChannel = make(chan *MsgOk)
	CoordChIn = make(chan *MsgCoordinator)

	go run()
	ListenToIncomingRMI()
}

func initWatchdogs() {
	// watchdog is actived after an OK message:
	// a COORDINATOR message must be received before Cfg.IDLE_WAIT_LIMIT,
	// elsewhere we know that the election starter has failed.
	Watchdogs = map[MsgType]*Watchdog{
		MSG_COORDINATOR: &Watchdog{
			Waiting: false,
			Timer:   time.NewTimer(time.Duration(Cfg.IDLE_WAIT_LIMIT) * time.Second),
		},
	}
	Watchdogs[MSG_COORDINATOR].Timer.Stop()
}

func run() {
	CurState.NodeInfo = AskForJoining()
	go startElection()
	for {
		select {
		case <-ElectionTimer.C:
			smlog.Info(LOG_ELECTION, "*** I declare myself as the new coordinator ***")
			SetMonitoringState(MONITORING_SEND)
			SetWatchdog(MSG_COORDINATOR, false)
			ElectionTimer.Stop()
			CurState.Coordinator = CurState.NodeInfo.GetId()
			CurState.Participant = false
			CurState.DirtyNetCache = true
			for _, dest := range AskForAllNodes() {
				if dest.GetId() != CurState.NodeInfo.GetId() {
					go sendCoord(NewCoordinatorMsg(CurState.NodeInfo.GetId(), CurState.NodeInfo.GetId()), dest)
				}
			}
			break
		case inp := <-ElectionChannel_bully:
			smlog.Debug(LOG_ELECTION, "Handling ELECTION message")
			SetMonitoringState(MONITORING_HALT)
			CurState.Participant = true
			// other elections are occurring
			if CurState.NodeInfo.GetId() > inp.GetStarter() {
				sendOk(NewOkMsg(CurState.NodeInfo.GetId()), AskForNodeInfo(inp.GetStarter()))
				go startElection()
			}
			break
		case <-OkChannel:
			smlog.Debug(LOG_ELECTION, "Handling OK message - someone is bullier than me...")
			SetMonitoringState(MONITORING_HALT)
			SetWatchdog(MSG_COORDINATOR, true)
			CurState.Participant = true
			ElectionTimer.Stop()
			break
		case inp := <-CoordChIn:
			smlog.Debug(LOG_ELECTION, "Handling COORDINATOR message")
			SetMonitoringState(MONITORING_HALT)
			SetWatchdog(MSG_COORDINATOR, false)
			CurState.Coordinator = inp.GetCoordinator()
			if inp.GetCoordinator() != CurState.NodeInfo.GetId() {
				smlog.Info(LOG_ELECTION, "New coordinator: %d", inp.GetCoordinator())
				SetMonitoringState(MONITORING_LISTEN)
				CurState.Participant = false
			} else {
				smlog.Fatal(LOG_UNDEFINED, "unreachable")
			}
			break
		case <-MonitoringChannel:
			smlog.Critical(LOG_ELECTION, "Coordinator failed!")
			SetMonitoringState(MONITORING_HALT)
			go startElection()
			break
		case <-Watchdogs[MSG_COORDINATOR].Timer.C:
			SetMonitoringState(MONITORING_HALT)
			SetWatchdog(MSG_COORDINATOR, false)
			smlog.Error(LOG_NETWORK, "COORDINATOR message not received within time limit after received OK.")
			go startElection()
			break
		}
	}
}

func startElection() {
	SetMonitoringState(MONITORING_HALT)
	CurState.DirtyNetCache = true
	CurState.Participant = true
	nodes := AskForNodesWithGreaterIds(CurState.NodeInfo.GetId())
	for _, nextNode := range nodes {
		if nextNode.GetId() != CurState.NodeInfo.GetId() {
			go sendElection(NewElectionBullyMsg(), nextNode)
		}
	}
	ElectionTimer.Reset(time.Duration(Cfg.ELECTION_ESPIRY+Cfg.ELECTION_ESPIRY_TOLERANCE) * time.Millisecond)
}
