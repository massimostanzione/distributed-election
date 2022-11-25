// Registry management: the actual registry and all the local operation
// that could be effectively performed on it.
package behavior

import (
	. "distributedelection/tools/api"
	. "distributedelection/tools/smlog"
	smlog "distributedelection/tools/smlog"
	"fmt"
)

// The actual registry.
var registry []SMNode

// Given the host and the port of a node, register it into the registry,
// or return its information if already registered.
func RegisterNewNode(host string, port int32) SMNode {
	node, existent := fetchRecordbyAddr(host, port)
	if existent {
		smlog.Info(LOG_SERVREG, "Node is already in the registry wth id = %d", node.Id)
	} else {
		// register new network node
		registry = append(registry, node)
		smlog.Info(LOG_SERVREG, "Logging new node into the registry with ID = %d", node.Id)
	}
	printRegistry()
	return node
}

// Fetch from the registry information about all the nodes
// with ID greater than baseId
func GetNodesWithBaseId(baseId int32) []SMNode {
	var array []SMNode
	for i := (baseId + 1); i <= int32(len(registry)); i++ {
		node := registry[i-1]
		array = append(array, node)
	}
	return array
}

// Given a node ID, fetch the related information if available
// into the registry, elsewhere generate a new one.
func FetchRecordById(id int) SMNode {
	// Baseline assumptions:
	//  i. nodes are identified starting from ID = 1
	// ii. order in the registry array is based on ID order
	smlog.Trace(LOG_SERVREG, "Request for node with id = %d", id)
	i := id
	for {
		i = i % len(registry)
		if i == 0 {
			i = len(registry)
		}
		return registry[i-1]
	}
}

// Given an hostname and a port, fetch the related information if available
// into the registry, elsewhere generate a new one.
func fetchRecordbyAddr(host string, port int32) (SMNode, bool) {
	smlog.Trace(LOG_SERVREG, "Request for node with host = %s and port = %d", host, port)
	for i := range registry {
		if registry[i].GetFullAddr() == (host + ":" + fmt.Sprint(port)) {
			// if found, the related node is active or is coming back alive
			// after a failure, so the previous ID will be returned back
			return registry[i], true
		}
	}
	return SMNode{getNewId(), host, port}, false
}

// Generate ID for a new node that is being registered.
// Assumption: IDs are assigned incrementally.
func getNewId() int32 {
	return int32(len(registry)) + 1
}

// Print the registry content.
func printRegistry() {
	smlog.Info(LOG_SERVREG, "Printing registry content:")
	smlog.Info(LOG_SERVREG, "id\taddr")
	smlog.Info(LOG_SERVREG, "---\t-------------------")
	for _, node := range registry {
		smlog.Info(LOG_SERVREG, "%d\t%s", node.Id, node.GetFullAddr())
	}
}
