package fredricksonlynch

import (
	. "distributedelection/node/env"
	. "distributedelection/node/structure/net"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
)

func sendElection(msg *MsgElectionFL, dest *SMNode) bool {
	smlog.Trace(LOG_ELECTION, "Sending ELECTION message to the next node")
	rmiErr := SafeRMI_Ring(MSG_ELECTION_FL, dest, true, msg, nil)
	return rmiErr
}

func sendCoord(msg *MsgCoordinator, dest *SMNode) {
	smlog.Trace(LOG_ELECTION, "Sending COORDINATOR message to the next node")
	SafeRMI_Ring(MSG_COORDINATOR, dest, true, nil, msg)
}
