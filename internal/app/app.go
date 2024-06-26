package app

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"zg_processing/internal/app/grpc_server"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			grpc_server.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig,
		),
	)
}
