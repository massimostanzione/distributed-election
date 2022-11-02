// smcommunication.go
// sm behavior is in statemachine.go
package bully

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/monitoring"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
)

func sendElection(msg *MsgElectionBully, dest *SMNode) {
	SafeRMI(MSG_ELECTION_BULLY, dest, msg, nil, nil)
}

func sendOk(msg *MsgOk, dest *SMNode) {
	SafeRMI(MSG_OK, dest, nil, msg, nil)
}

func sendCoord(msg *MsgCoordinator, dest *SMNode) {
	SafeRMI(MSG_COORDINATOR, dest, nil, nil, msg)
}

//TODO spostare?
func startElection() {
	SetMonitoringState(MONITORING_HALT)

	DirtyNetList = true

	CurState.Participant = true
	smlog.InfoU("aaaaaaaa") //TODO gestire quando rimane solo uno, succede anche altrove

	nodes := AskForNodesWithGreaterIds(CurState.NodeInfo.GetId())
	smlog.InfoU("ricevo: %s", nodes)
	for _, nextNode := range nodes {
		if nextNode.GetId() != CurState.NodeInfo.GetId() {
			smlog.Error(LOG_ELECTION, "invio a %s", nextNode.GetId())
			go sendElection(NewElectionBullyMsg(), nextNode)
		}
	}
	ElectionTimer.Reset(time.Duration(Cfg.ELECTION_ESPIRY+Cfg.ELECTION_ESPIRY_TOLERANCE) * time.Millisecond)

}
