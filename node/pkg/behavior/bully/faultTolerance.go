package bully

import (
	"context"
	pb "distributedelection/node/pb"

	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/api"
	. "distributedelection/tools/formatting"
	"time"

	//"bully/node/pkg/statemachine"

	//. "bully/node/pkg/net"

	//"bully/node/pkg/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	"google.golang.org/grpc"
)

func SafeRMI(tipo MsgType, dest *SMNode, tryNextWhenFailed bool, elezione *MsgElection, okMsg *MsgOk, coord *MsgCoordinator) (failedNodeExistence bool) { //opt ...interface{}) {

	attempts := 0
	nextNode := dest
	prossimoId := nextNode.GetId()
	prossimoAddr := nextNode.GetFullAddr()
	failedNodeExistence = false
	//	smlog.Printf("il nodo NON sono io, quindi provo a contattarlo")
	// L'ALTRO NODO FUNGE DA SERVER NEI MIEI CONFRONTI
	var errq error
	var starter int32
	time.Sleep(time.Duration(GenerateDelay()) * time.Millisecond)
	for {
		starter = -1
		errq = nil
		connN := ConnectToNode(prossimoAddr)
		defer connN.Close()
		// New server instance and service registering
		nodoServer := grpc.NewServer()
		pb.RegisterDistrElectNodeServer(nodoServer, &DGnode{})
		csN := pb.NewDistrElectNodeClient(connN)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.RESPONSE_TIME_LIMIT)*time.Millisecond)
		//	locCtx = ctx
		defer cancel()
		ok := false
		if !tryNextWhenFailed {
			ok = true
		}
		attempts++
		switch tipo {
		case MSG_ELECTION:
			netMsg := ToNetElectionMsg(elezione)
			//smlog.Println("\033[42m\033[1;30mSENDING ELECTION [", elezione, "] to ", prossimoAddr, "\033[0m")
			smlog.Warn(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"SENDING ELECT %v to %s"+ColorReset, netMsg, prossimoAddr)
			_, errq = csN.ForwardElection(ctx, netMsg)
			break
		case MSG_OK:
			netMsg := ToNetOkMsg(okMsg)
			//smlog.Println("\033[42m\033[1;30mSENDING COORDINATOR [", coord, "] to ", prossimoAddr, "\033[0m")
			smlog.Warn(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING OK %v to %s"+ColorReset, netMsg, prossimoAddr)
			//	log.Printf("\033[42m\033[1;30mSENDING COORD %[1]v to %[2]s \033[0m", coord, prossimoAddr)
			_, errq = csN.ForwardOk(ctx, netMsg)
			break
		case MSG_COORDINATOR:
			netMsg := ToNetCoordinatorMsg(coord)
			//smlog.Println("\033[42m\033[1;30mSENDING COORDINATOR [", coord, "] to ", prossimoAddr, "\033[0m")
			smlog.Warn(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING COORD %v to %s"+ColorReset, netMsg, prossimoAddr)
			//	log.Printf("\033[42m\033[1;30mSENDING COORD %[1]v to %[2]s \033[0m", coord, prossimoAddr)
			_, errq = csN.ForwardCoordinator(ctx, netMsg)
			break

		default:
			break
		}
		//_, errq := csN.InoltraElezione(ctx, elezione)
		if errq != nil {
			if attempts != Cfg.RMI_RETRY_TOLERANCE {
				smlog.Warn(LOG_UNDEFINED, "tentativo n. %d non andato a buon fine per %v, riprovo... %s", attempts, prossimoAddr, errq)

			} else {
				if (tipo == MSG_ELECTION || tipo == MSG_COORDINATOR) && prossimoId == starter {
					// lo starter è fallito, smetto di far circolare una seconda volta il messaggio
					// ci penserà quando tornerà in piedi a far ripartire elezione
					// così tolgo messaggi inutili dalla rete
					smlog.Critical(LOG_NETWORK, "\n\n\n\n\n\n\n\n\nlo starter è failed, smetto di far circolare")
					ok = true
					tryNextWhenFailed = false
				}
				//DeclareNodeState(nextNode, false)
				failedNodeExistence = true
				smlog.Error(LOG_UNDEFINED, "impossibile chiamare RMI su %v:\n %v", prossimoAddr, errq)

				if tryNextWhenFailed {
					smlog.Error(LOG_UNDEFINED, "provo col prossimo...")
					//nextNode = AskForNodeInfo(nextNode.GetId()+1, true)
					nextNode = AskForNodeInfo(nextNode.GetId() + 1)
					smlog.InfoU("Prossimo: %v presso %v", nextNode.GetId(), nextNode.GetFullAddr())
					prossimoId = nextNode.GetId()
					prossimoAddr = nextNode.GetFullAddr()
					if prossimoAddr == State.NodeInfo.GetFullAddr() {
						smlog.InfoU("Sono rimasto solo io")
						ok = true
					}
				}
				//break
			}
		} else {
			smlog.InfoU("Ok, andata bene! posso uscire")
			ok = true
		}
		if ok {
			smlog.InfoU("esco...")
			break
		}
		attempts = 0
		smlog.InfoU("torno nel ciclo col prossimo nodo...")
	}
	return failedNodeExistence
}
