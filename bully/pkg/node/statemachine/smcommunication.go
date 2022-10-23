// smcommunication.go
// sm behavior is in statemachine.go
package statemachine

import (
	. "bully/pkg/node/env"
	. "bully/pkg/node/net"

	. "bully/tools/smlog"
	smlog "bully/tools/smlog"
)

func sendElection(msg *MsgElection, dest *SMNode) {
	abortIfFailedCoord := false
	if dest.GetId() == CoordId {
		abortIfFailedCoord = true
	}

	rmiErr := SafeRMI(MSG_ELECTION, dest, !abortIfFailedCoord, msg, nil, nil)
	smlog.InfoU("(rimuovere) %s", rmiErr)
	// se il coordinatore è fallito prima che l'elezione terminasse, essa va abortita e va iniziata una nuova
	/*	if abortIfFailedCoord && rmiErr {
		smlog.Critical(LOG_ELECTION, "coordinatore andato mentre ancora c'era l'elezione, ne inizio una nuova...")

		startElection()
		setState(STATE_ELECTION_STARTER)
	}*/

}
func sendOk(msg *MsgOk, dest *SMNode) {
	// controllo AIFC?
	abortIfFailedCoord := false
	rmiErr := SafeRMI(MSG_OK, dest, !abortIfFailedCoord, nil, msg, nil)
	smlog.InfoU("(rimuovere) rmiErr=%s", rmiErr)
}
func sendCoord(msg *MsgCoordinator, dest *SMNode) {
	abortIfFailedCoord := false
	if dest.GetId() == CoordId {
		abortIfFailedCoord = true
	}
	rmiErr := SafeRMI(MSG_COORDINATOR, dest, !abortIfFailedCoord, nil, nil, msg)
	smlog.InfoU("(rimuovere) rmiErr=%s", rmiErr)
	/*
		// se il coordinatore è fallito prima che l'elezione terminasse, essa va abortita e va iniziata una nuova
		if abortIfFailedCoord && rmiErr {
			smlog.Critical(LOG_ELECTION, "coordinatore andato mentre ancora c'era l'elezione, ne inizio una nuova...")
			startElection()
			setState(STATE_ELECTION_STARTER)
		}*/
	//setState(STATE_ELECTION_STARTER)
}

/*
func vote(inp *MsgElection) {
	if !RedudantElectionCheck(Me.GetId(), inp) {
		smlog.Info(LOG_ELECTION, "- voting...")
		nextNode := AskForNodeInfo(Me.GetId()+1, true)
		inp.AddVoter(Me.GetId())
		sendElection(inp, nextNode)
		//sendCompiledMessage(Me.GetId(), MSG_ELECTION, nextNode, inp.GetStarter(), inp.GetIds())
		//setState(STATE_ELECTION_VOTER)
	} else {
		smlog.Fatal(LOG_ELECTION, "avevo già votato, TODO implementare abortElection")
	}
}
*/
/*
func endElection(coordinatorMsg *MsgCoordinator, forwardMsg bool) {
	// aggiorna coord
	CoordId = coordinatorMsg.GetCoordinator()
	smlog.Info(LOG_ELECTION, "new coordinator is ", CoordId)
	// vedi se C/NC
	if forwardMsg {
		nextNode := AskForNodeInfo(Me.GetId()+1, true)
		go sendCoord(coordinatorMsg, nextNode)
		//go send(MSG_COORDINATOR, starter, nextNode, )
	}
	if CoordId == Me.GetId() {
		setState(STATE_COORDINATOR)
	} else {
		setState(STATE_NON_COORDINATOR)
	}
}
*/ /*
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
*/
func startElection() {
	smlog.InfoU("aaaaaaaa") //TODO gestire quando rimane solo uno, succede anche altrove

	//nextNode := AskForNodeInfo(Me.GetId()+1, true)
	//Oss. secondo parametro sempre false per fault tolerance,
	//     per rilevare eventuali nodi offline
	nodes := AskForNodesWithGreaterIds(Me.GetId(), false)
	for _, nextNode := range nodes {
		if nextNode.GetId() != Me.GetId() {
			smlog.Error(LOG_ELECTION, "invio a %s", nextNode.GetId())
			go sendElection(NewElectionMsg(Me.GetId()), nextNode)
		}
	}
	//setState(STATE_ELECTION)
	// se sono rimasto solo io non faccio nemmeno iniziare l'elezione, è inutile
	/*if nextNode.GetId() != Me.GetId() {
		//sendEmptyMessage(Me.GetId(), MSG_ELECTION, nextNode)
		//	starter = true //TODO ricordare di resettare
		setState(STATE_ELECTION)
	} else {
		smlog.InfoU("Sono rimasto solo io/2")
		setState(STATE_COORDINATOR)
	}*/
}
