package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kafka-golang/internal/platform/awslibs"
	k "kafka-golang/internal/platform/kafka"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type PaymentConsumer struct {
	consumer *kafka.Consumer
	topic    string
	s3Client *awslibs.S3Client
}

func NewConsumer(consumer *k.KafkaConsumer, s3Client *awslibs.S3Client) *PaymentConsumer {
	return &PaymentConsumer{
		consumer: consumer.Consumer,
		topic:    consumer.Topic,
		s3Client: s3Client,
	}
}

func (pc *PaymentConsumer) Consume() error {
	err := pc.consumer.Subscribe(pc.topic, nil)
	if err != nil {
		return err
	}

	for {
		msg, ok := pc.consumer.Poll(100).(*kafka.Message)

		if !ok {
			continue
		}

		var evt struct {
			Type  string
			Value PaymentMessage
		}

		if err := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&evt); err != nil {
			fmt.Println("ignoring message, invalid", err)

			continue
		}

		ok = false

		switch evt.Type {
		case "payments.event.pay":

			fmt.Printf("processing payment: %s\n", evt.Value.ClientReferenceInformation.Code)

			var payload bytes.Buffer

			enc := json.NewEncoder(&payload)
			enc.SetIndent("", "    ")
			err := enc.Encode(evt.Value)
			if err != nil {
				fmt.Println("failed to encode payload", err)

				continue
			}

			err = pc.s3Client.Upload(fmt.Sprintf("%s.json", evt.Value.ClientReferenceInformation.Code), payload.Bytes())

			if err != nil {
				fmt.Printf("failed to upload file: %v", err)
			}

			fmt.Println("finish processing payment")
		}
	}
}
