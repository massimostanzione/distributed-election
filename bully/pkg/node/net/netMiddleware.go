// netMiddleware
package net

//var cs pb.DistGrepClient
import (
	pb "bully/pb/node"
	"context"

	. "bully/pkg/node/env"
	//"bully/pkg/node/statemachine"
	. "bully/tools/formatting"
	//"bully/pkg/node"

	. "bully/tools/smlog"
	smlog "bully/tools/smlog"
	"net"
	"time"

	"google.golang.org/grpc"
)

type DGnode struct {
	pb.UnimplementedDistGrepServer
}

var NONE = &pb.NONE{}
var cs pb.DistGrepClient

var w *grpc.Server
var lis net.Listener
var serverConn *grpc.ClientConn //server

const RMI_RETRY_TOLERANCE = 3
const LATE_HB_TOLERANCE = 3

func InitializeNetMW() {

	// il centrale espone il servizio di identificazione dei nodi

	conn := ConnectToNode(ServRegAddr)
	serverConn = conn
	/*	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		serverConn = conn
		if err != nil {
			smlog.Fatal("Error while contacting server on %v:\n %v", serverAddr, err)
		}
	*/
	//defer conn.Close()
	//	locCtx = ctx
	// MI METTO IN ASCOLTO: la porta su cui ascolto è
	// la stessa che invierò al gestore dell'anello,
	// in quanto è quella sulla quale sarò contattato
	// Start listening for incoming calls
	//port := "40046"
	liss, err := net.Listen("tcp", Me.GetFullAddr())
	lis = liss
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to listen to port %v:\n%v", Me.GetPort(), err)
	}
	// New server instance and service registering
	w = grpc.NewServer()
	pb.RegisterDistGrepServer(w, &DGnode{})
	// Serve incoming calls

	// Defining client interface, to be used to invoke the bully service
	cs = pb.NewDistGrepClient(serverConn)
}

func ConnectToNode(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while contacting server on %v:\n %v", addr, err)
	}
	return conn
}

func Listen() {
	smlog.Info(LOG_NETWORK, "Listening on port %v.", Me.GetPort())
	if err := w.Serve(lis); err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while trying to serve request: %v", err)
	}
	//	for pause {
	//	}
}

func StopServer() {
	smlog.Info(LOG_NETWORK, "Stopping server")
	w.Stop()
	//	for !pause {
	//	}
}
func SwitchServerState(run bool) {
	if !run {
		Listen()
	} else {
		StopServer()
	}
}
func contactServiceReg() *grpc.ClientConn {

	conn := ConnectToNode(ServRegAddr)
	defer conn.Close() //chiusura, se porta problemi controllare
	return conn
}
func AskForJoining() *SMNode {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	smlog.Info(LOG_SERVREG, "asking for joining the ring...")
	node, err := cs.JoinRing(ctx, &pb.NodeAddr{Host: Me.GetHost(), Port: Me.GetPort()})
	if err != nil {
		smlog.Fatal(LOG_NETWORK, "Error while executing bully:\n%v", err)
	}
	return ToSMNode(node)
}

//func askForNodeInfo(i int32, forceRunningNode bool) (int32, string) {
//TODO B: strutture come forceRunningNode non servono più, non c'è struttura ad anello
func AskForNodeInfo(i int32, forceRunningNode bool) *SMNode {
	smlog.Info(LOG_SERVREG, "Chiedo al centrale informazioni sul nodo %d", i)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()
	if forceRunningNode {
		ret, errr := cs.GetNextRunningNode(ctx, &pb.NodeId{Id: int32(i)})
		if errr != nil {
			smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
			return nil
		}
		return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	} else {
		ret, errr := cs.GetNode(ctx, &pb.NodeId{Id: int32(i)})
		if errr != nil {
			smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
			return nil
		}
		return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	}

}

