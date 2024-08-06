package kafka

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"sync"
	"time"
	message "zg_processing/pkg/message_v1"
)

type Kafka struct {
	Config *Config
	Logger *zap.Logger
	Writer *kafka.Writer
	wg     sync.WaitGroup
}

func NewKafka(logger *zap.Logger, config *Config) *Kafka {
	return &Kafka{
		Config: config,
		Logger: logger,
	}
}

func (k *Kafka) StartKafka() {
	k.Writer = &kafka.Writer{
		Addr:                   kafka.TCP(k.Config.Address),
		Topic:                  k.Config.Topics,
		AllowAutoTopicCreation: true,
		BatchTimeout:           10 * time.Millisecond,
	}
	k.Logger.Info("Kafka writer initialized", zap.String("address", k.Config.Address), zap.String("topic", k.Config.Topics))
}

func (k *Kafka) StopKafka() {
	k.wg.Wait()

	if err := k.Writer.Close(); err != nil {
		k.Logger.Error("Failed to close writer", zap.Error(err))
	} else {
		k.Logger.Info("Kafka writer closed successfully")
	}
}

func (k *Kafka) Send(ctx context.Context, message *message.Message) {
	k.wg.Add(1)
	defer k.wg.Done()

	var body bytes.Buffer
	enc := gob.NewEncoder(&body)
	err := enc.Encode(message)
	if err != nil {
		k.Logger.Error("Failed to encode message", zap.Error(err))
		return
	}

	key := []byte(message.Uuid)
	err = k.Writer.WriteMessages(ctx,
		kafka.Message{
			Value: body.Bytes(),
			Key:   key,
			Headers: []kafka.Header{
				{Key: "MessageKey", Value: key},
			},
		})
	if err != nil {
		k.Logger.Error("Failed to write message", zap.Error(err))
	} else {
		k.Logger.Info("Message written successfully", zap.ByteString("key", key))
	}
}
