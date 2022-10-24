// registrymgt.go
package registrymgt

import (
	"fmt"
	pb "fredricksonLynch/pb/serviceRegistry"
	. "fredricksonLynch/pkg/serviceRegistry/env"

	//	. "fredricksonLynch/tools/smlog"
	smlog "fredricksonLynch/tools/smlog"

	//. "fredricksonLynch/pkg/serviceRegistry/net"
	. "fredricksonLynch/tools/formatting"
	//"log"
)

/*
func StartServiceRegistry() {
	InitializeNetMW()
	//go run()
	Listen()
}*/

func FetchRecordbyAddr(host string, port int32) (NodeRecord, bool) {
	for i := range Nodes {
		if Nodes[i].GetFullAddress() == (host + ":" + fmt.Sprint(port)) {
			// se l'ho trovato è perché mi ha chiesto di entrare,
			// quindi o è nato attivo o è stato attivo prima di fallire,
			// e ora sta rientrando (gli ritorno il posto che aveva prima)

			//Nodes[i].ReportedAsFailed = false
			return Nodes[i], true
		}
	}
	return NodeRecord{getNewId(), host, port, false}, false
}

func FetchRecordbyId(id int, forceRunningNodeOnly bool) (NodeRecord, bool) {
	// assunzione che i nodi siano identificati a partire da 1
	smlog.InfoU("ricevo richiesta di trovare il nodo %d", id)
	i := id
	searchedNodeWasFailed := false
	/*normalized := id % len(Nodes)
	if normalized == 0 {
		//smlog.InfoU("norm=0")
		normalized = len(Nodes)
	}
	*/
	//smlog.InfoU("Node vs normalizzato: %d %d", id, normalized)
	// assumo per ipotesi che l'ordinamento dei nodi nell'array
	// coincida con l'ordinamento degli indici
	// ossia: Nodes[i] è il nodo i-esimo
	for {
		//i = (i % len(Nodes)==0)?i % len(Nodes):0
		i = i % len(Nodes)
		if i == 0 {
			i = len(Nodes)
		}

		if (forceRunningNodeOnly && !Nodes[i-1].ReportedAsFailed) || !forceRunningNodeOnly {
			return Nodes[i-1], searchedNodeWasFailed
		} else {
			// whatever node was to be searched, the (first) search
			// resulted in a failed node
			i++
			searchedNodeWasFailed = true
		}

	}
}

func GetAllNodesExecutive(forceRunningNodeOnly bool) *pb.NodeList {
	var array []*pb.Node
	for _, node := range Nodes {
		// TODO documentazione: qui è il centrale che si occupa di sapere lo stato di ciascuno,
		// anche perché è lo stesso che mantiene la struttura dell'anello,
		// quindi è anche giusto che se ne occupi lui che distribuire il tutto
		// in modo distribuito
		if (forceRunningNodeOnly && !node.ReportedAsFailed) || !forceRunningNodeOnly {
			grpcNode := ToNetNode(node) //&pb.Node{Id: int32(node.Id), Host: node.Host, Port: node.Port}
			array = append(array, grpcNode)
		}
	}
	return &pb.NodeList{List: array}

}
func getNewId() int {
	return len(Nodes) + 1
}

func PrintRing() {
	smlog.InfoU("L'anello adesso è fatto così:")
	smlog.InfoU("id\taddr\t\t\tstatus")
	smlog.InfoU("---\t-------------------\t---------")
	for _, node := range Nodes {
		statusStr := "N.D."
		if node.ReportedAsFailed {
			statusStr = ColorRed + Bold + "FAILED" + ColorReset
		} else {
			statusStr = ColorGreen + Bold + "RUNNING" + ColorReset
		}
		smlog.InfoU("%d\t%s\t%s", node.Id, node.GetFullAddress(), statusStr)
	}
}
