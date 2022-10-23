// hbmgt
package statemachine

import (
	. "bully/pkg/node/env"
	. "bully/pkg/node/net"
	. "bully/tools/smlog"
	smlog "bully/tools/smlog"

	"time"
)

var EventsSend chan (string)
var EventsList chan (string)

const HB_TIMEOUT = 1000 * time.Millisecond
const HB_TOLERANCE = HB_TIMEOUT * 2 //250 * time.Millisecond
type hbMgtState int8

const (
	HB_HALT hbMgtState = iota
	HB_SEND
	HB_LISTEN
)

type hbMgtProperty int8

const (
	HB_PROPERTY_LISTENING hbMgtProperty = iota
	HB_PROPERTY_SENDING
)

var coordTimer *time.Ticker
var noncoordTimer *time.Ticker
var SendingHB = false
var ListeningtoHb = false

func setHbMgt(state hbMgtState) {
	smlog.InfoU("setting HBmgt to %v", state)
	switch state {
	case HB_HALT:
		setProperty(HB_PROPERTY_LISTENING, false)
		setProperty(HB_PROPERTY_SENDING, false)
		break
	case HB_SEND:
		setProperty(HB_PROPERTY_LISTENING, false)
		setProperty(HB_PROPERTY_SENDING, true)
		break
	case HB_LISTEN:
		setProperty(HB_PROPERTY_LISTENING, true)
		setProperty(HB_PROPERTY_SENDING, false)
		break
	}
}

func setProperty(prop hbMgtProperty, val bool) {
	switch prop {
	case HB_PROPERTY_LISTENING:
		if val {
			if !ListeningtoHb {
				ListeningtoHb = true
				go ListenToHb()
			}
		} else {
			if ListeningtoHb {
				ListeningtoHb = false
				EventsList <- "Stop1"
				noncoordTimer.Stop()
			}
		}
		break
	case HB_PROPERTY_SENDING:
		if val {
			if !SendingHB {
				SendingHB = true
				go InviaHB()
			}
		} else {
			if SendingHB {
				SendingHB = false
				EventsSend <- "STOP2"
				coordTimer.Stop()
			}
		}
		break
	}
}

/*
func startListening() {
	if !ListeningtoHb {
		ListeningtoHb = true
		go ListenToHb()
	}
}
func stopListening() {
	if ListeningtoHb {
		ListeningtoHb = false
		EventsList <- "Stop1"
		noncoordTimer.Stop()
	}
}
func startSending() {
	if !SendingHB {
		SendingHB = true
		go InviaHB()
	}
}
func stopSending() {
	if SendingHB {
		SendingHB = false
		EventsSend <- "STOP2"
		coordTimer.Stop()
	}
}

func haltHb() {
	if ListeningtoHb {
		ListeningtoHb = false
		EventsList <- "Stop1"
		noncoordTimer.Stop()
	}
	if SendingHB {
		SendingHB = false
		EventsSend <- "STOP2"
		coordTimer.Stop()
	}

}
*/
// da separare nel comportamento
func ListenToHb() {
	smlog.InfoU("inizio routine di ascolto hb...")
	interrupt := false
	noncoordTimer = time.NewTicker(HB_TIMEOUT + HB_TOLERANCE)
	for {
		select {
		case in := <-Heartbeat:
			smlog.Info(LOG_HB, "confermo hb")
			if in.GetId() != CoordId && CoordId != -1 {
				// more than one coordinators in the network!
				// it can happen when a large (>=15) number
				// of nodes join at the same time
				// in this case we need a new election
				smlog.Error(LOG_HB, "Received HB from %d, that is not my coordinator %d", in.GetId(), CoordId)
				smlog.Error(LOG_HB, "Starting new election...")
				go startElection()
				IsElectionStarted = true
				interrupt = true
			}
			smlog.Info(LOG_HB, "confermo hb")
			noncoordTimer.Reset(HB_TIMEOUT + HB_TOLERANCE)
			break
		case <-noncoordTimer.C:
			smlog.Critical(LOG_HB, "non sento più il coord")
			//Events <- "STOP1"
			go startElection()
			IsElectionStarted = true
			interrupt = true
			break
		case <-EventsList:
			//if in == "Stop1" {
			smlog.Info(LOG_HB, "devo smettere di ascoltare perché c'è una elezione in corso")
			interrupt = true
			//}
			break

		}
		if interrupt {
			break
		}
	}
	noncoordTimer.Stop()
	smlog.Info(LOG_HB, "Esco dalla routine di ascolto HB...")
	ListeningtoHb = false
	//setHbMgt(HB_HALT)
}

func InviaHB() {
	interrupt := false
	coordTimer = time.NewTicker(HB_TIMEOUT)
	//defer coordTimer.Stop()
	//failedNodeExistence := true
	/*
		if len(allNodesList.GetList()) == 1 {
			// se ci sono solo io, evito direttamente
			smlog.Info(LOG_UNDEFINED, "Sono rimasto solo io")
			//events <- "STOP"
			interrupt = true // vedere
		}*/
	allNodesList := AskForAllNodes()
	hbMsg := &MsgHeartbeat{Id: Me.GetId()}
	for {
		select {
		case <-coordTimer.C:
			smlog.Critical(LOG_UNDEFINED, "SuccessfulHB=%v", SuccessfulHB)
			if SuccessfulHB == 0 {
				interrupt = true
				SuccessfulHB = -1
				break
			}
			SuccessfulHB = len(allNodesList) - 1
			//smlog.Critical(LOG_UNDEFINED, "***** INIZIALIZZO SuccessfulHB=%v", SuccessfulHB)
			for _, node := range allNodesList {
				//				node := ToSMNode(nodenet)
				if node.GetFullAddr() != Me.GetFullAddr() {
					smlog.Info(LOG_HB, "Invio HB al nodo %d, presso %s", node.GetId(), node.GetFullAddr())
					// mando gli hb in parallelo! tanto devo mandarli a tutti
					go SafeRMI(MSG_HEARTBEAT, node, false, nil, nil, nil, hbMsg)
					/*if rmiErr {
						failedNodeExistence = true
					}*/
				}
			} // qui ho inviato gli hb a tutti i nodi

			smlog.Info(LOG_HB, "Inviati tutti gli HB, attendo timer...")
			break
		case in := <-EventsSend:
			if in == "STOP2" { // fare canali differenti?
				smlog.InfoU("arrivato evento di STOP: %s", in)
				interrupt = true
			}
			break
		}
		if interrupt {
			break
		}
	}
	coordTimer.Stop()
	smlog.Info(LOG_HB, "Esco dalla routine di invio HB...")
	SendingHB = false
	//setHbMgt(HB_HALT)
}
func initHbMgt() {

	EventsSend = make(chan string)
	EventsList = make(chan string)
}
