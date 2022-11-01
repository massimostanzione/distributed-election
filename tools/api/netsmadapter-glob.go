package api

import (
	pb "distributedelection/serviceregistry/pb"
)

func ToSMNode(net *pb.Node) *SMNode {
	return &SMNode{Id: net.GetId(), Host: net.GetHost(), Port: net.GetPort()}
}

func ToNetNode(sm SMNode) *pb.Node {
	return &pb.Node{Id: sm.GetId(), Host: sm.GetHost(), Port: int32(sm.GetPort())}
}

func ToNetNodeList(list []SMNode) *pb.NodeList {
	var array []*pb.Node
	for i := 0; i < len(list); i++ {
		node := list[i]
		array = append(array, ToNetNode(node))
	}
	return &pb.NodeList{List: array}
}
