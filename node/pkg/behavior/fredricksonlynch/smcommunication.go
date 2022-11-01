// smcommunication.go
// sm behavior is in statemachine.go
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
	if !RedudantElectionCheck(State.NodeInfo.GetId(), inp) {
		smlog.Trace(LOG_ELECTION, "- voting...")
		//nextNode := AskForNodeInfo(State.NodeInfo.GetId()+1, true)
		ret = inp.AddVoter(State.NodeInfo.GetId())
		//spostare fuori il seguente
		//go sendElection(inp, nextNode)
		//sendCompiledMessage(State.NodeInfo.GetId(), MSG_ELECTION, nextNode, inp.GetStarter(), inp.GetIds())
		//setState(STATE_ELECTION_VOTER)
	} else {
		// due opzioni:
		// 1 - me ne accoro a posteriori (che è questa): se lo starter ha fallito me ne
		//     accorgo perché mi arriva un ELECTION in cui avevo già votato,
		//     e ciò significa che il messaggio è andato OLTRE il suo giro
		// 2 - me ne accorgo prima di inviarlo: se il destinatario è lo starter e non risponde
		//     allora abortisco l'elezione, tanto se torna ha il suo timer lungo
		// ==> DA PREFERIRSI LA N. 2,
		//     perché ciò vale anche per COORD, e lì non ho il controllo per vedere se ho votato o meno
		smlog.Fatal(LOG_ELECTION, "already voted! something is wrong")
	}
	return ret
}

func endElection(coordinatorMsg *MsgCoordinator, forwardMsg bool) {
	// aggiorna coord
	State.Coordinator = coordinatorMsg.GetCoordinator()
	smlog.Info(LOG_ELECTION, "new coordinator is ", State.Coordinator)
	// vedi se C/NC
	if forwardMsg {
		//nextNode := AskForNodeInfo(State.NodeInfo.GetId()+1, true)
		go sendCoord(coordinatorMsg, NextNode)
		//go send(MSG_COORDINATOR, starter, nextNode, )
	}
}

func elect(candidates []int32) int32 {
	max := candidates[0]
	for _, val := range candidates {
		if val > max {
			max = val
		}
	}
	//	smlog.Println("PROCLAMO ELETTO IL NUMERO", max)
	return max
}
