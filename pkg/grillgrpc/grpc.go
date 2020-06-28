package grillgrpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GRPC struct {
	server   *grpc.Server
	recorder *recorder
	host     string
	port     string
}

func (gg *GRPC) Start(ctx context.Context) error {
	recorder := newGRPCRecorder()
	server := grpc.NewServer(grpc.UnaryInterceptor(recorder.unaryInterceptor))

	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	go server.Serve(listen)

	port := listen.Addr().(*net.TCPAddr).Port

	gg.server = server
	gg.host = "localhost"
	gg.port = fmt.Sprintf("%d", port)
	gg.recorder = recorder

	return nil
}

func (gg *GRPC) RegisterServices(fn func(server *grpc.Server)) {
	fn(gg.server)
}

func (gg *GRPC) Host() string {
	return gg.host
}

func (gg *GRPC) Port() string {
	return gg.port
}

func (gg *GRPC) Stop(ctx context.Context) error {
	gg.server.Stop()
	return nil
}
