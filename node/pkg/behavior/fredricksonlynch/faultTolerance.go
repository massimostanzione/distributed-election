package fredricksonlynch

import (
	"context"
	pb "distributedelection/node/pb"

	//b "distributedelection/node/pkg/behavior/bully"
	//. "distributedelection/node/pkg/behavior/fredricksonlynch"
	. "distributedelection/node/pkg/env"
	. "distributedelection/tools/formatting"

	//	"math/rand"
	"time"

	//"distributedelection/node/pkg/statemachine"

	. "distributedelection/node/pkg/net"

	//"distributedelection/node/pkg/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	"google.golang.org/grpc"
)

func RedudantElectionCheck(voter int32, electionMsg *MsgElection) bool {
	for _, i := range electionMsg.GetVoters() {
		if i == voter {
			return true
		}
	}
	return false
}
func SafeRMI(tipo MsgType, dest *SMNode, tryNextWhenFailed bool, elezione *MsgElection, coord *MsgCoordinator) (failedNodeExistence bool) { //opt ...interface{}) {

	// update local cache the first time a sequential message is sent
	if NextNode.GetId() == 0 {
		requested := AskForNodeInfo(Me.GetId() + 1)
		if requested.GetId() != 1 {
			NextNode = requested
			dest = requested
			smlog.Debug(LOG_UNDEFINED, "NextNode initialized.", NextNode)
		}
	}

	attempts := 0
	nextNode := dest
	prossimoId := nextNode.GetId()
	prossimoAddr := nextNode.GetFullAddr()
	failedNodeExistence = false
	//TODO funzione unica "connetti a..." con paramtero addr, qui e altrove
	//	smlog.Printf("il nodo NON sono io, quindi provo a contattarlo")
	// L'ALTRO NODO FUNGE DA SERVER NEI MIEI CONFRONTI
	var errq error
	var starter int32
	time.Sleep(time.Duration(GenerateDelay()) * time.Millisecond)
	for {
		starter = -1
		errq = nil
		connN := ConnectToNode(prossimoAddr)
		//defer connN.Close()
		// New server instance and service registering
		nodoServer := grpc.NewServer()
		pb.RegisterDistGrepServer(nodoServer, &DGnode{})
		csN := pb.NewDistGrepClient(connN)
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
			starter = elezione.GetStarter()
			netMsg := ToNetElectionMsg(elezione)

			//smlog.Println("\033[42m\033[1;30mSENDING ELECTION [", elezione, "] to ", prossimoAddr, "\033[0m")
			smlog.Info(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"SENDING ELECT %v to %s"+ColorReset, netMsg, prossimoAddr)
			_, errq = csN.ForwardElection(ctx, netMsg)
			break
		case MSG_COORDINATOR:
			starter = coord.GetStarter()
			netMsg := ToNetCoordinatorMsg(coord)
			//smlog.Println("\033[42m\033[1;30mSENDING COORDINATOR [", coord, "] to ", prossimoAddr, "\033[0m")
			smlog.Info(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING COORD %v to %s"+ColorReset, netMsg, prossimoAddr)
			//	log.Printf("\033[42m\033[1;30mSENDING COORD %[1]v to %[2]s \033[0m", coord, prossimoAddr)
			_, errq = csN.ForwardCoordinator(ctx, netMsg)
			break
		default:
			break
		}
		//_, errq := csN.InoltraElezione(ctx, elezione)
		if errq != nil {
			if attempts != Cfg.RMI_RETRY_TOLERANCE {
				smlog.Warn(LOG_NETWORK, "Failed attempt n. %d to contact %v, trying again...", attempts, prossimoAddr)
				smlog.Debug(LOG_NETWORK, "(%s)", errq)

			} else {
				if (tipo == MSG_ELECTION || tipo == MSG_COORDINATOR) && prossimoId == starter {
					// lo starter è fallito, smetto di far circolare una seconda volta il messaggio
					// ci penserà quando tornerà in piedi a far ripartire elezione
					// così tolgo messaggi inutili dalla rete

					smlog.Error(LOG_NETWORK, "Election starter failed! Stopping forwarding...")
					smlog.Debug(LOG_NETWORK, "(%s)", errq)
					ok = true
					tryNextWhenFailed = false
				}
				if tipo == MSG_HEARTBEAT {
					SuccessfulHB--
				}
				//DeclareNodeState(nextNode, false)
				failedNodeExistence = true
				smlog.Error(LOG_NETWORK, "Could not invoke RMI on %v", prossimoAddr)
				smlog.Debug(LOG_NETWORK, "(%s)", errq)

				if tryNextWhenFailed {
					//nextNode = AskForNodeInfo(nextNode.GetId()+1, true)
					nextNode = AskForNodeInfo(nextNode.GetId() + 1)
					smlog.Info(LOG_NETWORK, "Trying next node: %v@%v", nextNode.GetId(), nextNode.GetFullAddr())
					prossimoId = nextNode.GetId()
					prossimoAddr = nextNode.GetFullAddr()
					if prossimoAddr == Me.GetFullAddr() {
						smlog.InfoU("Sono rimasto solo io")
						ok = true
					}
				}
				//break
			}
		} else {
			smlog.Trace(LOG_NETWORK, "RMI invoked correctly, exiting from SafeRMI...")
			ok = true
		}
		if ok {
			//smlog.Critical(LOG_NETWORK, "esco...")
			break
		}
		attempts = 0
		///smlog.Critical(LOG_NETWORK, "torno nel ciclo col prossimo nodo...")
	}
	return failedNodeExistence
}
