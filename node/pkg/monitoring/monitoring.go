// Monitoring tool
// Use it in other .go files via the following API:
// - invoke SetMonitoringState(monitoringState) to change the monitoring state
// - use MonitoringChannel to receive messages when the monitoring ticker has expired
// Ticker parameters are described into Cfg structure, env package.
package monitoring

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/misc"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	"time"
)

// Monitoring tool state. Available states are described below
type monitoringState uint8

const (
	// Don't send, nor listen to any heartbeat messages
	MONITORING_HALT monitoringState = iota

	// Send heartbeat messages
	MONITORING_SEND

	// Listen for heartbeat messages,
	// triggering a message into MonitoringChannel
	// when the time has expired
	MONITORING_LISTEN
)

var timer *time.Ticker
var stateChan chan (monitoringState)
var curMonitoringState monitoringState

// Receive booleans into this channel when the monitoring ticker has expired
var MonitoringChannel chan (bool)

func initialize() {
	MonitoringChannel = make(chan bool)
	Heartbeat = make(chan *MsgHeartbeat)
	stateChan = make(chan monitoringState)
	curMonitoringState = MONITORING_HALT
	timer = TickerFictitiousInit(timer)
	go monitoring()
}

// Set monitoring state to newState
func SetMonitoringState(newState monitoringState) {
	if MonitoringChannel == nil {
		initialize()
	}
	stateChan <- newState
}

func monitoring() {
	for {
		select {
		case curMonitoringState = <-stateChan:
			smlog.Debug(LOG_MONITORING, "new state: %d", curMonitoringState)
			updateTimer()
			break
		case <-Heartbeat:
			smlog.Trace(LOG_MONITORING, "HB received")
			if curMonitoringState == MONITORING_LISTEN {
				acknowledgeHb()
			}
			break
		case <-timer.C:
			smlog.Trace(LOG_MONITORING, "ticker expired")
			handleTickerExpiry()
			break
		}
	}
}

func updateTimer() {
	switch curMonitoringState {
	case MONITORING_HALT:
		timer.Stop()
		break
	case MONITORING_LISTEN:
		timer.Reset(time.Duration(Cfg.MONITORING_TIMEOUT+Cfg.MONITORING_TOLERANCE) * time.Millisecond)
		break
	case MONITORING_SEND:
		timer.Reset(time.Duration(Cfg.MONITORING_TIMEOUT) * time.Millisecond)
		break
	}
}

func handleTickerExpiry() {
	switch curMonitoringState {
	case MONITORING_HALT:
		timer.Stop()
		break
	case MONITORING_LISTEN:
		timer.Stop()
		curMonitoringState = MONITORING_HALT
		smlog.Critical(LOG_MONITORING, "Monitoring timer expired!")
		MonitoringChannel <- true
		break
	case MONITORING_SEND:
		sendHb()
		break
	}
}
func acknowledgeHb() {
	timer.Reset(time.Duration(Cfg.MONITORING_TIMEOUT+Cfg.MONITORING_TOLERANCE) * time.Millisecond)
}

func sendHb() {
	timer.Reset(time.Duration(Cfg.MONITORING_TIMEOUT) * time.Millisecond)
	allNodesList := AskForAllNodesList()
	hbMsg := &MsgHeartbeat{Id: CurState.NodeInfo.GetId()}
	smlog.Debug(LOG_UNDEFINED, "SuccessfulHB = %v", SuccessfulHB)
	if SuccessfulHB == 0 {
		SuccessfulHB = -1
	}
	SuccessfulHB = len(allNodesList) - 1
	for _, node := range allNodesList {
		if node.GetFullAddr() != CurState.NodeInfo.GetFullAddr() {
			smlog.Info(LOG_MONITORING, "Sending HB to node %d, at %s", node.GetId(), node.GetFullAddr())
			time.Sleep(time.Duration(GenerateDelay()) * time.Millisecond)
			go SafeHB(ToNetHeartbeatMsg(hbMsg), node)
		}
	}
}
