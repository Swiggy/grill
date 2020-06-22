package grillgrpc

import (
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

func Start() (*GRPC, error) {
	recorder := newGRPCRecorder()
	server := grpc.NewServer(grpc.UnaryInterceptor(recorder.unaryInterceptor))

	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	go server.Serve(listen)

	port := listen.Addr().(*net.TCPAddr).Port
	return &GRPC{
		server:   server,
		host:     "localhost",
		port:     fmt.Sprintf("%d", port),
		recorder: recorder,
	}, nil
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

func (gg *GRPC) Stop() error {
	gg.server.Stop()
	return nil
}
