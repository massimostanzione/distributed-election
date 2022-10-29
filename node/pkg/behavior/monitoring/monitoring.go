package monitoring

import (
	. "distributedelection/node/pkg/env"
	. "distributedelection/node/pkg/net"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"

	"time"
)

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

var timer *time.Ticker
var coordTimer *time.Ticker

//var COORDTIMER_DURATION =
//var NONCOORDTIMER_DURATION =

var noncoordTimer *time.Ticker
var SendingHB = false
var ListeningtoHb = false
var EventsSend chan (string)
var EventsList chan (string)
var MonitoringChannel chan (string)
var stateChan chan (string)
var statee string

func InitMonitoring() {
	Heartbeat = make(chan *MsgHeartbeat)
	EventsSend = make(chan string)
	EventsList = make(chan string)
	stateChan = make(chan string)
	statee = ""
	coordTimer = time.NewTicker(time.Duration(Cfg.HB_TIMEOUT) * time.Millisecond)
	noncoordTimer = time.NewTicker(time.Duration(Cfg.HB_TIMEOUT+Cfg.HB_TOLERANCE) * time.Millisecond)
	// https://github.com/golang/go/issues/12721
	coordTimer.Stop()
	noncoordTimer.Stop()
	timer = time.NewTicker(1 * time.Millisecond)
	timer.Stop()
	go monitoring()
	go InviaHB()
	go ListenToHb()
	SetMonitoringState(HB_HALT)
}
func SetMonitoringState(state hbMgtState) {
	switch state {
	case HB_HALT:
		stateChan <- "halt"
		//setProperty(HB_PROPERTY_LISTENING, false)
		//setProperty(HB_PROPERTY_SENDING, false)
		break
	case HB_SEND:
		stateChan <- "send"
		//setProperty(HB_PROPERTY_LISTENING, false)
		//setProperty(HB_PROPERTY_SENDING, true)
		break
	case HB_LISTEN:
		stateChan <- "listen"
		//setProperty(HB_PROPERTY_LISTENING, true)
		//setProperty(HB_PROPERTY_SENDING, false)
		break
	}
}

func setProperty(prop hbMgtProperty, val bool) {
	switch prop {
	case HB_PROPERTY_LISTENING:
		if val {
			/*if !ListeningtoHb {
				ListeningtoHb = true
				go ListenToHb()
			}*/
			EventsList <- "run"
		} else {
			/*if ListeningtoHb {
				noncoordTimer.Stop()
				ListeningtoHb = false
				smlog.Debug(LOG_STATEMACHINE, "inizio")
				MonitoringChannel <- "elect"
				EventsList <- "Stop1"
				smlog.Debug(LOG_STATEMACHINE, "fine")
			}*/
			noncoordTimer.Stop()
			EventsList <- "stop"
		}
		break
	case HB_PROPERTY_SENDING:
		if val {
			EventsSend <- "run"
			/*
				if !SendingHB {
					SendingHB = true
					go InviaHB()
				}*/
		} else {
			coordTimer.Stop()
			EventsSend <- "stop"
			/*
				if SendingHB {
					SendingHB = false
					EventsSend <- "STOP2"
					coordTimer.Stop()
				}*/
		}
		break
	}
}

