package grillgrpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GrillGRPC struct {
	server   *grpc.Server
	recorder *recorder
	host     string
	port     string
}

func Start() (*GrillGRPC, error) {
	recorder := newGRPCRecorder()
	server := grpc.NewServer(grpc.UnaryInterceptor(recorder.unaryInterceptor))

	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	go server.Serve(listen)

	port := listen.Addr().(*net.TCPAddr).Port
	return &GrillGRPC{
		server:   server,
		host:     "localhost",
		port:     fmt.Sprintf("%d", port),
		recorder: recorder,
	}, nil
}

func (gg *GrillGRPC) RegisterServices(fn func(server *grpc.Server)) {
	fn(gg.server)
}

func (gg *GrillGRPC) Host() string {
	return gg.host
}

func (gg *GrillGRPC) Port() string {
	return gg.port
}

func (gg *GrillGRPC) Stop() error {
	gg.server.Stop()
	return nil
}
