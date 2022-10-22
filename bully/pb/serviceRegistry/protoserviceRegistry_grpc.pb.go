// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package serviceRegistry

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DistGrepClient is the client API for DistGrep service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DistGrepClient interface {
	JoinRing(ctx context.Context, in *NodeAddr, opts ...grpc.CallOption) (*Node, error)
	GetNode(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*Node, error)
	GetNextRunningNode(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*Node, error)
	GetAllNodes(ctx context.Context, in *NONE, opts ...grpc.CallOption) (*NodeList, error)
	GetAllRunningNodes(ctx context.Context, in *NONE, opts ...grpc.CallOption) (*NodeList, error)
	GetAllNodesWithIdGreaterThan(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*NodeList, error)
	ReportAsFailed(ctx context.Context, in *Node, opts ...grpc.CallOption) (*NONE, error)
	ReportAsRunning(ctx context.Context, in *Node, opts ...grpc.CallOption) (*NONE, error)
}

type distGrepClient struct {
	cc grpc.ClientConnInterface
}

func NewDistGrepClient(cc grpc.ClientConnInterface) DistGrepClient {
	return &distGrepClient{cc}
}

func (c *distGrepClient) JoinRing(ctx context.Context, in *NodeAddr, opts ...grpc.CallOption) (*Node, error) {
	out := new(Node)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/joinRing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distGrepClient) GetNode(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*Node, error) {
	out := new(Node)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/getNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distGrepClient) GetNextRunningNode(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*Node, error) {
	out := new(Node)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/getNextRunningNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distGrepClient) GetAllNodes(ctx context.Context, in *NONE, opts ...grpc.CallOption) (*NodeList, error) {
	out := new(NodeList)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/getAllNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distGrepClient) GetAllRunningNodes(ctx context.Context, in *NONE, opts ...grpc.CallOption) (*NodeList, error) {
	out := new(NodeList)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/getAllRunningNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distGrepClient) GetAllNodesWithIdGreaterThan(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*NodeList, error) {
	out := new(NodeList)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/getAllNodesWithIdGreaterThan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distGrepClient) ReportAsFailed(ctx context.Context, in *Node, opts ...grpc.CallOption) (*NONE, error) {
	out := new(NONE)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/reportAsFailed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distGrepClient) ReportAsRunning(ctx context.Context, in *Node, opts ...grpc.CallOption) (*NONE, error) {
	out := new(NONE)
	err := c.cc.Invoke(ctx, "/bully.DistGrep/reportAsRunning", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DistGrepServer is the server API for DistGrep service.
// All implementations must embed UnimplementedDistGrepServer
// for forward compatibility
type DistGrepServer interface {
	JoinRing(context.Context, *NodeAddr) (*Node, error)
	GetNode(context.Context, *NodeId) (*Node, error)
	GetNextRunningNode(context.Context, *NodeId) (*Node, error)
	GetAllNodes(context.Context, *NONE) (*NodeList, error)
	GetAllRunningNodes(context.Context, *NONE) (*NodeList, error)
	GetAllNodesWithIdGreaterThan(context.Context, *NodeId) (*NodeList, error)
	ReportAsFailed(context.Context, *Node) (*NONE, error)
	ReportAsRunning(context.Context, *Node) (*NONE, error)
	mustEmbedUnimplementedDistGrepServer()
}

// UnimplementedDistGrepServer must be embedded to have forward compatible implementations.
type UnimplementedDistGrepServer struct {
}

func (UnimplementedDistGrepServer) JoinRing(context.Context, *NodeAddr) (*Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinRing not implemented")
}
func (UnimplementedDistGrepServer) GetNode(context.Context, *NodeId) (*Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNode not implemented")
}
func (UnimplementedDistGrepServer) GetNextRunningNode(context.Context, *NodeId) (*Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNextRunningNode not implemented")
}
func (UnimplementedDistGrepServer) GetAllNodes(context.Context, *NONE) (*NodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNodes not implemented")
}
func (UnimplementedDistGrepServer) GetAllRunningNodes(context.Context, *NONE) (*NodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllRunningNodes not implemented")
}
func (UnimplementedDistGrepServer) GetAllNodesWithIdGreaterThan(context.Context, *NodeId) (*NodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNodesWithIdGreaterThan not implemented")
}
func (UnimplementedDistGrepServer) ReportAsFailed(context.Context, *Node) (*NONE, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportAsFailed not implemented")
}
func (UnimplementedDistGrepServer) ReportAsRunning(context.Context, *Node) (*NONE, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportAsRunning not implemented")
}
func (UnimplementedDistGrepServer) mustEmbedUnimplementedDistGrepServer() {}

// UnsafeDistGrepServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DistGrepServer will
// result in compilation errors.
type UnsafeDistGrepServer interface {
	mustEmbedUnimplementedDistGrepServer()
}

func RegisterDistGrepServer(s grpc.ServiceRegistrar, srv DistGrepServer) {
	s.RegisterService(&DistGrep_ServiceDesc, srv)
}

func _DistGrep_JoinRing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeAddr)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).JoinRing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/joinRing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).JoinRing(ctx, req.(*NodeAddr))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistGrep_GetNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).GetNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/getNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).GetNode(ctx, req.(*NodeId))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistGrep_GetNextRunningNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).GetNextRunningNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/getNextRunningNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).GetNextRunningNode(ctx, req.(*NodeId))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistGrep_GetAllNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NONE)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).GetAllNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/getAllNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).GetAllNodes(ctx, req.(*NONE))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistGrep_GetAllRunningNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NONE)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).GetAllRunningNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/getAllRunningNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).GetAllRunningNodes(ctx, req.(*NONE))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistGrep_GetAllNodesWithIdGreaterThan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).GetAllNodesWithIdGreaterThan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/getAllNodesWithIdGreaterThan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).GetAllNodesWithIdGreaterThan(ctx, req.(*NodeId))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistGrep_ReportAsFailed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).ReportAsFailed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/reportAsFailed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).ReportAsFailed(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _DistGrep_ReportAsRunning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistGrepServer).ReportAsRunning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bully.DistGrep/reportAsRunning",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistGrepServer).ReportAsRunning(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

// DistGrep_ServiceDesc is the grpc.ServiceDesc for DistGrep service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DistGrep_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bully.DistGrep",
	HandlerType: (*DistGrepServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "joinRing",
			Handler:    _DistGrep_JoinRing_Handler,
		},
		{
			MethodName: "getNode",
			Handler:    _DistGrep_GetNode_Handler,
		},
		{
			MethodName: "getNextRunningNode",
			Handler:    _DistGrep_GetNextRunningNode_Handler,
		},
		{
			MethodName: "getAllNodes",
			Handler:    _DistGrep_GetAllNodes_Handler,
		},
		{
			MethodName: "getAllRunningNodes",
			Handler:    _DistGrep_GetAllRunningNodes_Handler,
		},
		{
			MethodName: "getAllNodesWithIdGreaterThan",
			Handler:    _DistGrep_GetAllNodesWithIdGreaterThan_Handler,
		},
		{
			MethodName: "reportAsFailed",
			Handler:    _DistGrep_ReportAsFailed_Handler,
		},
		{
			MethodName: "reportAsRunning",
			Handler:    _DistGrep_ReportAsRunning_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "serviceRegistry/protoserviceRegistry.proto",
}