package app

import (
	"go.uber.org/fx"
	"zg_processing/internal/app/grpc_server"
	"zg_processing/internal/app/kafka"
	"zg_processing/internal/app/log"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			grpc_server.NewModule(),
			kafka.NewModule(),
			log.NewModule(),
		),
		fx.Provide(
			NewConfig,
		),
	)
}
