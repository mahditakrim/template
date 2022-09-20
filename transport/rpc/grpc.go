package rpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/mahditakrim/template/service"
	"github.com/mahditakrim/template/transport"
	"github.com/mahditakrim/template/transport/rpc/pb"
	"google.golang.org/grpc"
)

type rpc struct {
	grpc    *grpc.Server
	service service.Library
	pb.UnimplementedLibraryServiceServer
}

func NewRPC(service service.Library) (transport.Transport, error) {

	if service == nil {
		return nil, errors.New("nil service reference")
	}

	return &rpc{grpc.NewServer(), service, pb.UnimplementedLibraryServiceServer{}}, nil
}

func (rpc *rpc) Run(addr string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	pb.RegisterLibraryServiceServer(rpc.grpc, rpc)
	fmt.Println("rpc server listening on", addr)

	return rpc.grpc.Serve(listener)
}

func (rpc *rpc) Shutdown() error {

	rpc.grpc.GracefulStop()
	return nil
}
