package app

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"zg_processing/internal/app/grpc_server"
	"zg_processing/internal/app/kafka"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			grpc_server.NewModule(),
			kafka.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig,
		),
	)
}
