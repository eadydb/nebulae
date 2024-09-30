// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.2
// source: v1/nebulae.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NebulaeServiceClient is the client API for NebulaeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NebulaeServiceClient interface {
	HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type nebulaeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNebulaeServiceClient(cc grpc.ClientConnInterface) NebulaeServiceClient {
	return &nebulaeServiceClient{cc}
}

func (c *nebulaeServiceClient) HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/proto.NebulaeService/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NebulaeServiceServer is the server API for NebulaeService service.
// All implementations must embed UnimplementedNebulaeServiceServer
// for forward compatibility
type NebulaeServiceServer interface {
	HealthCheck(context.Context, *emptypb.Empty) (*HealthCheckResponse, error)
	mustEmbedUnimplementedNebulaeServiceServer()
}

// UnimplementedNebulaeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNebulaeServiceServer struct {
}

func (UnimplementedNebulaeServiceServer) HealthCheck(context.Context, *emptypb.Empty) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedNebulaeServiceServer) mustEmbedUnimplementedNebulaeServiceServer() {}

// UnsafeNebulaeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NebulaeServiceServer will
// result in compilation errors.
type UnsafeNebulaeServiceServer interface {
	mustEmbedUnimplementedNebulaeServiceServer()
}

func RegisterNebulaeServiceServer(s grpc.ServiceRegistrar, srv NebulaeServiceServer) {
	s.RegisterService(&NebulaeService_ServiceDesc, srv)
}

func _NebulaeService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NebulaeServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NebulaeService/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NebulaeServiceServer).HealthCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// NebulaeService_ServiceDesc is the grpc.ServiceDesc for NebulaeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NebulaeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.NebulaeService",
	HandlerType: (*NebulaeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _NebulaeService_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/nebulae.proto",
}