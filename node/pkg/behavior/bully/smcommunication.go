// smcommunication.go
// sm behavior is in statemachine.go
package bully

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"time"
)

func sendElection(msg *MsgElection, dest *SMNode) {
	/*abortIfFailedCoord := false
	if dest.GetId() == CoordId {
		abortIfFailedCoord = true
	}*/

	rmiErr := SafeRMI(MSG_ELECTION, dest, false, msg, nil, nil)
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
	//abortIfFailedCoord := false
	rmiErr := SafeRMI(MSG_OK, dest, false, nil, msg, nil)
	smlog.InfoU("(rimuovere) rmiErr=%s", rmiErr)
}
func sendCoord(msg *MsgCoordinator, dest *SMNode) {
	/*	abortIfFailedCoord := false
		if dest.GetId() == CoordId {
			abortIfFailedCoord = true
		}*/
	rmiErr := SafeRMI(MSG_COORDINATOR, dest, false, nil, nil, msg)
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
	if !RedudantElectionCheck(State.NodeInfo.GetId(), inp) {
		smlog.Info(LOG_ELECTION, "- voting...")
		nextNode := AskForNodeInfo(State.NodeInfo.GetId()+1, true)
		inp.AddVoter(State.NodeInfo.GetId())
		sendElection(inp, nextNode)
		//sendCompiledMessage(State.NodeInfo.GetId(), MSG_ELECTION, nextNode, inp.GetStarter(), inp.GetIds())
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
		nextNode := AskForNodeInfo(State.NodeInfo.GetId()+1, true)
		go sendCoord(coordinatorMsg, nextNode)
		//go send(MSG_COORDINATOR, starter, nextNode, )
	}
	if CoordId == State.NodeInfo.GetId() {
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
	DirtyNetList = true

	State.Participant = true
	smlog.InfoU("aaaaaaaa") //TODO gestire quando rimane solo uno, succede anche altrove

	//nextNode := AskForNodeInfo(State.NodeInfo.GetId()+1, true)
	//Oss. secondo parametro sempre false per fault tolerance,
	//     per rilevare eventuali nodi offline
	nodes := AskForNodesWithGreaterIds(State.NodeInfo.GetId())
	smlog.InfoU("ricevo: %s", nodes)
	for _, nextNode := range nodes {
		if nextNode.GetId() != State.NodeInfo.GetId() {
			smlog.Error(LOG_ELECTION, "invio a %s", nextNode.GetId())
			go sendElection(NewElectionMsg(), nextNode)
		}
	}
	ElectionTimer.Reset(time.Duration(Cfg.ELECTION_ESPIRY+Cfg.ELECTION_ESPIRY_TOLERANCE) * time.Millisecond)
	//setState(STATE_ELECTION)
	// se sono rimasto solo io non faccio nemmeno iniziare l'elezione, è inutile
	/*if nextNode.GetId() != State.NodeInfo.GetId() {
		//sendEmptyMessage(State.NodeInfo.GetId(), MSG_ELECTION, nextNode)
		//	starter = true //TODO ricordare di resettare
		setState(STATE_ELECTION)
	} else {
		smlog.InfoU("Sono rimasto solo io/2")
		setState(STATE_COORDINATOR)
	}*/
}
