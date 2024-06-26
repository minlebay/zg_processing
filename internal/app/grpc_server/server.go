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
}

func NewServer(logger *zap.Logger, config *Config) *Server {
	return &Server{
		Done:   make(chan struct{}),
		Logger: logger,
		Config: config,
	}
}

func (r *Server) StartServer() {
	listener, err := net.Listen("tcp", r.Config.ListenAddress)
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	server := grpc.NewServer()

	router.RegisterMessageRouterServer(server, r)

	if err = server.Serve(listener); err != nil {
		r.Logger.Fatal(err.Error())
	}

	r.Logger.Info("Server started at address " + r.Config.ListenAddress)
}

func (r *Server) StopServer() {
	r.Logger.Info("Server stopped")
	r.Done <- struct{}{}
}

func (r *Server) ReceiveMessage(ctx context.Context, m *router.Message) (*router.Response, error) {

	r.Logger.Info("message received: ", zap.Any("message", m))

	resp := router.Response{
		Success: true,
		Message: fmt.Sprintf("message received %v", time.Now()),
	}

	return &resp, nil
}
