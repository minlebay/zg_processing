package grpc_server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
	"zg_processing/internal/app/kafka"
	"zg_processing/pkg/message_v1"
)

type Server struct {
	Done   chan struct{}
	Logger *zap.Logger
	Config *Config
	message.UnimplementedMessageRouterServer
	GRPCServer *grpc.Server
	Kafka      *kafka.Kafka
	wg         sync.WaitGroup
}

func NewServer(logger *zap.Logger, config *Config, kafka *kafka.Kafka) *Server {
	return &Server{
		Done:       make(chan struct{}),
		Logger:     logger,
		Config:     config,
		GRPCServer: grpc.NewServer(),
		Kafka:      kafka,
	}
}

func (s *Server) StartServer(ctx context.Context) {
	go func() {
		listener, err := net.Listen("tcp", s.Config.ListenAddress)
		if err != nil {
			s.Logger.Fatal(err.Error())
		}

		message.RegisterMessageRouterServer(s.GRPCServer, s)

		if err = s.GRPCServer.Serve(listener); err != nil {
			s.Logger.Fatal(err.Error())
		}
	}()
}

func (s *Server) StopServer(ctx context.Context) {
	s.wg.Wait()
	s.GRPCServer.Stop()
	s.Logger.Info("Server stopped")
}

func (s *Server) ReceiveMessage(ctx context.Context, m *message.Message) (*message.Response, error) {
	s.wg.Add(1)
	defer s.wg.Done()

	resp := message.Response{
		Success: true,
		Message: fmt.Sprintf("message received %v", time.Now()),
	}

	go s.Kafka.Send(context.Background(), m)

	return &resp, nil
}
