package grpc_server

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {

	return fx.Module(
		"server",
		fx.Provide(
			NewServerConfig,
			NewServer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, r *Server) {
				lc.Append(fx.Hook{
					OnStart: func(context.Context) error {
						go r.StartServer()
						return nil
					},
					OnStop: func(context.Context) error {
						r.StopServer()
						return nil
					},
				})
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("server")
		}),
	)
}
