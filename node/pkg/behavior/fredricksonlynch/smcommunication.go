package fredricksonlynch

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
)

func sendElection(msg *MsgElection, dest *SMNode) bool {
	smlog.Trace(LOG_ELECTION, "Sending ELECTION message to the next node")
	rmiErr := SafeRMI(MSG_ELECTION, dest, true, msg, nil)
	return rmiErr
}

func sendCoord(msg *MsgCoordinator, dest *SMNode) {
	smlog.Trace(LOG_ELECTION, "Sending COORDINATOR message to the next node")
	SafeRMI(MSG_COORDINATOR, dest, true, nil, msg)
}

func vote(inp *MsgElection) *MsgElection {
	var ret *MsgElection
	if !RedudantElectionCheck(CurState.NodeInfo.GetId(), inp) {
		smlog.Trace(LOG_ELECTION, "- voting...")
		//nextNode := AskForNodeInfo(State.NodeInfo.GetId()+1, true)
		ret = inp.AddVoter(CurState.NodeInfo.GetId())
		//spostare fuori il seguente
		//go sendElection(inp, nextNode)
		//sendCompiledMessage(State.NodeInfo.GetId(), MSG_ELECTION, nextNode, inp.GetStarter(), inp.GetIds())
		//setState(STATE_ELECTION_VOTER)
	} else {
		// this case should be already managed in SafeRMI
		// if this else is reached, something is wrong
		smlog.Fatal(LOG_ELECTION, "already voted! something is wrong")
	}
	return ret
}

func elect(candidates []int32) int32 {
	max := candidates[0]
	for _, val := range candidates {
		if val > max {
			max = val
		}
	}
	smlog.Trace(LOG_ELECTION, "Elected node no. %d", max)
	return max
}
