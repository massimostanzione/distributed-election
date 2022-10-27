package net

import (
	pb "distributedelection/node/pb"
	. "distributedelection/node/pkg/env"
)

func ToSMNode(net *pb.Node) *SMNode {
	return &SMNode{Id: net.GetId(), Host: net.GetHost(), Port: net.GetPort()}
}

func ToNetNode(sm SMNode) *pb.Node {
	return &pb.Node{Id: sm.GetId(), Host: sm.GetHost(), Port: int32(sm.GetPort())}
}
