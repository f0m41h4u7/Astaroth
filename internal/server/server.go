package server

import (
	"net"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//go:generate protoc --proto_path=../../pkg/api/ --go_out=plugins=grpc:../../pkg/api ../../pkg/api/stats.proto
type Server struct {
	grpc *grpc.Server
}

func InitServer() *Server {
	//	app = cl
	s := &Server{}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(zap.L())))
	api.RegisterCalendarServer(grpcServer, s)
	s.grpc = grpcServer
	return s
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", "127.0.0.1:1337")
	if err != nil {
		return err
	}
	err = s.grpc.Serve(lis)
	return err
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
