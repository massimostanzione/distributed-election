// centrale, mantiene id e indirizzi
package main

import (
	"context"
	"flag"
	"fmt"
	pb "fredricksonLynch/pb/serviceRegistry"
	. "fredricksonLynch/tools/formatting"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DGserver struct {
	pb.UnimplementedDistGrepServer
}
type NodeRecord struct {

	// TODO documentazione: reportedAsFailed in modo da poter riassegnare stesso numero se il falllimento è temporaneo
	id               int
	Host             string
	Port             int32
	reportedAsFailed bool
}

func (record *NodeRecord) getFullAddress() string {
	return record.Host + ":" + fmt.Sprint(record.Port)
}

var nodes []NodeRecord
var NONE = &pb.NONE{}

func fetchRecordbyAddr(host string, port int32) (NodeRecord, bool) {
	for i := range nodes {
		if nodes[i].getFullAddress() == (host + ":" + fmt.Sprint(port)) {
			// se l'ho trovato è perché mi ha chiesto di entrare,
			// quindi o è nato attivo o è stato attivo prima di fallire,
			// e ora sta rientrando (gli ritorno il posto che aveva prima)
			nodes[i].reportedAsFailed = false
			return nodes[i], true
		}
	}
	return NodeRecord{getNewId(), host, port, false}, false
}

// TODO per tutte le funzioni, ovunque, ritornare anche il booleano di errore
func fetchRecordbyId(id int, forceRunningNodeOnly bool) (NodeRecord, bool, bool) {
	// assunzione che i nodi siano identificati a partire da 1
	//TODO uniformare logging basato su loggo
	log.Printf("ricevo richiesta di trovare il nodo %d", id)
	i := id
	anf := false
	/*normalized := id % len(nodes)
	if normalized == 0 {
		//log.Printf("norm=0")
		normalized = len(nodes)
	}
	*/
	//log.Printf("Node vs normalizzato: %d %d", id, normalized)
	// assumo per ipotesi che l'ordinamento dei nodi nell'array
	// coincida con l'ordinamento degli indici
	// ossia: nodes[i] è il nodo i-esimo
	// TODO implementare ordinamento in base al campo id?
	for {
		//i = (i % len(nodes)==0)?i % len(nodes):0
		i = i % len(nodes)
		if i == 0 {
			i = len(nodes)
		}

		if (forceRunningNodeOnly && !nodes[i-1].reportedAsFailed) || !forceRunningNodeOnly {
			// FIXME il terzo parametro?
			return nodes[i-1], anf, true
		} else {
			i++
			anf = true
		}

	}
	/*
		for _, node := range nodes {
			if node.id == normalized {
				return node, true
			}
		}*/

	log.Fatalf("what? unreachable code")
	//return fetchRecordbyId(id - 1) //TODO gestire, "ritorna se stesso" ma da generalizzare
	return NodeRecord{-99, "wtf", -1, false}, true, false
}

func (s *DGserver) JoinRing(ctx context.Context, in *pb.NodeAddr) (*pb.Node, error) {
	log.Printf("*** REQUEST RECEIVED ***")
	log.Printf("Un nodo vuole entrare, all'indirizzo %s", in.GetHost()+":"+fmt.Sprint(in.GetPort()))

	node, existent := fetchRecordbyAddr(in.GetHost(), in.GetPort())
	if existent {
		log.Printf("Il nodo all'indirizzo %s risulta registrato, con id=%d", in.GetHost()+":"+fmt.Sprint(in.GetPort()), node.id)
	} else {
		nodes = append(nodes, node)
		log.Printf("Il nodo mi è nuovo, gli assegno id=%d", node.id)
	}
	printRing()
	//compreso map della struct locale in quella grpc-compatibile:
	return &pb.Node{Id: int32(node.id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}

func (s *DGserver) ReportAsFailed(ctx context.Context, in *pb.Node) (*pb.NONE, error) {
	log.Printf("Il nodo n. %d, presso %s, mi è stato segnalato come failed", in.GetHost()+":"+fmt.Sprint(in.GetPort()), in.GetHost()+":"+fmt.Sprint(in.GetPort()))
	for i := range nodes {
		if int32(nodes[i].id) == in.GetId() {
			nodes[i].reportedAsFailed = true
			printRing()
			return NONE, status.New(codes.OK, "").Err()
		}
	}
	//TODO ritornare un errore: non sono riuscito a trovare il nodo da marcare come fallito
	log.Fatalf("NON SONO RIUSCITO A TROVARE IL NODO %d, da marcare come FAILED", in.GetId())
	return NONE, status.New(codes.OK, "").Err()
}

// TODO documentazione: può succedere che un nodo veda un altro temporaneamente out,
// quindi lo segnala come tale, e al ricevere di un messaggio da tale nodo,
// non risponde in quanto il server gli dice ancora che è out, quindi bisogna aggiornare
func (s *DGserver) ReportAsRunning(ctx context.Context, in *pb.Node) (*pb.NONE, error) {
	log.Printf("Il nodo n. %d, presso %s, mi è stato segnalato come RUNNING", in.GetId(), in.GetHost()+":"+fmt.Sprint(in.GetPort()))
	for i := range nodes {
		if int32(nodes[i].id) == in.GetId() {
			nodes[i].reportedAsFailed = false
			printRing()
			return NONE, status.New(codes.OK, "").Err()
		}
	}
	//TODO ritornare un errore: non sono riuscito a trovare il nodo da marcare come funzionante
	log.Fatalf("NON SONO RIUSCITO A TROVARE IL NODO %d, da marcare come RUNNING", in.GetId())
	return NONE, status.New(codes.OK, "").Err()
}
func printRing() {
	log.Printf("L'anello adesso è fatto così:")
	log.Printf("id\taddr\t\tstatus")
	log.Printf("---\t---------------\t---------")
	for _, node := range nodes {
		statusStr := "N.D."
		if node.reportedAsFailed {
			statusStr = ColorRed + Bold + "FAILED" + ColorReset
		} else {
			statusStr = ColorGreen + Bold + "RUNNING" + ColorReset
		}
		//TODO implementare anche qui un getFullAddr?
		log.Printf("%d\t%s\t%s", node.id, node.Host+":"+fmt.Sprint(node.Port), statusStr)
	}
}
func (s *DGserver) GetNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	log.Printf("*** REQUEST RECEIVED ***")
	log.Printf("Serve conoscere chi è %d", in.Id)
	//anf := false
	// TODO i tre parametri ritornati
	node, _, _ := fetchRecordbyId(int(in.Id), false) //TODO ricostruire catene di cast
	// TODO controllo su err
	/*if err != nil {
		log.Fatalf("gestire")
		return NONE, false
	}*/
	log.Printf("Ti ritorno il nodo richiesto, che è %s", node.Host+":"+fmt.Sprint(node.Port))
	printRing()
	return &pb.Node{Id: int32(node.id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}

func (s *DGserver) GetNextRunningNode(ctx context.Context, in *pb.NodeId) (*pb.Node, error) {
	log.Printf("*** REQUEST RECEIVED ***")
	log.Printf("Serve conoscere chi è %d", in.Id)
	node, _, _ := fetchRecordbyId(int(in.Id), true)
	// TODO controllo su err
	/*if err != nil {
		log.Fatalf("gestire")
		return NONE, false
	}*/
	log.Printf("Ti ritorno il nodo richiesto, che è %s", node.Host+":"+fmt.Sprint(node.Port))
	printRing()
	return &pb.Node{Id: int32(node.id), Host: node.Host, Port: node.Port}, status.New(codes.OK, "").Err()
}
func ggetAllNodes(forceRunningNodeOnly bool) *pb.NodeList {
	var array []*pb.Node
	for _, node := range nodes {
		// TODO documentazione: qui è il centrale che si occupa di sapere lo stato di ciascuno,
		// anche perché è lo stesso che mantiene la struttura dell'anello,
		// quindi è anche giusto che se ne occupi lui che distribuire il tutto
		// in modo distribuito
		if (forceRunningNodeOnly && !node.reportedAsFailed) || !forceRunningNodeOnly {
			//TODO funzione di mappatura/demappatura da grpc.Node a Node di go, locale qui
			grpcNode := &pb.Node{Id: int32(node.id), Host: node.Host, Port: node.Port}
			array = append(array, grpcNode)
		}
	}
	return &pb.NodeList{List: array}

}
func (s *DGserver) GetAllNodes(ctx context.Context, in *pb.NONE) (*pb.NodeList, error) {
	/*var array []*pb.Node
	for _, node := range nodes {
		// TODO documentazione: qui è il centrale che si occupa di sapere lo stato di ciascuno,
		// anche perché è lo stesso che mantiene la struttura dell'anello,
		// quindi è anche giusto che se ne occupi lui che distribuire il tutto
		// in modo distribuito
		if node.reportedAsFailed == false {
			//TODO funzione di mappatura/demappatura da grpc.Node a Node di go, locale qui
			grpcNode := &pb.Node{Id: int32(node.id), Addr: node.addr}
			array = append(array, grpcNode)
		}
	}*/
	return ggetAllNodes(false), status.New(codes.OK, "").Err()
}
func (s *DGserver) GetAllRunningNodes(ctx context.Context, in *pb.NONE) (*pb.NodeList, error) {
	// copiato dal precedente
	/*	var array []*pb.Node
		for _, node := range nodes {
			// TODO documentazione: qui è il centrale che si occupa di sapere lo stato di ciascuno,
			// anche perché è lo stesso che mantiene la struttura dell'anello,
			// quindi è anche giusto che se ne occupi lui che distribuire il tutto
			// in modo distribuito
			if node.reportedAsFailed == false {
				//TODO funzione di mappatura/demappatura da grpc.Node a Node di go, locale qui
				grpcNode := &pb.Node{Id: int32(node.id), Addr: node.addr}
				array = append(array, grpcNode)
			}
		}
		return &pb.NodeList{List: array}, status.New(codes.OK, "").Err()*/

	return ggetAllNodes(true), status.New(codes.OK, "").Err()
}
func getNewId() int {
	return len(nodes) + 1
}
func main() {
	log.Println("*** DISTGREP SERVER ***")
	// Parsing input arguments
	port := flag.String("p", "40042", "port to listen for distgrep requests")
	//workers := flag.String("w", "localhost:40043;localhost:40044;localhost:40045", "addresses and ports of the workers to be bound with, in the following format:\nADDRESS_1:PORT_1;ADDRESS2:PORT_2;...;ADDRESS_N:PORT_N\nMust be between 1 and 15")
	help := flag.Bool("help", false, "show this message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	/*
		workersArray = strings.Split(*workers, ";")
		if len(workersArray) < MIN_WORKERSNO || len(workersArray) > MAX_WORKERSNO {
			fmt.Printf("Workers to be bound with must be between %v and %v\n", MIN_WORKERSNO, MAX_WORKERSNO)
			os.Exit(-1)
		}
		log.Println("Will bind (on-demand) to the following workers:")
		for i := range workersArray {
			log.Printf("- %v", workersArray[i])
		}*/
	// Start listening for incoming calls
	lis, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatalf("Error while trying to listen to port %v:\n%v", *port, err)
	}
	log.Printf("--------------------------")
	log.Printf("Listening on port %v...", *port)
	// New server instance and service registering
	s := grpc.NewServer()
	pb.RegisterDistGrepServer(s, &DGserver{})
	// Serve incoming calls
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while trying to serve request: %v", err)
	}
}
