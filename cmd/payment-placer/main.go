package main

import (
	"fmt"
	"kafka-golang/internal/payment"
	"kafka-golang/internal/platform/kafka"
	"time"
)

func main() {
	topic := "test-topic"
	p, err := kafka.NewKafkaProducer("localhost:9092", "foo", topic)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
	}

	pp := payment.NewPlacer(p)

	for i := 0; i < 1000; i++ {
		fakePayment := payment.NewMessage()

		if err := pp.Pay(fakePayment); err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Second * 3)
	}
}
