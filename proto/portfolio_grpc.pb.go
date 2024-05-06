// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: proto/portfolio.proto

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

const (
	PortfolioService_CreateAsset_FullMethodName = "/portfolio.PortfolioService/CreateAsset"
	PortfolioService_GetAsset_FullMethodName    = "/portfolio.PortfolioService/GetAsset"
	PortfolioService_UpdateAsset_FullMethodName = "/portfolio.PortfolioService/UpdateAsset"
	PortfolioService_DeleteAsset_FullMethodName = "/portfolio.PortfolioService/DeleteAsset"
)

// PortfolioServiceClient is the client API for PortfolioService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PortfolioServiceClient interface {
	CreateAsset(ctx context.Context, in *CreateAssetRequest, opts ...grpc.CallOption) (*Asset, error)
	GetAsset(ctx context.Context, in *GetAssetRequest, opts ...grpc.CallOption) (*Asset, error)
	UpdateAsset(ctx context.Context, in *UpdateAssetRequest, opts ...grpc.CallOption) (*Asset, error)
	DeleteAsset(ctx context.Context, in *DeleteAssetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type portfolioServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPortfolioServiceClient(cc grpc.ClientConnInterface) PortfolioServiceClient {
	return &portfolioServiceClient{cc}
}

func (c *portfolioServiceClient) CreateAsset(ctx context.Context, in *CreateAssetRequest, opts ...grpc.CallOption) (*Asset, error) {
	out := new(Asset)
	err := c.cc.Invoke(ctx, PortfolioService_CreateAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portfolioServiceClient) GetAsset(ctx context.Context, in *GetAssetRequest, opts ...grpc.CallOption) (*Asset, error) {
	out := new(Asset)
	err := c.cc.Invoke(ctx, PortfolioService_GetAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portfolioServiceClient) UpdateAsset(ctx context.Context, in *UpdateAssetRequest, opts ...grpc.CallOption) (*Asset, error) {
	out := new(Asset)
	err := c.cc.Invoke(ctx, PortfolioService_UpdateAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portfolioServiceClient) DeleteAsset(ctx context.Context, in *DeleteAssetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PortfolioService_DeleteAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PortfolioServiceServer is the server API for PortfolioService service.
// All implementations must embed UnimplementedPortfolioServiceServer
// for forward compatibility
type PortfolioServiceServer interface {
	CreateAsset(context.Context, *CreateAssetRequest) (*Asset, error)
	GetAsset(context.Context, *GetAssetRequest) (*Asset, error)
	UpdateAsset(context.Context, *UpdateAssetRequest) (*Asset, error)
	DeleteAsset(context.Context, *DeleteAssetRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedPortfolioServiceServer()
}

// UnimplementedPortfolioServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPortfolioServiceServer struct {
}

func (UnimplementedPortfolioServiceServer) CreateAsset(context.Context, *CreateAssetRequest) (*Asset, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAsset not implemented")
}
func (UnimplementedPortfolioServiceServer) GetAsset(context.Context, *GetAssetRequest) (*Asset, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAsset not implemented")
}
func (UnimplementedPortfolioServiceServer) UpdateAsset(context.Context, *UpdateAssetRequest) (*Asset, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAsset not implemented")
}
func (UnimplementedPortfolioServiceServer) DeleteAsset(context.Context, *DeleteAssetRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAsset not implemented")
}
func (UnimplementedPortfolioServiceServer) mustEmbedUnimplementedPortfolioServiceServer() {}

// UnsafePortfolioServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PortfolioServiceServer will
// result in compilation errors.
type UnsafePortfolioServiceServer interface {
	mustEmbedUnimplementedPortfolioServiceServer()
}

func RegisterPortfolioServiceServer(s grpc.ServiceRegistrar, srv PortfolioServiceServer) {
	s.RegisterService(&PortfolioService_ServiceDesc, srv)
}

func _PortfolioService_CreateAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioServiceServer).CreateAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioService_CreateAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioServiceServer).CreateAsset(ctx, req.(*CreateAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortfolioService_GetAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioServiceServer).GetAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioService_GetAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioServiceServer).GetAsset(ctx, req.(*GetAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortfolioService_UpdateAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioServiceServer).UpdateAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioService_UpdateAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioServiceServer).UpdateAsset(ctx, req.(*UpdateAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortfolioService_DeleteAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioServiceServer).DeleteAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioService_DeleteAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioServiceServer).DeleteAsset(ctx, req.(*DeleteAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PortfolioService_ServiceDesc is the grpc.ServiceDesc for PortfolioService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PortfolioService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "portfolio.PortfolioService",
	HandlerType: (*PortfolioServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAsset",
			Handler:    _PortfolioService_CreateAsset_Handler,
		},
		{
			MethodName: "GetAsset",
			Handler:    _PortfolioService_GetAsset_Handler,
		},
		{
			MethodName: "UpdateAsset",
			Handler:    _PortfolioService_UpdateAsset_Handler,
		},
		{
			MethodName: "DeleteAsset",
			Handler:    _PortfolioService_DeleteAsset_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/portfolio.proto",
}
