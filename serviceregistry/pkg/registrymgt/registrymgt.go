// registrymgt.go
package registrymgt

import (
	pb "distributedelection/serviceregistry/pb"
	. "distributedelection/serviceregistry/pkg/env"

	//	. "distributedelection/serviceregistry/tools"
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"fmt"
	//	. "distributedelection/tools/formatting"
)

// Given an hostname and a port, fetch the related SMNode if available
// into the registry, elsewhere generate a new one.
func FetchRecordbyAddr(host string, port int32) (SMNode, bool) {
	smlog.Trace(LOG_SERVREG, "Request for node with host = %s and port = %d", host, port)
	for i := range Nodes {
		if Nodes[i].GetFullAddr() == (host + ":" + fmt.Sprint(port)) {
			// if found, the related node is active or is coming back alive
			// after a failure, so the previous ID will be returned back
			return Nodes[i], true
		}
	}
	return SMNode{getNewId(), host, port}, false
}

// Given a node ID, fetch the related SMNode if available
// into the registry, elsewhere generate a new one.
func FetchRecordbyId(id int) SMNode {
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

func getNewId() int32 {
	return int32(len(Nodes)) + 1
}

func printRegistry() {
	smlog.Info(LOG_SERVREG, "Printing registry content:")
	smlog.Info(LOG_SERVREG, "id\taddr\t\t\tstatus")
	smlog.Info(LOG_SERVREG, "---\t-------------------\t---------")
	for _, node := range Nodes {
		smlog.Info(LOG_SERVREG, "%d\t%s", node.Id, node.GetFullAddr())
	}
}
func ManageJoining(host string, port int32) SMNode {
	node, existent := FetchRecordbyAddr(host, port)
	if existent {
		smlog.Info(LOG_SERVREG, "Node is already in the registry wth id = %d", node.Id)
	} else {
		// register new network node
		Nodes = append(Nodes, node)
		smlog.Info(LOG_SERVREG, "Logging new node into the registry with ID = %d", node.Id)
	}
	printRegistry()
	return node
}
