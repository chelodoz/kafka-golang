package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	AWS_ACCESS_KEY  = "test"
	AWS_SECRET_KEY  = "test"
	AWS_S3_REGION   = "us-east-2"             // Region
	AWS_S3_ENDPOINT = "http://localhost:4566" // Endpoint
	AWS_S3_BUCKET   = "my-bucket"             // Bucket
)

func main() {
	// s3
	creds := credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY, AWS_SECRET_KEY, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(AWS_S3_REGION),

		// localstack config
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: AWS_S3_ENDPOINT}, nil
			}),
		))

	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// localstack config
		o.UsePathStyle = true
	})
	uploader := manager.NewUploader(client)

	//kafka
	topic := "test-topic"
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "foo",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:

			fmt.Printf("processing order: %s\n", string(e.Value))

			err := UploadToS3(uploader, AWS_S3_BUCKET, fmt.Sprintf("%s.txt", string(e.Value)), e.Value)
			if err != nil {
				fmt.Printf("Failed to upload file: %v", err)
			}

			fmt.Println("finish processing order")
		case *kafka.Error:
			fmt.Printf("%s\n", e)
		}
	}
}

func UploadToS3(uploader *manager.Uploader, bucketName string, filename string, body []byte) error {

	fmt.Printf("uploading %s to s3\n", filename)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(body),
	})
	return err
}

// func ListBuckets(client *s3.Client) {
// 	res, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
// 	fmt.Println("listing buckets")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	for _, bucket := range res.Buckets {
// 		fmt.Printf("the bucket name is : %s\n", *bucket.Name)
// 	}
// }
