package app

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
	"testing"
	"zg_processing/internal/app/grpc_server"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(
		fx.Options(
			grpc_server.NewModule(),
			//kafka.NewModule(),
		),
		fx.Provide(
			zap.NewProduction,
			NewConfig,
		),
	)
	require.NoError(t, err)
}

func TestKafka(t *testing.T) {

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "processing_1",
	})

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello Kafka!"),
		},
	)

	if err != nil {
		log.Fatalf("could not write message %v", err)
	} else {
		log.Println("message written successfully")
	}

	writer.Close()

}
