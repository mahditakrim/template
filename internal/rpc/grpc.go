package rpc

import (
	"errors"
	"net"

	"github.com/mahditakrim/template/internal/rpc/pb"
	"github.com/mahditakrim/template/luncher"
	"github.com/mahditakrim/template/service"
	"google.golang.org/grpc"
)

type rpc struct {
	listener net.Listener
	grpc     *grpc.Server
	service  service.Service
	pb.UnimplementedLibraryServiceServer
}

func NewRPC(service service.Service, addr string) (luncher.Runnable, error) {

	if service == nil {
		return nil, errors.New("nil service reference")
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &rpc{
		listener,
		grpc.NewServer(),
		service,
		pb.UnimplementedLibraryServiceServer{},
	}, nil
}

func (rpc *rpc) Run() error {

	pb.RegisterLibraryServiceServer(rpc.grpc, rpc)
	return rpc.grpc.Serve(rpc.listener)
}

func (rpc *rpc) Shutdown() error {

	rpc.grpc.GracefulStop()
	return nil
}
