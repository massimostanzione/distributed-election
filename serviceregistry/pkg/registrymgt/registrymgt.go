// registrymgt.go
package registrymgt

import (
	pb "distributedelection/serviceregistry/pb"
	. "distributedelection/serviceregistry/pkg/env"
	. "distributedelection/serviceregistry/tools"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"fmt"
	//	. "distributedelection/tools/formatting"
)

// Given an hostname and a port, fetch the related NodeRecord if available
// into the registry, elsewhere generate a new one.
func FetchRecordbyAddr(host string, port int32) (NodeRecord, bool) {
	smlog.Trace(LOG_SERVREG, "Request for node with host = %s and port = %d", host, port)
	for i := range Nodes {
		if Nodes[i].GetFullAddress() == (host + ":" + fmt.Sprint(port)) {
			// if found, the related node is active or is coming back alive
			// after a failure, so the previous ID will be returned back
			return Nodes[i], true
		}
	}
	return NodeRecord{getNewId(), host, port}, false
}

// Given a node ID, fetch the related NodeRecord if available
// into the registry, elsewhere generate a new one.
func FetchRecordbyId(id int) NodeRecord {
	// Baseline assumptions:
	//  i. nodes are identified starting from ID = 1
	// ii. order in the registry array is based on ID order
	smlog.Trace(LOG_SERVREG, "Request for node with id = %d", id)
	i := id
	for {
		i = i % len(Nodes)
		if i == 0 {
			i = len(Nodes)
		}

		//if (forceRunningNodeOnly && !Nodes[i-1].ReportedAsFailed) || !forceRunningNodeOnly {
		return Nodes[i-1] //, searchedNodeWasFailed
		//} else {
		// whatever node was to be searched, the (first) search
		// resulted in a failed node
		//	i++
		//	searchedNodeWasFailed = true
		//}

	}
}

func GetAllNodesExecutive(baseId int32) *pb.NodeList {
	var array []*pb.Node
	for i := (baseId + 1); i <= int32(len(Nodes)); i++ {
		node := Nodes[i-1]
		grpcNode := ToNetNode(node)
		array = append(array, grpcNode)
	}
	return &pb.NodeList{List: array}
}

func getNewId() int {
	return len(Nodes) + 1
}

func PrintRing() {
	smlog.Info(LOG_SERVREG, "Printing registry content:")
	smlog.Info(LOG_SERVREG, "id\taddr\t\t\tstatus")
	smlog.Info(LOG_SERVREG, "---\t-------------------\t---------")
	for _, node := range Nodes {
		smlog.Info(LOG_SERVREG, "%d\t%s", node.Id, node.GetFullAddress())
	}
}
