// smcommunication.go
// sm behavior is in statemachine.go
package statemachine

import (
	. "fredricksonLynch/pkg/node/env"
	. "fredricksonLynch/pkg/node/net"
	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"
)

func sendElection(msg *MsgElection, dest *SMNode) bool {
	/*abortIfFailedCoord := false
	if dest.GetId() == CoordId {
		abortIfFailedCoord = true
	}
	*/
	rmiErr := SafeRMI(MSG_ELECTION, dest, true, msg, nil, nil)

	// se il coordinatore è fallito prima che l'elezione terminasse, essa va abortita e va iniziata una nuova
	/*	if abortIfFailedCoord && rmiErr {
		smlog.Critical(LOG_ELECTION, "coordinatore andato mentre ancora c'era l'elezione, ne inizio una nuova...")

		startElection()
		setState(STATE_ELECTION_STARTER)
	}*/
	return rmiErr

}

func sendCoord(msg *MsgCoordinator, dest *SMNode) {
	/*abortIfFailedCoord := false
	if dest.GetId() == CoordId {
		abortIfFailedCoord = true
	}*/
	//rmiErr := SafeRMI(MSG_COORDINATOR, dest, false, nil, msg, nil)
	SafeRMI(MSG_COORDINATOR, dest, true, nil, msg, nil)

	// se il coordinatore è fallito prima che l'elezione terminasse, essa va abortita e va iniziata una nuova
	/*if abortIfFailedCoord && rmiErr {
		smlog.Critical(LOG_ELECTION, "coordinatore andato mentre ancora c'era l'elezione, ne inizio una nuova...")
		startElection()
		setState(STATE_ELECTION_STARTER)
	}*/
	//setState(STATE_ELECTION_STARTER)
}

func vote(inp *MsgElection) {
	if !RedudantElectionCheck(Me.GetId(), inp) {
		smlog.Info(LOG_ELECTION, "- voting...")
		nextNode := AskForNodeInfo(Me.GetId()+1, true)
		inp.AddVoter(Me.GetId())
		//spostare fuori il seguente
		go sendElection(inp, nextNode)
		//sendCompiledMessage(Me.GetId(), MSG_ELECTION, nextNode, inp.GetStarter(), inp.GetIds())
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
		smlog.Fatal(LOG_ELECTION, "avevo già votato, TODO implementare abortElection")
	}
}

func endElection(coordinatorMsg *MsgCoordinator, forwardMsg bool) {
	// aggiorna coord
	CoordId = coordinatorMsg.GetCoordinator()
	smlog.Info(LOG_ELECTION, "new coordinator is ", CoordId)
	// vedi se C/NC
	if forwardMsg {
		//nextNode := AskForNodeInfo(Me.GetId()+1, true)
		go sendCoord(coordinatorMsg, NextNode)
		//go send(MSG_COORDINATOR, starter, nextNode, )
	}
	if CoordId == Me.GetId() {
		setState(STATE_COORDINATOR)
	} else {
		setState(STATE_NON_COORDINATOR)
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
