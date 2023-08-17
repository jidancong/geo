package repo

import (
	"context"

	"github.com/jidancong/geo/pkg/kafka"
)

type ConsumerKafka struct {
	kafkaClient *kafka.KafkaPkg
}

func NewConsumerKafka(kafkaClient *kafka.KafkaPkg) *ConsumerKafka {
	return &ConsumerKafka{kafkaClient}
}

func (c *ConsumerKafka) Consumer(processor kafka.Processor) {
	consumer := c.kafkaClient.NewConsumer(processor)
	channel := make(chan error, 1)
	consumer.ConsumeMessages(context.TODO(), channel)
}
