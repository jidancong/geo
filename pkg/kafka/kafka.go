package kafka

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
)

type KafkaPkg struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaPkg(host, topic, groupId string, numPart, replica int) (*KafkaPkg, error) {
	hosts := strings.Split(host, ",")
	if len(hosts) <= 0 {
		return nil, fmt.Errorf("kafka broker: %v", hosts)
	}

	conn, err := kafka.Dial("tcp", hosts[0])
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, fmt.Errorf("kafka controller error: %s", err)
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPart,
			ReplicationFactor: replica,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)

	if err != nil {
		return nil, fmt.Errorf("create topic error: %s", err)
	}

	// kafka writer
	w := kafka.NewWriter(kafka.WriterConfig{
		Topic:    topic,
		Brokers:  hosts,
		Balancer: &kafka.RoundRobin{},
		Async:    true,
	})

	// kafka reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: hosts,
		Topic:   topic,
		GroupID: groupId,
		// MaxBytes: 10e6, //10MB
	})

	return &KafkaPkg{w, r}, nil
}

func (k *KafkaPkg) NewConsumer(processor Processor) *Consumer {
	return NewConsumer(k.reader, processor)
}

func (k *KafkaPkg) NewProducer() *Producer {
	return NewProducer(k.writer)
}
