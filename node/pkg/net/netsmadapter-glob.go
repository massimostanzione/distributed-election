package net

import (
	. "distributedelection/node/pkg/env"
	pb "distributedelection/serviceregistry/pb"
)

func ToSMNode(net *pb.Node) *SMNode {
	return &SMNode{Id: net.GetId(), Host: net.GetHost(), Port: net.GetPort()}
}

func ToNetNode(sm SMNode) *pb.Node {
	return &pb.Node{Id: sm.GetId(), Host: sm.GetHost(), Port: int32(sm.GetPort())}
}
