package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Reader is a contract to a kafka writer.
type Reader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Processor is a contract to a kafka message processor.
// type Processor interface {
// 	Process(ctx context.Context, msg []byte) error
// }

type Processor func(ctx context.Context, msg []byte) error

// Consumer implements functionality to consume kafka messages.
type Consumer struct {
	reader           Reader
	processor        Processor
	consumeCtxCancel context.CancelFunc
}

// NewConsumer is a constructor function for Consumer.
func NewConsumer(reader Reader, processor Processor) *Consumer {
	return &Consumer{
		reader:    reader,
		processor: processor,
	}
}

// ConsumeMessage consumes and processes a single kafka message.
func (p *Consumer) ConsumeMessage(ctx context.Context) error {
	m, err := p.reader.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("fetching message: %w", err)
	}

	if err := p.processor(ctx, m.Value); err != nil {
		return fmt.Errorf("processing message: %w", err)
	}

	if err := p.reader.CommitMessages(ctx, m); err != nil {
		return fmt.Errorf("comitting message: %w", err)
	}

	return nil
}

// ConsumeMessages consumes messages continuously and sends errors to a channel. Stops consuming on ctx.Done().
func (p *Consumer) ConsumeMessages(ctx context.Context, errs chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	p.consumeCtxCancel = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(errs)
				return
			default:
			}

			err := p.ConsumeMessage(ctx)
			if err != nil {
				errs <- err
			}
		}
	}()
}

// Stop stops the consumer gracefully.
func (p *Consumer) Stop() error {
	p.consumeCtxCancel()

	if err := p.reader.Close(); err != nil {
		return fmt.Errorf("closing kafka reader: %w", err)
	}

	return nil
}
