package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	Topic    string
}

func NewKafkaConsumer(host string, groupID string, topic string) (*KafkaConsumer, error) {

	config := kafka.ConfigMap{
		"bootstrap.servers": host,
		"group.id":          groupID,
		"auto.offset.reset": "smallest",
	}

	client, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}

	if err := client.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		Consumer: client,
		Topic:    topic,
	}, nil
}
