package net

import (
	"context"
	pb "fredricksonlynch/node/pb"
	. "fredricksonlynch/node/pkg/env"
	. "fredricksonlynch/tools/formatting"
	"math/rand"
	"time"

	//"fredricksonlynch/node/pkg/statemachine"

	//. "fredricksonlynch/node/pkg/net"

	//"fredricksonlynch/node/pkg/net"
	. "fredricksonlynch/tools/smlog"
	smlog "fredricksonlynch/tools/smlog"

	"google.golang.org/grpc"
)

func ToNCL(input string) NetCongestionLevel {
	switch input {
	case "ABSENT":
		return NCL_ABSENT
	case "LIGHT":
		return NCL_LIGHT
	case "MEDIUM":
		return NCL_MEDIUM
	case "SEVERE":
		return NCL_SEVERE
	case "CUSTOM":
		return NCL_CUSTOM
	default:
		return NCL_ABSENT
	}
}
func generateDelay() int32 {
	var min float32
	var max float32
	switch Cfg.NCL_CONGESTION_LEVEL {
	case NCL_ABSENT:
		min = 0
		max = 0
	case NCL_LIGHT:
		min = 0
		max = .2 * Cfg.HB_TIMEOUT
	case NCL_MEDIUM:
		min = .3 * Cfg.HB_TIMEOUT
		max = .5 * Cfg.HB_TIMEOUT
	case NCL_SEVERE:
		min = .5 * Cfg.HB_TIMEOUT
		max = 1.5 * Cfg.HB_TIMEOUT
	case NCL_CUSTOM:
		min = Cfg.NCL_CUSTOM_DELAY_MIN
		max = Cfg.NCL_CUSTOM_DELAY_MAX

	}
	ret := (rand.Float32() * (max - min)) + min
	return int32(ret)
}

func RedudantElectionCheck(voter int32, electionMsg *MsgElection) bool {
	for _, i := range electionMsg.GetVoters() {
		if i == voter {
			return true
		}
	}
	return false
}

func SafeRMI(tipo MsgType, dest *SMNode, tryNextWhenFailed bool, elezione *MsgElection, coord *MsgCoordinator, hb *MsgHeartbeat) (failedNodeExistence bool) { //opt ...interface{}) {
	for Pause {
	}
	// update local cache the first time a sequential message is sent
	if NextNode.GetId() == 0 && tipo != MSG_HEARTBEAT {
		requested := AskForNodeInfo(Me.GetId() + 1)
		if requested.GetId() != 1 {
			NextNode = requested
			dest = requested
			smlog.Critical(LOG_UNDEFINED, "aggiornato NextNode", NextNode)
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
	time.Sleep(time.Duration(generateDelay()) * time.Millisecond)
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
			smlog.Warn(LOG_MSG_SENT, ColorBlkBckgrGreen+BoldBlack+"SENDING ELECT %v to %s"+ColorReset, netMsg, prossimoAddr)
			_, errq = csN.ForwardElection(ctx, netMsg)
			break
		case MSG_COORDINATOR:
			starter = coord.GetStarter()
			netMsg := ToNetCoordinatorMsg(coord)
			//smlog.Println("\033[42m\033[1;30mSENDING COORDINATOR [", coord, "] to ", prossimoAddr, "\033[0m")
			smlog.Warn(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING COORD %v to %s"+ColorReset, netMsg, prossimoAddr)
			//	log.Printf("\033[42m\033[1;30mSENDING COORD %[1]v to %[2]s \033[0m", coord, prossimoAddr)
			_, errq = csN.ForwardCoordinator(ctx, netMsg)
			break
		case MSG_HEARTBEAT:
			netMsg := ToNetHeartbeatMsg(hb)
			smlog.Warn(LOG_MSG_RECV, ColorBlkBckgrGreen+BoldBlack+"SENDING HB %v to %s"+ColorReset, netMsg, prossimoAddr)
			_, errq = csN.SendHeartBeat(ctx, netMsg)
			break
		default:
			break
		}
		//_, errq := csN.InoltraElezione(ctx, elezione)
		if errq != nil {
			if attempts != RMI_RETRY_TOLERANCE {
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
				if tipo == MSG_HEARTBEAT {
					SuccessfulHB--
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
					if prossimoAddr == Me.GetFullAddr() {
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
