package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	proto "goservice/pkg/api/grpc"
	"net"
)

type clients struct {
}

type instance struct {
	proto.UnimplementedGoServiceServer
	clients *clients
}

func (i *instance) ServiceMethod(ctx context.Context, req *proto.ServiceMethodRequest) (*proto.ServiceMethodResponse, error) {
	return &proto.ServiceMethodResponse{Message: req.Message}, nil
}

func Serve(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	srv := &instance{clients: &clients{}}

	grpcServer := grpc.NewServer()
	proto.RegisterGoServiceServer(grpcServer, srv)

	return grpcServer.Serve(l)
}
