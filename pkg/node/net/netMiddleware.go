// netMiddleware
package net

//var cs pb.DistGrepClient
import (
	"context"
	pb "fredricksonLynch/pb/node"

	. "fredricksonLynch/pkg/node/env"
	//"fredricksonLynch/pkg/node/statemachine"

	//"fredricksonLynch/pkg/node"

	"flag"
	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"
	"log"
	"net"
	"strconv"
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
	//TODO sistemare, ev. Address etc.
	portParam := flag.String("p", "40043", "porta")
	flag.Parse()
	/*	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}*/
	//port = *portParam
	//addr = "localHost:" + port
	port, _ := strconv.ParseInt(*portParam, 10, 32)
	smlog.Critical(LOG_UNDEFINED, "%d", int32(port))
	Me.SetPort(int32(port))
	Me.SetHost("localhost")
	//addr = "localHost:" + port
	// Parsing input arguments
	/*	filepath := flag.String("f", "../../ILIAD_1STBOOK_IT_ALTERED", "source file to be \"fredricksonLynchp-ed\"")
		substr := flag.String("substr", "Achille", "substr to be searched into the source file")
		serverAddr := flag.String("s", "localHost:40042", "server address and port, in the format ADDRESS:PORT")
		highlight := flag.String("hl", "classic", "[classic/asterisks/none] set substr highlighting in the output\nNOTICE: \"classic\" option may be not available on all systems.")
		help := flag.Bool("help", false, "show this message")

		flag.Parse()
		_, exists := HighlightType[*highlight]
		if !exists {
			fmt.Println("\"-hl\" flag not correctly set.\nSee 'fredricksonLynch -help' for allowed values.")
			os.Exit(-1)
		}
		if *help {
			flag.PrintDefaults()
			os.Exit(0)
		}*/

	// Contacting the server
	// il centrale espone il servizio di identificazione dei nodi

	serverAddr := "localHost:40042" //TODO configurare
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	serverConn = conn
	if err != nil {
		log.Fatalf("Error while contacting server on %v:\n %v", serverAddr, err)
	}
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
		log.Fatalf("Error while trying to listen to port %v:\n%v", Me.GetPort(), err)
	}
	// New server instance and service registering
	w = grpc.NewServer()
	pb.RegisterDistGrepServer(w, &DGnode{})
	// Serve incoming calls

	// Defining client interface, to be used to invoke the fredricksonLynch service
	cs = pb.NewDistGrepClient(serverConn)
}
func Listen() {
	smlog.Info(LOG_NETWORK, "Listening on port %v.", Me.GetPort())
	if err := w.Serve(lis); err != nil {
		log.Fatalf("Error while trying to serve request: %v", err)
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
	serverAddr := "localHost:40042" //TODO configurare
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while contacting server on %v:\n %v", serverAddr, err)
	}
	//defer conn.Close() TODO vedere quando chiudere, altrimenti porta problemi
	return conn
}
func AskForJoining() *SMNode {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	smlog.Info(LOG_SERVREG, "asking for joining the ring...")
	//TODO ma serve ancora NodeAddr?
	//node, err := cs.JoinRing(ctx, &pb.NodeAddr{Addr: Me.GetFullAddr()})
	node, err := cs.JoinRing(ctx, &pb.NodeAddr{Host: Me.GetHost(), Port: Me.GetPort()})
	if err != nil {
		log.Fatalf("Error while executing fredricksonLynch:\n%v", err)
	}
	return ToSMNode(node)
}

//func askForNodeInfo(i int32, forceRunningNode bool) (int32, string) {
func AskForNodeInfo(i int32, forceRunningNode bool) *SMNode {
	smlog.Info(LOG_SERVREG, "Chiedo al centrale informazioni sul nodo %d", i)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//	locCtx = ctx
	defer cancel()
	if forceRunningNode {
		ret, errr := cs.GetNextRunningNode(ctx, &pb.NodeId{Id: int32(i)})
		if errr != nil {
			log.Fatalf("errore in GETNODO:\n%v", errr)
			return nil
		}
		return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	} else {
		ret, errr := cs.GetNode(ctx, &pb.NodeId{Id: int32(i)})
		if errr != nil {
			log.Fatalf("errore in GETNODO:\n%v", errr)
			return nil
		}
		return &SMNode{Id: ret.GetId(), Host: ret.GetHost(), Port: ret.GetPort()}
	}

}

/*func forwardCoordinator(dest *SMNode, msg *MsgCoordinator) {
	coordinatore := msg
	//TODO check su esistenza invalidi?
	safeRMI("C", dest, true, nil, coordinatore, nil)

}*/

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

				// TODO questa di seguito è duplicata rispetto a safeRMI?
				connN, errN := grpc.Dial(node.GetFullAddr(), grpc.WithInsecure())
				if errN != nil {
					log.Printf("Error while contacting server (NODO) on %v:\n %v", node.GetFullAddr(), errN)
				}
				//defer connN.Close() //OSS. NOT DEFERRED!
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
				//_, err = csN.SendHeartBeat(ctx, &pb.Heartbeat{Id: int32(id)})
				//TODO quale err è questo?
				/*if err != nil {
					smlog.Error(LOG_UNDEFINED, "problema nel mandare l'HB: %v", err)
					cs.ReportAsFailed(ctx, node)
				}*/
				//log.Printf("%d", connN.GetState())
				//TODO specificare che ho chiuso qui di proposito
				// per evitare che le connessioni aperte lascino
				// i riceventi in busy su connessioni che non
				// servono più
				connN.Close()

			}
		} // qui ho inviato gli hb a tutti i nodi

		smlog.Info(LOG_HB, "Inviati tutti gli HB, attendo timer...")
		select {
		case <-coordTimer.C:
			//smlog.Printf("SCATTA IL TIMER")
			// semplicemente vai avanti
			break
		case val := <-Events:
			if val == "STOP" {
				log.Printf("------------------ ARRIVATO EVENTO DI \"STOP\"")
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

/*
cs.ReportAsRunning(ctx, &pb.Node{Id: Me.GetId(),
				Host: Me.GetHost(),
				Port: Me.GetPort()})

*/

//TODO vedere come aggiustarla
func generateContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	return ctx, cancel
}
