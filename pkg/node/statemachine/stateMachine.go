package statemachine

import (
	. "fredricksonLynch/pkg/node/env"
	. "fredricksonLynch/pkg/node/net"
	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"

	//	smlog "fredricksonLynch/smlog"
	"fmt"

	//	"log"
	//	"os"

	//"strconv"
	"time"
	//	loggo "github.com/juju/loggo"
)

type nodeState uint8 //TODO spostare altrove?

var currentState = STATE_UNDEFINED

var nonCoordTimer *time.Timer

const (
	STATE_UNDEFINED nodeState = iota
	STATE_JOINING
	STATE_ELECTION_STARTER
	STATE_ELECTION_VOTER
	STATE_COORDINATOR
	STATE_NON_COORDINATOR
)

//var smlog = log.New(os.Stderr, "[SM] "+time.Now().Format("15:04:05")+" "+currentState.Short()+" "+starterToLogout(starter)+" ", 0)

func (state nodeState) Short() string {
	switch state {
	case STATE_UNDEFINED:
		return "N/A  "
	case STATE_JOINING:
		return "JOINI"
	case STATE_ELECTION_STARTER:
		return "START"
	case STATE_ELECTION_VOTER:
		return "VOTER"
	case STATE_COORDINATOR:
		return "COORD"
	case STATE_NON_COORDINATOR:
		return "NONCO"
	}
	return "err"
}

func setState(state nodeState) {
	currentState = state
	//	smlog.SetState(stateToLogout(currentState))
	//logger := loggo.GetLogger("")
	//loggo.ConfigureLoggers("eeeeee")
	//logger.SetLogLevel(loggo.TRACE)
	//logger.Infof("new state: ", stateToLogout(currentState))
	smlog.SetState(state.Short())
	smlog.Info(LOG_STATEMACHINE, "new state: %s", currentState.Short())
}

func StartStateMachine() {
	Heartbeat = make(chan *MsgHeartbeat, 1)
	Events = make(chan string, 1)
	ElectionChannel = make(chan *MsgElection, 1)
	CoordChannel = make(chan *MsgCoordinator, 1)
	fmt.Println("[SM] Time     Lvl   State Event  Description")
	fmt.Println("[SM] -------- ----- ----- ------ ---------------")
	smlog.Info(LOG_UNDEFINED, "Starting SM...")
	smlog.InfoU("Type CTRL+C to terminate")
	smlog.InfoU("Type CTRL+Z to Pause/resume")
	InitializeNetMW()

	setState(STATE_JOINING)
	go run()
	Listen()
}

func run() {
	for {
		for Pause {
		}
		smlog.Info(LOG_STATEMACHINE, "Running state cycle")
		//smlog.Println("*** RUNNING STATE: ", (msgType)(currentState))
		switch currentState {
		case STATE_JOINING: // 1
			state_joining()
			break
		case STATE_ELECTION_STARTER: // 2
			state_election_starter()
			break
		case STATE_ELECTION_VOTER: // 3
			state_election_voter()
			break
		case STATE_COORDINATOR: // 4
			state_coordinator()
			break
		case STATE_NON_COORDINATOR: // 5
			state_nonCoordinator()
			break
		default: //TODO
		}
	}
}

func state_joining() {
	for Pause {
	}
	//	starter = false
	Me = AskForJoining()
	startElection()
	//TODO di seguito togliere?
	//if nextId != Me.GetId() {
	//	} else {
	//	setState(STATE_COORDINATOR)
	//	}
	//startElezione(ctx)
}

