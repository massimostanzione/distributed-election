package fredricksonlynch

import (
	. "distributedelection/node/env"
	. "distributedelection/node/structure/net"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
)

func sendElection(msg *MsgElectionFL, dest *SMNode) {
	MsgOrderOut <- MSG_ELECTION_FL
	ElectChOut <- msg
}

func sendCoord(msg *MsgCoordinator, dest *SMNode) {
	MsgOrderOut <- MSG_COORDINATOR
	CoordChOut <- msg
}

// Process outbound traffic while:
// - guaranteeing FIFO
// - avoiding stalls
func outboundQueue() {
	for {
		out := <-MsgOrderOut
		switch out {
		case MSG_ELECTION_FL:
			msg := <-ElectChOut
			smlog.Trace(LOG_ELECTION, "Sending ELECTION message to the next node")
			SafeRMI_Ring(MSG_ELECTION_FL, NextNode, msg, nil)
			break
		case MSG_COORDINATOR:
			msg := <-CoordChOut
			smlog.Trace(LOG_ELECTION, "Sending COORDINATOR message to the next node")
			SafeRMI_Ring(MSG_COORDINATOR, NextNode, nil, msg)
			break
		}
	}
}
