package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Writer is a contract to a kafka writer.
type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

// Producer implements functionality to produce kafka messages.
type Producer struct {
	writer Writer
}

// NewProducer is a constructor function for Producer.
func NewProducer(writer Writer) *Producer {
	return &Producer{
		writer: writer,
	}
}

// Produce produces a single kafka message.
func (p *Producer) Produce(ctx context.Context, channel string, message any) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: channel,
		Value: bytes,
	}); err != nil {
		return fmt.Errorf("could not write message: %w", err)
	}

	return nil
}