func state_election_starter() {
	late_hb_received := 0
	for {
		for Pause {
		}
		select {
		case inp := <-ElectionChannel:
			if inp.GetStarter() == Me.GetId() {
				smlog.Info(LOG_ELECTION, "- nomino il coord")
				electedId := elect(inp.GetVoters())
				// --------------
				nextNode := AskForNodeInfo(Me.GetId()+1, true)
				sendCoord(NewCoordinatorMsg(Me.GetId(), electedId), nextNode)
			} else {
				// c'è un altro starter, quindi più elezioni in giro
				// ma va bene così, voto e faccio girare
				vote(inp)
			}
			break
		case inp := <-CoordChannel:
			if inp.GetStarter() == Me.GetId() {
				endElection(inp, inp.GetStarter() != Me.GetId())
			} else {
				// sono in attesa della mia elezione, non posso terminarla
				// quindi semplicemente faccio girare
				// non accetto quindi coordinatori che non vengano dalla mia elezione
				nextNode := AskForNodeInfo(Me.GetId()+1, true)
				sendCoord(NewCoordinatorMsg(inp.GetStarter(), inp.GetCoordinator()), nextNode)
			}
			//endElection(inp, inp.GetStarter() != Me.GetId())
			break
		case <-Heartbeat:
			// non dovrei riceverli qui, quindi ne ignoro alcuni
			// (potrebbero essere residui di elezioni precedenti)
			// se continuano ad arrivare vuol dire che c'è qualcosa che non va
			// quindi nel caso faccio partire nuova elezione
			late_hb_received++
			if late_hb_received == LATE_HB_TOLERANCE {
				startElection()
			}
			break
		}
		break
	}

}

func state_election_voter() {
	late_hb_received := 0
	for {
		for Pause {
		}
		select {
		case inp := <-ElectionChannel:
			/*	if inp.GetStarter() == Me.GetId() {
				smlog.Info(LOG_ELECTION, "- nomino il coord, ma da ELECTION_VOTER")
				electedId := elect(inp.GetVoters())
				nextNode := AskForNodeInfo(Me.GetId()+1, true)
				sendCoord(NewCoordinatorMsg(Me.GetId(), electedId), nextNode)
			} else {*/
			// ho già votato, l'unica cosa che devo aspettare è un COORD,
			// quindi se ci sono altre elezioni semplicemente voto
			vote(inp)
			//}
			break
		case inp := <-CoordChannel:
			endElection(inp, inp.GetStarter() != Me.GetId())
			break
		case <-Heartbeat:
			// non dovrei riceverli qui, quindi ne ignoro alcuni
			// (potrebbero essere residui di elezioni precedenti)
			// se continuano ad arrivare vuol dire che c'è qualcosa che non va
			// quindi nel caso faccio partire nuova elezione
			late_hb_received++
			if late_hb_received == LATE_HB_TOLERANCE {
				startElection()
			}
			break
		}
		break
	}
}

func state_coordinator() {
	late_hb_received := 0
	//setState(STATE_ELECTION_VOTER)
	confirmedCoord := false
	for {
		for Pause {
		}
		if !confirmedCoord {
			go HBroutine()
		}
		select {
		case inp := <-ElectionChannel:
			Events <- "STOP"
			if inp.GetStarter() == Me.GetId() {
				smlog.Fatal(LOG_ELECTION, "sono coord ma mi è arrivata elezione da Me!")
				/*smlog.Info(LOG_ELECTION, "- nomino il coord, ma da COORDINATOR")
				electedId := elect(inp.GetVoters())
				// --------------
				nextNode := AskForNodeInfo(Me.GetId()+1, true)
				sendCoord(NewCoordinatorMsg(Me.GetId(), electedId), nextNode)*/
				//sendCoord(electedId, nextAddr)
			} else {
				setState(STATE_ELECTION_VOTER)
				vote(inp)
			}
			break
		case inp := <-CoordChannel:
			if inp.GetCoordinator() == Me.GetId() {
				confirmedCoord = true
			} else {
				smlog.Fatal(LOG_ELECTION, "sono coord, ma mi è arrivato un altro COORD senza elezioni, ad elezioni chiuse")
				Events <- "STOP"
				//TODO se coord resto io, posso risparmiarmelo?
				//if inp.GetCoordinator()!=Me.GetId(){
				//NOTA anche se sono lo stesso, chiamo endElection visto che il comportamento è lo stesso

				//smlog.Println("pppppppqqqqqqq")
				endElection(inp, inp.GetStarter() != Me.GetId())
			}
			break
		case <-Heartbeat:
			// non dovrei riceverli qui, quindi ne ignoro alcuni
			// (potrebbero essere residui di elezioni precedenti)
			// se continuano ad arrivare vuol dire che c'è qualcosa che non va
			// quindi nel caso faccio partire nuova elezione
			late_hb_received++
			if late_hb_received == LATE_HB_TOLERANCE {
				startElection()
			}
			break
		}
		//if !confirmedCoord {
		break
		//}
	}
}