func monitoring() {
	for {
		select {
		case staterecv := <-stateChan:
			statee = staterecv
			smlog.Critical(LOG_HB, "staterecv: RICEVO %s", statee)
			if statee == "send" {
				timer.Reset(time.Duration(Cfg.HB_TIMEOUT) * time.Millisecond)
			}
			if statee == "listen" {
				timer.Reset(time.Duration(Cfg.HB_TIMEOUT+Cfg.HB_TOLERANCE) * time.Millisecond)
			}
			if statee == "halt" {
				timer.Stop()
			}
			break
		case <-Heartbeat:
			smlog.Critical(LOG_HB, "arriva hb")
			if statee == "listen" {
				acknowledgeHb()
			}
			break
		case <-timer.C:
			smlog.Critical(LOG_HB, "arriva timer")
			if statee == "send" {
				sendhb()
			}
			if statee == "listen" {
				timer.Stop()
				statee = "halt"
				smlog.Critical(LOG_HB, "non sento più il coord")
				MonitoringChannel <- "elect"
			}
			if statee == "halt" {
				timer.Stop()
			}
			break
		}
	}
}
func acknowledgeHb() {
	/*if in.GetId() != CoordId && CoordId != -1 {
		// more than one coordinators in the network!
		// it can happen when a large (>=15) number
		// of nodes join at the same time
		// in this case we need a new election
		smlog.Error(LOG_HB, "Received HB from %d, that is not my coordinator %d", in.GetId(), CoordId)
		smlog.Error(LOG_HB, "Starting new election...")
		MonitoringChannel <- "elect"
		interrupt = true
	}*/
	smlog.Info(LOG_HB, "confermo hb")
	timer.Reset(time.Duration(Cfg.HB_TIMEOUT+Cfg.HB_TOLERANCE) * time.Millisecond)
}

func sendhb() {
	timer.Reset(time.Duration(Cfg.HB_TIMEOUT) * time.Millisecond)
	//defer coordTimer.Stop()
	//failedNodeExistence := true
	/*
		if len(allNodesList.GetList()) == 1 {
			// se ci sono solo io, evito direttamente
			smlog.Info(LOG_UNDEFINED, "Sono rimasto solo io")
			//events <- "STOP"
			interrupt = true // vedere
		}*/
	allNodesList := AskForAllNodesList()
	hbMsg := &MsgHeartbeat{Id: Me.GetId()}
	smlog.Critical(LOG_UNDEFINED, "SuccessfulHB=%v", SuccessfulHB)
	if SuccessfulHB == 0 {
		//		interrupt = true
		SuccessfulHB = -1
		//		break
	}
	SuccessfulHB = len(allNodesList) - 1
	//smlog.Critical(LOG_UNDEFINED, "***** INIZIALIZZO SuccessfulHB=%v", SuccessfulHB)
	for _, node := range allNodesList {
		//				node := ToSMNode(nodenet)
		if node.GetFullAddr() != Me.GetFullAddr() {
			smlog.Info(LOG_HB, "Invio HB al nodo %d, presso %s", node.GetId(), node.GetFullAddr())
			// mando gli hb in parallelo! tanto devo mandarli a tutti
			//go SafeRMI(MSG_HEARTBEAT, node, false, nil, nil, hbMsg)
			go SafeHB(ToNetHeartbeatMsg(hbMsg), node)

		}
	}
}

