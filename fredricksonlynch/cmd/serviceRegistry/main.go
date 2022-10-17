// centrale, mantiene id e indirizzi
package main

import (
	"flag"
	"fmt"
	pb "fredricksonLynch/pb/serviceRegistry"
	. "fredricksonLynch/tools/formatting"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

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

func fetchRecordbyId(id int, forceRunningNodeOnly bool) (NodeRecord, bool) {
	// assunzione che i nodi siano identificati a partire da 1
	//TODO uniformare logging basato su loggo
	log.Printf("ricevo richiesta di trovare il nodo %d", id)
	i := id
	searchedNodeWasFailed := false
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
			return nodes[i-1], searchedNodeWasFailed
		} else {
			// whatever node was to be searched, the (first) search
			// resulted in a failed node
			i++
			searchedNodeWasFailed = true
		}

	}
}

func PrintRing() {
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
		log.Printf("%d\t%s\t%s", node.id, node.getFullAddress(), statusStr)
	}
}
func getAllNodesExecutive(forceRunningNodeOnly bool) *pb.NodeList {
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

	// TODO sistemare queste chiamate
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