func state_nonCoordinator() {
	late_hb_received := 0
	nonCoordTimer = time.NewTimer(HB_TIMEOUT + HB_TOLERANCE)
	//go listenHB()
	for {
		if Pause {
			nonCoordTimer.Stop()
			for Pause {
			}
			nonCoordTimer.Reset(HB_TIMEOUT + HB_TOLERANCE)
		}
		select {
		/*case <-Heartbeat:
		//TODO if val == "HB" reset, invece di farlo altrove?
		//smlog.Printf("Ricevuto HB dal canale")
		//nonCoordTimer.Stop()
		//nonCoordTimer = time.NewTimer(HB_TIMEOUT + HB_TOLERANCE)
		nonCoordTimer.Reset(HB_TIMEOUT + HB_TOLERANCE)
		break // TODO serve?*/
		case inp := <-Heartbeat:
			if inp.GetId() != CoordId {
				late_hb_received++
				if late_hb_received == LATE_HB_TOLERANCE {
					nonCoordTimer.Stop()
					startElection()
				}
			} else {
				nonCoordTimer.Reset(HB_TIMEOUT + HB_TOLERANCE)
			}
			break
		case <-nonCoordTimer.C:
			nonCoordTimer.Stop()
			smlog.Critical(LOG_UNDEFINED, "\033[41m*** COORDINATOR FAILURE DETECTED! ***\033[0m")
			//nonCoordTimer.Stop()
			//			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			//locCtx = ctx
			//		defer cancel()
			//nonCoordTimer.Stop()
			failedNode := AskForNodeInfo(CoordId, false) //cs.GetNode(ctx, &pb.NodeId{Id: CoordId})
			DeclareNodeState(failedNode, false)

			// qui non posso segnalare nulla sugli altri: so solo che COORD è failed,
			// e che io sono vivo
			// aggiorno quindi il centrale su di Me

			DeclareNodeState(Me, true)

			//			startElezione(ctx)
			//stopList = true
			startElection()
			break

		case inp := <-ElectionChannel:
			nonCoordTimer.Stop()
			//Events <- "STOP"
			if inp.GetStarter() == Me.GetId() {
				smlog.Fatal(LOG_ELECTION, "sono coord ma mi è arrivata elezione da Me!")
				/*			//TODO questa proceura è ripetuta, unire

							//TODO fault tolerance - non bloccare tutto con Fatal
										smlog.Info(LOG_ELECTION, "- nomino il coord, ma da NON_COORDINATOR")
										electedId := elect(inp.GetVoters())
										// --------------
										nextNode := AskForNodeInfo(Me.GetId()+1, true)
										sendCoord(NewCoordinatorMsg(Me.GetId(), electedId), nextNode)*/
				//sendCoord(electedId, nextAddr)
			} else {
				setState(STATE_ELECTION_VOTER)
				vote(inp)
			}
			break
		case inp := <-CoordChannel:
			nonCoordTimer.Stop()
			//Events <- "STOP" // TODO serve?
			//TODO se coord resto io, posso risparmiarmelo?
			//if inp.GetCoordinator()!=Me.GetId(){
			//NOTA anche se sono lo stesso, chiamo endElection visto che il comportamento è lo stesso

			//smlog.Println("pppppppqqqqqqq")
			//NOTA anche se sono lo stesso, chiamo endElection visto che il comportamento è lo stesso
			endElection(inp, inp.GetStarter() != Me.GetId())
			//}
			break
		}
		break
	}
}
