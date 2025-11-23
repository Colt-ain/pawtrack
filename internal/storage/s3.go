package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Storage(region, bucket string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Storage{
		client: client,
		bucket: bucket,
		region: region,
	}, nil
}

func (s *S3Storage) Upload(file io.Reader, filename string, contentType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("events/%s%s", uuid.New().String(), ext)

	// Upload to S3
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(uniqueFilename),
		Body:        file.(io.ReadSeeker),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Return the public URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, uniqueFilename)
	return fileURL, nil
}

func (s *S3Storage) Delete(fileURL string) error {
	// Extract key from URL and delete from S3
	// TODO: Implement deletion
	return nil
}

func (s *S3Storage) GetSignedURL(fileURL string, expiryDuration time.Duration) (string, error) {
	// TODO: Implement signed URL generation using presigner
	return fileURL, nil
}