func AskForNodesWithGreaterIds(baseId int32, forceRunningNode bool) []*SMNode {
	smlog.Trace(LOG_SERVREG, "Chiedo al centrale informazioni sui nodi con id > %d", baseId)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()
	if forceRunningNode {
		//de facto not implemented
		_, errr := cs.GetNextRunningNode(ctx, &pb.NodeId{Id: int32(baseId)})
		if errr != nil {
			smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
			return nil
		}
		return nil //&SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	} else {
		ret, errr := cs.GetAllNodesWithIdGreaterThan(ctx, &pb.NodeId{Id: int32(baseId)})
		if errr != nil {
			smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
			return nil
		}
		//conversion
		//TODO implementare anche nelle altre chiamate simili
		var array []*SMNode
		for _, node := range ret.GetList() {
			array = append(array, ToSMNode(node))
		}
		return array //&SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	}

}
func AskForAllNodes() []*SMNode {
	smlog.Trace(LOG_SERVREG, "Chiedo al centrale informazioni su TUTTI i nodi,")
	smlog.Trace(LOG_SERVREG, "attivi o meno, per verificare stato di essi")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()

	ret, errr := cs.GetAllNodes(ctx, NONE)
	if errr != nil {
		smlog.Fatal(LOG_NETWORK, "errore in GETNODO:\n%v", errr)
		return nil
	}
	//conversion
	//TODO implementare anche nelle altre chiamate simili
	var array []*SMNode
	for _, node := range ret.GetList() {
		array = append(array, ToSMNode(node))
	}
	return array //&SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
}

/*
func HBroutine() {
	interrupt := false
	coordTimer := time.NewTicker(HB_TIMEOUT)
	//defer coordTimer.Stop()
	failedNodeExistence := true
	var allNodesList *pb.NodeList

	for {
		if failedNodeExistence {
			smlog.Debug(LOG_SERVREG, "vado a chiedere i nodi in piedi")
			var errl error = error(nil)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			allNodesList, errl = cs.GetAllRunningNodes(ctx, NONE)
			if errl != nil {
				smlog.Fatal(LOG_SERVREG, "problema nel reperire AllRunningNodes: %v", errl)
			}
			failedNodeExistence = false

			if len(allNodesList.GetList()) == 1 {
				// se ci sono solo io, evito direttamente
				smlog.Info(LOG_UNDEFINED, "Sono rimasto solo io")
				//events <- "STOP"
				interrupt = true
			}
		}
		for _, nodenet := range allNodesList.GetList() {
			node := ToSMNode(nodenet)
			if node.GetFullAddr() != Me.GetFullAddr() {
				// TODO fare funzione che, dato un nodo, ritorna la connessione con esso

				nodoServer := grpc.NewServer()
				pb.RegisterDistGrepServer(nodoServer, &DGnode{})
				smlog.Info(LOG_HB, "Invio HB al nodo %d, presso %s", node.GetId(), node.GetFullAddr())

				// mi segnalo come vivo
				ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
				defer cancel2()
				_, erro := cs.ReportAsRunning(ctx2, &pb.Node{Id: Me.GetId()})
				if erro != nil {
					smlog.Error(LOG_UNDEFINED, "problema nel segnalarmi vivo mentre mando HB: %v", erro)
				}
				// TODO nota per relazione e documentazione:
				// uso il parametro FALSE perché, a differenza degli altri casi,
				// sto *enumerando* i nodi a cui inviare l'HB,
				// mentre negli altri casi vado alla meglio cercando il successivo in piedi
				// in questo caso il parametro di ritorno mi indicherà non solo la presenza generica
				// di un nodo fallito, ma il fatto che il nodo fallito è proprio quello che ho provato
				rmiErr := SafeRMI(MSG_HEARTBEAT, node, false, nil, nil, &pb.Heartbeat{Id: Me.GetId()})
				if rmiErr {
					failedNodeExistence = true
				}
			}
		} // qui ho inviato gli hb a tutti i nodi

		smlog.Info(LOG_HB, "Inviati tutti gli HB, attendo timer...")
		select {
		case <-coordTimer.C:
			//smlog.InfoU("SCATTA IL TIMER")
			// semplicemente vai avanti
			break
		case val := <-Events:
			if val == "STOP" {
				smlog.InfoU("------------------ ARRIVATO EVENTO DI \"STOP\"")
			}
			coordTimer.Stop()
			interrupt = true
			break
		}

		if interrupt {
			break
		}

	}
	smlog.Info(LOG_HB, "Esco dalla routine di invio HB...")
}
*/
//TODO gestire fallimenti temporanei rispetto al servReg
func DeclareNodeState(node *SMNode, running bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if running {
		//smlog.Info(LOG_NETWORK, "\033[41m*** dichiaro running il nodo %d ***\033[0m", node.GetId())
		_, err := cs.ReportAsRunning(ctx, ToNetNode(*node))
		if err != nil {
			smlog.Fatal(LOG_UNDEFINED, err.Error())
		}
	} else {
		smlog.Error(LOG_NETWORK, ColorRedBckgrWhite+"*** dichiaro failed il nodo %d ***"+ColorReset, node.GetId())
		_, err := cs.ReportAsFailed(ctx, ToNetNode(*node))
		if err != nil {
			smlog.Fatal(LOG_UNDEFINED, err.Error())
		}
	}
}

/*
func DeclareFailed(failedNode *SMNode) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	smlog.Error(LOG_NETWORK, "\033[41m*** dichiaro failed il nodo %d ***\033[0m", failedNode.GetId())
	_, err := cs.ReportAsFailed(ctx, &pb.Node{Id: failedNode.GetId(),
		Host: failedNode.GetHost(),
		Port: failedNode.GetPort()})
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED, err.Error())
	}
}

//TODO unire con la precedente
func DeclareRunning(runningNode *SMNode) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	smlog.Info(LOG_NETWORK, "\033[41m*** dichiaro running il nodo %d ***\033[0m", runningNode.GetId())
	_, err := cs.ReportAsFailed(ctx, &pb.Node{Id: runningNode.GetId(),
		Host: runningNode.GetHost(),
		Port: runningNode.GetPort()})
	if err != nil {
		smlog.Fatal(LOG_UNDEFINED, err.Error())
	}
}
*/
