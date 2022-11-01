// Behavior specification for the bully algorithm.
package bully

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/monitoring"
	. "distributedelection/node/pkg/net"
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
	ElectionChannel = make(chan *MsgElection)
	OkChannel = make(chan *MsgOk)
	CoordChannel = make(chan *MsgCoordinator)
	smlog.Info(LOG_UNDEFINED, "Starting SM...")
	smlog.InfoU("Type CTRL+C to terminate")

	go run()
	Listen()
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
	State.NodeInfo = AskForJoining()
	go startElection()
	for {
		select {
		case <-ElectionTimer.C:
			SetMonitoringState(MONITORING_SEND)
			SetWatchdog(MSG_COORDINATOR, false)
			State.Participant = false
			ElectionTimer.Stop()
			smlog.Info(LOG_ELECTION, "*** I declare myself as the new coordinator ***")
			State.Coordinator = State.NodeInfo.GetId()
			for _, dest := range AskForAllNodes() {
				if dest.GetId() != State.NodeInfo.GetId() {
					go sendCoord(NewCoordinatorMsg(State.NodeInfo.GetId(), State.NodeInfo.GetId()), dest)
				}
			}
			break
		case inp := <-ElectionChannel:
			SetMonitoringState(MONITORING_HALT)
			//SetWatchdog(MSG_COORDINATOR, true)
			State.Participant = true
			smlog.InfoU("arrivato E")
			// other elections are occurring
			if State.NodeInfo.GetId() > inp.GetStarter() {
				sendOk(NewOkMsg(State.NodeInfo.GetId()), AskForNodeInfo(inp.GetStarter()))
				go startElection()
			}
			break
		case <-OkChannel:
			SetMonitoringState(MONITORING_HALT)
			SetWatchdog(MSG_COORDINATOR, true)
			State.Participant = true
			smlog.Info(LOG_ELECTION, "Someone is bullier than me...")
			ElectionTimer.Stop()
			break
		case inp := <-CoordChannel:
			SetMonitoringState(MONITORING_HALT)
			SetWatchdog(MSG_COORDINATOR, false)
			State.Coordinator = inp.GetCoordinator()
			smlog.InfoU("arrivato C")
			if inp.GetCoordinator() != State.NodeInfo.GetId() {
				smlog.InfoU("nuovo coord è %d", inp.GetCoordinator())
				SetMonitoringState(MONITORING_LISTEN)
				State.Participant = false
			} else {
				smlog.Fatal(LOG_UNDEFINED, "unreachable")
			}
			break
		case <-MonitoringChannel:
			SetMonitoringState(MONITORING_HALT)
			go startElection()
			break
		case <-Watchdogs[MSG_COORDINATOR].Timer.C:
			SetMonitoringState(MONITORING_HALT)
			smlog.Error(LOG_NETWORK, "COORDINATOR message not received within time limit after received OK.")
			SetWatchdog(MSG_COORDINATOR, false)
			go startElection()
			break
		}
	}
}
