package grpc_server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"time"
	"zg_processing/pkg/message_v1/router"
)

type Server struct {
	Done   chan struct{}
	Logger *zap.Logger
	Config *Config
	router.UnimplementedMessageRouterServer
	GRPCServer *grpc.Server
}

func NewServer(logger *zap.Logger, config *Config) *Server {
	return &Server{
		Done:       make(chan struct{}),
		Logger:     logger,
		Config:     config,
		GRPCServer: grpc.NewServer(),
	}
}

func (s *Server) StartServer(ctx context.Context) {
	go func() {
		listener, err := net.Listen("tcp", s.Config.ListenAddress)
		if err != nil {
			s.Logger.Fatal(err.Error())
		}

		router.RegisterMessageRouterServer(s.GRPCServer, s)

		if err = s.GRPCServer.Serve(listener); err != nil {
			s.Logger.Fatal(err.Error())
		}

		s.Logger.Info("Server started at address " + s.Config.ListenAddress)
	}()
}

func (s *Server) StopServer(ctx context.Context) {
	s.Logger.Info("Server stopped")
	s.Done <- struct{}{}
}

func (s *Server) ReceiveMessage(ctx context.Context, m *router.Message) (*router.Response, error) {

	s.Logger.Info("message received: ", zap.Any("message", m))

	resp := router.Response{
		Success: true,
		Message: fmt.Sprintf("message received %v", time.Now()),
	}

	return &resp, nil
}
