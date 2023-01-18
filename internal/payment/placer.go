package payment

import (
	"bytes"
	"encoding/json"
	"fmt"

	k "kafka-golang/internal/platform/kafka"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type PaymentPlacer struct {
	producer *kafka.Producer
	topic    string
}

func NewPlacer(producer *k.KafkaProducer) *PaymentPlacer {
	return &PaymentPlacer{
		producer: producer.Producer,
		topic:    producer.Topic,
	}
}

type event struct {
	Type  string
	Value *PaymentMessage
}

func (pp *PaymentPlacer) Pay(payment *PaymentMessage) error {
	return pp.publish(payment, "payments.event.pay")
}

func (pp *PaymentPlacer) publish(payment *PaymentMessage, msgType string) error {
	var payload bytes.Buffer

	event := event{
		Type:  msgType,
		Value: payment,
	}

	err := json.NewEncoder(&payload).Encode(event)
	if err != nil {
		return err
	}

	err = pp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &pp.topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload.Bytes(),
	},
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println("placed payment on the queue")

	return nil
}
