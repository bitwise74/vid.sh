// Package aws defines functions used to interact with the AWS API
package aws

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"go.uber.org/zap"
)

type S3Client struct {
	C      *s3.Client
	Bucket *string
}

func NewS3() (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("ACCESS_KEY_ID"),
			os.Getenv("SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	bucket := aws.String(os.Getenv("BUCKET"))

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = os.Getenv("REGION")
	})

	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		var apiErr smithy.APIError

		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NotFound" {
				return nil, fmt.Errorf("bucket '%s' does not exist", *bucket)
			}
		}

		return nil, fmt.Errorf("failed to check if bucket exists, %w", err)
	}

	// Create any missing folders
	_, err = client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    aws.String("avatars/"),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			// If the object does not exist, create it
			if apiErr.ErrorCode() == "NotFound" {
				_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
					Bucket: bucket,
					Key:    aws.String("avatars/"),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to create 'avatars' folder, %w", err)
				}

				zap.L().Info("Created avatars/ directory in S3 bucket")
			} else {
				return nil, fmt.Errorf("failed to check if 'avatars/' exists, %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to check if 'avatars/' exists, %w", err)
		}
	}

	zap.L().Info("S3 client initialized", zap.String("bucket", *bucket))

	return &S3Client{
		C:      client,
		Bucket: bucket,
	}, nil
}
