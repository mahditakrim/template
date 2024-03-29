// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: internal/rpc/pb/library.proto

package pb

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

// LibraryServiceClient is the client API for LibraryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LibraryServiceClient interface {
	GetBook(ctx context.Context, in *BookID, opts ...grpc.CallOption) (*Book, error)
	GetBooks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Books, error)
	EditeBook(ctx context.Context, in *Book, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateBook(ctx context.Context, in *Book, opts ...grpc.CallOption) (*BookID, error)
	RemoveBook(ctx context.Context, in *BookID, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type libraryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLibraryServiceClient(cc grpc.ClientConnInterface) LibraryServiceClient {
	return &libraryServiceClient{cc}
}

func (c *libraryServiceClient) GetBook(ctx context.Context, in *BookID, opts ...grpc.CallOption) (*Book, error) {
	out := new(Book)
	err := c.cc.Invoke(ctx, "/pb.LibraryService/GetBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *libraryServiceClient) GetBooks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Books, error) {
	out := new(Books)
	err := c.cc.Invoke(ctx, "/pb.LibraryService/GetBooks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *libraryServiceClient) EditeBook(ctx context.Context, in *Book, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/pb.LibraryService/EditeBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *libraryServiceClient) CreateBook(ctx context.Context, in *Book, opts ...grpc.CallOption) (*BookID, error) {
	out := new(BookID)
	err := c.cc.Invoke(ctx, "/pb.LibraryService/CreateBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *libraryServiceClient) RemoveBook(ctx context.Context, in *BookID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/pb.LibraryService/RemoveBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LibraryServiceServer is the server API for LibraryService service.
// All implementations must embed UnimplementedLibraryServiceServer
// for forward compatibility
type LibraryServiceServer interface {
	GetBook(context.Context, *BookID) (*Book, error)
	GetBooks(context.Context, *emptypb.Empty) (*Books, error)
	EditeBook(context.Context, *Book) (*emptypb.Empty, error)
	CreateBook(context.Context, *Book) (*BookID, error)
	RemoveBook(context.Context, *BookID) (*emptypb.Empty, error)
	mustEmbedUnimplementedLibraryServiceServer()
}

// UnimplementedLibraryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLibraryServiceServer struct {
}

func (UnimplementedLibraryServiceServer) GetBook(context.Context, *BookID) (*Book, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBook not implemented")
}
func (UnimplementedLibraryServiceServer) GetBooks(context.Context, *emptypb.Empty) (*Books, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBooks not implemented")
}
func (UnimplementedLibraryServiceServer) EditeBook(context.Context, *Book) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditeBook not implemented")
}
func (UnimplementedLibraryServiceServer) CreateBook(context.Context, *Book) (*BookID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBook not implemented")
}
func (UnimplementedLibraryServiceServer) RemoveBook(context.Context, *BookID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveBook not implemented")
}
func (UnimplementedLibraryServiceServer) mustEmbedUnimplementedLibraryServiceServer() {}

// UnsafeLibraryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LibraryServiceServer will
// result in compilation errors.
type UnsafeLibraryServiceServer interface {
	mustEmbedUnimplementedLibraryServiceServer()
}

func RegisterLibraryServiceServer(s grpc.ServiceRegistrar, srv LibraryServiceServer) {
	s.RegisterService(&LibraryService_ServiceDesc, srv)
}

func _LibraryService_GetBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LibraryServiceServer).GetBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LibraryService/GetBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LibraryServiceServer).GetBook(ctx, req.(*BookID))
	}
	return interceptor(ctx, in, info, handler)
}

func _LibraryService_GetBooks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LibraryServiceServer).GetBooks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LibraryService/GetBooks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LibraryServiceServer).GetBooks(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LibraryService_EditeBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Book)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LibraryServiceServer).EditeBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LibraryService/EditeBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LibraryServiceServer).EditeBook(ctx, req.(*Book))
	}
	return interceptor(ctx, in, info, handler)
}

func _LibraryService_CreateBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Book)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LibraryServiceServer).CreateBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LibraryService/CreateBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LibraryServiceServer).CreateBook(ctx, req.(*Book))
	}
	return interceptor(ctx, in, info, handler)
}

func _LibraryService_RemoveBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LibraryServiceServer).RemoveBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LibraryService/RemoveBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LibraryServiceServer).RemoveBook(ctx, req.(*BookID))
	}
	return interceptor(ctx, in, info, handler)
}

// LibraryService_ServiceDesc is the grpc.ServiceDesc for LibraryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LibraryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.LibraryService",
	HandlerType: (*LibraryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBook",
			Handler:    _LibraryService_GetBook_Handler,
		},
		{
			MethodName: "GetBooks",
			Handler:    _LibraryService_GetBooks_Handler,
		},
		{
			MethodName: "EditeBook",
			Handler:    _LibraryService_EditeBook_Handler,
		},
		{
			MethodName: "CreateBook",
			Handler:    _LibraryService_CreateBook_Handler,
		},
		{
			MethodName: "RemoveBook",
			Handler:    _LibraryService_RemoveBook_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/rpc/pb/library.proto",
}
