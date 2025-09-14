package pb

import (
	"net"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/server"
	"github.com/kirillmashkov/shortener.git/internal/service"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

type GRPCServer struct {
	UnimplementedShortenerServer
	server    *grpc.Server
	service service.Service
}

func (s *GRPCServer) Run() error {
	RegisterShortenerServer(s.server, s)

	listen, err := net.Listen("tcp", "localhost:3200")
	if err != nil {
		app.Log.Fatal("can't start grpc server", zap.Error(err))
		return err
	}

	return s.server.Serve(listen)
}

func (s *GRPCServer) Shutdown() error {
	s.server.GracefulStop()
	return nil
}

func New(service service.Service) server.Server {
	s := grpc.NewServer()

	return &GRPCServer{server: s, service: service}
}