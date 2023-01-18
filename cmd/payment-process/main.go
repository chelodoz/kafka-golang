package main

import (
	"context"
	"kafka-golang/internal/payment"
	"kafka-golang/internal/platform/awslibs"
	"log"

	k "kafka-golang/internal/platform/kafka"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const (
	AWS_ACCESS_KEY = "test"
	AWS_SECRET_KEY = "test"
	AWS_REGION     = "us-east-2"             // Region
	AWS_ENDPOINT   = "http://localhost:4566" // Endpoint
	AWS_S3_BUCKET  = "my-bucket"             // Bucket
)

func main() {
	creds := credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY, AWS_SECRET_KEY, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(AWS_REGION),
		// localstack config
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: AWS_ENDPOINT}, nil
			}),
		))

	if err != nil {
		log.Fatal(err)
	}

	s3 := awslibs.NewS3Client(cfg, AWS_S3_BUCKET)

	topic := "test-topic"
	consumer, err := k.NewKafkaConsumer("localhost:9092", "foo", topic)
	if err != nil {
		log.Fatal(err)
	}

	paymentConsumer := payment.NewConsumer(consumer, s3)

	err = paymentConsumer.Consume()
	if err != nil {
		log.Fatal(err)
	}
}
