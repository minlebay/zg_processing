package app

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
	"zg_processing/internal/app/grpc_server"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(
		fx.Options(
			grpc_server.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig,
		),
	)
	require.NoError(t, err)
}
