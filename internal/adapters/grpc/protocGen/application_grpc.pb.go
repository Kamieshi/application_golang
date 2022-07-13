// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: application.proto

package protocGen

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

// EntityClient is the client API for Entity service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EntityClient interface {
	GetEntityById(ctx context.Context, in *GetEntityByIdRequest, opts ...grpc.CallOption) (*GetEntityByIdResponse, error)
	GetAllEntity(ctx context.Context, in *GetAllEntityRequest, opts ...grpc.CallOption) (*GetAllEntityResponse, error)
	UpdateEntity(ctx context.Context, in *UpdateEntityRequest, opts ...grpc.CallOption) (*UpdateEntityResponse, error)
	DeleteEntity(ctx context.Context, in *DeleteEntityRequest, opts ...grpc.CallOption) (*DeleteEntityResponse, error)
}

type entityClient struct {
	cc grpc.ClientConnInterface
}

func NewEntityClient(cc grpc.ClientConnInterface) EntityClient {
	return &entityClient{cc}
}

func (c *entityClient) GetEntityById(ctx context.Context, in *GetEntityByIdRequest, opts ...grpc.CallOption) (*GetEntityByIdResponse, error) {
	out := new(GetEntityByIdResponse)
	err := c.cc.Invoke(ctx, "/applicationGolang.Entity/GetEntityById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityClient) GetAllEntity(ctx context.Context, in *GetAllEntityRequest, opts ...grpc.CallOption) (*GetAllEntityResponse, error) {
	out := new(GetAllEntityResponse)
	err := c.cc.Invoke(ctx, "/applicationGolang.Entity/GetAllEntity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityClient) UpdateEntity(ctx context.Context, in *UpdateEntityRequest, opts ...grpc.CallOption) (*UpdateEntityResponse, error) {
	out := new(UpdateEntityResponse)
	err := c.cc.Invoke(ctx, "/applicationGolang.Entity/UpdateEntity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *entityClient) DeleteEntity(ctx context.Context, in *DeleteEntityRequest, opts ...grpc.CallOption) (*DeleteEntityResponse, error) {
	out := new(DeleteEntityResponse)
	err := c.cc.Invoke(ctx, "/applicationGolang.Entity/DeleteEntity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EntityServer is the server API for Entity service.
// All implementations must embed UnimplementedEntityServer
// for forward compatibility
type EntityServer interface {
	GetEntityById(context.Context, *GetEntityByIdRequest) (*GetEntityByIdResponse, error)
	GetAllEntity(context.Context, *GetAllEntityRequest) (*GetAllEntityResponse, error)
	UpdateEntity(context.Context, *UpdateEntityRequest) (*UpdateEntityResponse, error)
	DeleteEntity(context.Context, *DeleteEntityRequest) (*DeleteEntityResponse, error)
	mustEmbedUnimplementedEntityServer()
}

// UnimplementedEntityServer must be embedded to have forward compatible implementations.
type UnimplementedEntityServer struct {
}

func (UnimplementedEntityServer) GetEntityById(context.Context, *GetEntityByIdRequest) (*GetEntityByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEntityById not implemented")
}
func (UnimplementedEntityServer) GetAllEntity(context.Context, *GetAllEntityRequest) (*GetAllEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllEntity not implemented")
}
func (UnimplementedEntityServer) UpdateEntity(context.Context, *UpdateEntityRequest) (*UpdateEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEntity not implemented")
}
func (UnimplementedEntityServer) DeleteEntity(context.Context, *DeleteEntityRequest) (*DeleteEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEntity not implemented")
}
func (UnimplementedEntityServer) mustEmbedUnimplementedEntityServer() {}

// UnsafeEntityServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EntityServer will
// result in compilation errors.
type UnsafeEntityServer interface {
	mustEmbedUnimplementedEntityServer()
}

func RegisterEntityServer(s grpc.ServiceRegistrar, srv EntityServer) {
	s.RegisterService(&Entity_ServiceDesc, srv)
}

func _Entity_GetEntityById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEntityByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityServer).GetEntityById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/applicationGolang.Entity/GetEntityById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityServer).GetEntityById(ctx, req.(*GetEntityByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Entity_GetAllEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityServer).GetAllEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/applicationGolang.Entity/GetAllEntity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityServer).GetAllEntity(ctx, req.(*GetAllEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Entity_UpdateEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityServer).UpdateEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/applicationGolang.Entity/UpdateEntity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityServer).UpdateEntity(ctx, req.(*UpdateEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Entity_DeleteEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityServer).DeleteEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/applicationGolang.Entity/DeleteEntity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityServer).DeleteEntity(ctx, req.(*DeleteEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Entity_ServiceDesc is the grpc.ServiceDesc for Entity service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Entity_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "applicationGolang.Entity",
	HandlerType: (*EntityServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEntityById",
			Handler:    _Entity_GetEntityById_Handler,
		},
		{
			MethodName: "GetAllEntity",
			Handler:    _Entity_GetAllEntity_Handler,
		},
		{
			MethodName: "UpdateEntity",
			Handler:    _Entity_UpdateEntity_Handler,
		},
		{
			MethodName: "DeleteEntity",
			Handler:    _Entity_DeleteEntity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "application.proto",
}

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "applicationGolang.User",
	HandlerType: (*UserServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "application.proto",
}
