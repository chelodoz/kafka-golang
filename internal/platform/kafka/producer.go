package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

func NewKafkaProducer(host string, clientID string, topic string) (*KafkaProducer, error) {

	config := kafka.ConfigMap{
		"bootstrap.servers": host,
		"client.id":         clientID,
		"acks":              "all",
	}

	client, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: client,
		Topic:    topic,
	}, nil
}
