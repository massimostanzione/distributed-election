package tools

import (
	pb "bully/serviceregistry/pb"
	. "bully/serviceregistry/pkg/env"
)

func ToNetNode(sm NodeRecord) *pb.Node {
	return &pb.Node{Id: int32(sm.GetId()), Host: sm.GetHost(), Port: int32(sm.GetPort())}
}
