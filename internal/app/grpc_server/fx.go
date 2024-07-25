package grpc_server

import (
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
				lc.Append(fx.StartStopHook(r.StartServer, r.StopServer))
			},
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("processing_server") // TODO add identifier of which server is being decorated
		}),
	)
}