func ListenToHb() {
	smlog.InfoU("inizio routine di ascolto hb...")
	for {
		state := <-EventsList
		smlog.Critical(LOG_HB, "LISTEN: RICEVO %s", state)
		if state == "run" {
			interrupt := false
			noncoordTimer.Reset(time.Duration(Cfg.HB_TIMEOUT+Cfg.HB_TOLERANCE) * time.Millisecond)
			for {
				select {
				case <-Heartbeat:
					/*if in.GetId() != CoordId && CoordId != -1 {
						// more than one coordinators in the network!
						// it can happen when a large (>=15) number
						// of nodes join at the same time
						// in this case we need a new election
						smlog.Error(LOG_HB, "Received HB from %d, that is not my coordinator %d", in.GetId(), CoordId)
						smlog.Error(LOG_HB, "Starting new election...")
						MonitoringChannel <- "elect"
						interrupt = true
					}*/
					smlog.Info(LOG_HB, "confermo hb")
					noncoordTimer.Reset(time.Duration(Cfg.HB_TIMEOUT+Cfg.HB_TOLERANCE) * time.Millisecond)
					break
				case <-noncoordTimer.C:
					smlog.Critical(LOG_HB, "non sento più il coord")
					MonitoringChannel <- "elect"
					//Events <- "STOP1"
					noncoordTimer.Stop()
					//interrupt = true
					break
				case in := <-EventsList:
					smlog.Critical(LOG_HB, ">>> leggo da EVENTSLIST: %s", in)
					if in == "stop" {
						smlog.Info(LOG_HB, "devo smettere di ascoltare perché c'è una elezione in corso")
						interrupt = true

					}
					break

				}
				if interrupt {
					break
				}
			}
			noncoordTimer.Stop()
			smlog.Info(LOG_HB, "Esco dalla routine di ascolto HB...")
		}
	}
	//ListeningtoHb = false
	//-------------------
	/*
		for {
			interrupt := false
			noncoordTimer = time.NewTicker(time.Duration(Cfg.HB_TIMEOUT+Cfg.HB_TOLERANCE) * time.Millisecond)
			for {
				select {
				case <-Heartbeat:
					/*if in.GetId() != CoordId && CoordId != -1 {
						// more than one coordinators in the network!
						// it can happen when a large (>=15) number
						// of nodes join at the same time
						// in this case we need a new election
						smlog.Error(LOG_HB, "Received HB from %d, that is not my coordinator %d", in.GetId(), CoordId)
						smlog.Error(LOG_HB, "Starting new election...")
						MonitoringChannel <- "elect"
						interrupt = true
					}*/ /*
					smlog.Info(LOG_HB, "confermo hb")
					noncoordTimer.Reset(time.Duration(Cfg.HB_TIMEOUT+Cfg.HB_TOLERANCE) * time.Millisecond)
					break
				case <-noncoordTimer.C:
					smlog.Critical(LOG_HB, "non sento più il coord")
					MonitoringChannel <- "elect"
					//Events <- "STOP1"
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
		}
		noncoordTimer.Stop()
		smlog.Info(LOG_HB, "Esco dalla routine di ascolto HB...")
		ListeningtoHb = false
		//SetMonitoringState(HB_HALT)
	*/
}
func InviaHB() {
	smlog.InfoU("inizio routine di invio hb...")
	for {
		state := <-EventsSend
		smlog.Critical(LOG_HB, "SEND: RICEVO %s", state)
		if state == "run" {

			interrupt := false
			coordTimer.Reset(time.Duration(Cfg.HB_TIMEOUT) * time.Millisecond)
			//defer coordTimer.Stop()
			//failedNodeExistence := true
			/*
				if len(allNodesList.GetList()) == 1 {
					// se ci sono solo io, evito direttamente
					smlog.Info(LOG_UNDEFINED, "Sono rimasto solo io")
					//events <- "STOP"
					interrupt = true // vedere
				}*/
			allNodesList := AskForAllNodesList()
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
							//go SafeRMI(MSG_HEARTBEAT, node, false, nil, nil, hbMsg)
							go SafeHB(ToNetHeartbeatMsg(hbMsg), node)
						}
					} // qui ho inviato gli hb a tutti i nodi

					smlog.Info(LOG_HB, "Inviati tutti gli HB, attendo timer...")
					break
				case in := <-EventsSend:
					if in == "stop" { // fare canali differenti?
						smlog.InfoU("arrivato evento di STOP: %s", in)
						interrupt = true
					}
					break
				}
				if interrupt {
					break
				}
			}
			//SendingHB = false
			//SetMonitoringState(HB_HALT)

			coordTimer.Stop()
			smlog.Info(LOG_HB, "Esco dalla routine di invio HB...")

		}
	}
	//---------------------------
	/*
		interrupt := false
		coordTimer = time.NewTicker(time.Duration(Cfg.HB_TIMEOUT) * time.Millisecond)
		//defer coordTimer.Stop()
		//failedNodeExistence := true
		/*
			if len(allNodesList.GetList()) == 1 {
				// se ci sono solo io, evito direttamente
				smlog.Info(LOG_UNDEFINED, "Sono rimasto solo io")
				//events <- "STOP"
				interrupt = true // vedere
			}*/ /*
		allNodesList := AskForAllNodesList()
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
						//go SafeRMI(MSG_HEARTBEAT, node, false, nil, nil, hbMsg)
						go SafeHB(ToNetHeartbeatMsg(hbMsg), node)
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
		//SetMonitoringState(HB_HALT)*/
}
