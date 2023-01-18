package awslibs

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	bucketName string
	client     *s3.Client
	uploader   *manager.Uploader
}

func NewS3Client(cfg aws.Config, bucketName string) *S3Client {
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// localstack config
		o.UsePathStyle = true
	})

	uploader := manager.NewUploader(client)
	return &S3Client{
		client:     client,
		uploader:   uploader,
		bucketName: bucketName,
	}
}

func (c *S3Client) Upload(filename string, body []byte) error {
	fmt.Printf("uploading %s to s3\n", filename)
	_, err := c.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(body),
	})
	return err
}

func (c *S3Client) ListBuckets() {
	res, err := c.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	fmt.Println("listing buckets")
	if err != nil {
		fmt.Println(err)
	}
	for _, bucket := range res.Buckets {
		fmt.Printf("the bucket name is : %s\n", *bucket.Name)
	}
}
