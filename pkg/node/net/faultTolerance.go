package net

import (
	"context"
	pb "fredricksonLynch/pb/node"
	. "fredricksonLynch/pkg/node/env"

	//"fredricksonLynch/pkg/node/statemachine"

	//. "fredricksonLynch/pkg/node/net"

	//"fredricksonLynch/pkg/node/net"
	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"
	"log"
	"time"

	"google.golang.org/grpc"
)

// TODO list:
// - [ok] gestione failure starter durante elezione (ma anche coordinamento)
// - timeout elezione se rimango dopo timeout senza sapere nulla
// - abortElection, da chiamare in vote()
// - [se c'è tempo] pause/resume per ogni processo, handling CTRL-Z
// --- [ok] ciò è correlato all'ascolto degli HB: se mi arrivano da chi non mi aspetto c'è problema => avvio elezione

func RedudantElectionCheck(voter int32, electionMsg *MsgElection) bool {
	for _, i := range electionMsg.GetVoters() {
		if i == voter {
			return true
		}
	}
	return false
}

func SafeRMI(tipo MsgType, dest *SMNode, tryNextWhenFailed bool, elezione *MsgElection, coord *MsgCoordinator, hb *pb.Heartbeat) (failedNodeExistence bool) { //opt ...interface{}) {
	//TODO gestione delay con parametro (separata da SM)
	for Pause {
	}
	//FIXME da risolvere: scenario in cui provo a far girare i messaggi ma lo starter fallisce prima
	// che gli ritorni
	// esempio: se il mio successivo è lo starter e non è in piedi,
	// faccio girare per un certo numero di volte, oppure ABORT
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()
	attempts := 0
	nextNode := dest
	//	prossimoId := nextNode.GetId()
	prossimoAddr := nextNode.GetFullAddr()
	failedNodeExistence = false
	for {
		//TODO funzione unica "connetti a..." con paramtero addr, qui e altrove
		//	smlog.Printf("il nodo NON sono io, quindi provo a contattarlo")
		// L'ALTRO NODO FUNGE DA SERVER NEI MIEI CONFRONTI
		connN, errN := grpc.Dial(prossimoAddr, grpc.WithInsecure())
		if errN != nil {
			log.Fatalf("Error while contacting server (NODO) on %v:\n %v", prossimoAddr, errN)
		}
		//defer connN.Close()
		// New server instance and service registering
		nodoServer := grpc.NewServer()
		pb.RegisterDistGrepServer(nodoServer, &DGnode{})
		csN := pb.NewDistGrepClient(connN)
		for {
			attempts++
			errq := error(nil)
			switch tipo {
			case MSG_ELECTION:
				netMsg := ToNetElectionMsg(elezione)
				//smlog.Println("\033[42m\033[1;30mSENDING ELECTION [", elezione, "] to ", prossimoAddr, "\033[0m")
				smlog.Warn(LOG_MSG_SENT, "\033[42m\033[1;30mSENDING ELECT %v to %s \033[0m", netMsg, prossimoAddr)
				_, errq = csN.ForwardElection(ctx, netMsg)
				break
			case MSG_COORDINATOR:
				netMsg := ToNetCoordinatorMsg(coord)
				//smlog.Println("\033[42m\033[1;30mSENDING COORDINATOR [", coord, "] to ", prossimoAddr, "\033[0m")
				smlog.Warn(LOG_MSG_RECV, "\033[42m\033[1;30mSENDING COORD %v to %s \033[0m", netMsg, prossimoAddr)
				//	log.Printf("\033[42m\033[1;30mSENDING COORD %[1]v to %[2]s \033[0m", coord, prossimoAddr)
				_, errq = csN.ForwardCoordinator(ctx, netMsg)
				break
			case MSG_HEARTBEAT:
				_, errq = csN.SendHeartBeat(ctx, hb)
				break
			default:
				break
			}
			//_, errq := csN.InoltraElezione(ctx, elezione)
			if errq != nil {
				if attempts != RMI_RETRY_TOLERANCE {
					smlog.Warn(LOG_UNDEFINED, "tentativo n. %d non andato a buon fine per %v, riprovo...", attempts, prossimoAddr)

				} else {
					DeclareNodeState(nextNode, false)
					failedNodeExistence = true
					smlog.Error(LOG_UNDEFINED, "impossibile chiamare RMI su %v:\n %v", prossimoAddr, errq)

					if tryNextWhenFailed {
						smlog.Error(LOG_UNDEFINED, "provo col prossimo che risulta in piedi...")
						nextNode = AskForNodeInfo(nextNode.GetId()+1, true)
						//						prossimoId = nextNode.GetId()
						prossimoAddr = nextNode.GetFullAddr()
					}
					break
				}
			} else {
				break
			}
		}
		if attempts == 1 || !tryNextWhenFailed {
			break
		}
		attempts = 0
	}
	return failedNodeExistence
}
