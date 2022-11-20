// smcommunication.go
// sm behavior is in statemachine.go
package bully

import (
	. "distributedelection/node/env"
	. "distributedelection/node/structure/net"
	. "distributedelection/tools/api"
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
