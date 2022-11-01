// Network Middleware: adapt gRPC network requests to local registry behavior.
package net

import (
	pb "distributedelection/serviceregistry/pb"
	reg "distributedelection/serviceregistry/pkg/behavior"
	api "distributedelection/tools/api"
)

func ManageJoining(host string, port int32) *pb.Node {
	return api.ToNetNode(reg.RegisterNewNode(host, port))
}

func GetAllNodesExecutive(baseId int32) *pb.NodeList {
	return api.ToNetNodeList(reg.GetNodesWithBaseId(baseId))
}

func FetchRecordById(id int32) *pb.Node {
	return api.ToNetNode(reg.FetchRecordById(int(id)))
}
