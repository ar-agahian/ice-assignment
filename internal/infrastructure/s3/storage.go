package s3

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// FileStorage implements the FileStorage interface using AWS S3
type FileStorage struct {
	client     *s3.Client
	bucketName string
}

// NewFileStorage creates a new S3 FileStorage
func NewFileStorage(ctx context.Context) (*FileStorage, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	endpoint := os.Getenv("S3_ENDPOINT")
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		cfg.BaseEndpoint = aws.String(endpoint)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.UsePathStyle = true
		}
	})
	storage := &FileStorage{
		client:     client,
		bucketName: bucketName,
	}
	if err := storage.ensureBucket(ctx); err != nil {
		return nil, err
	}
	return storage, nil
}

// Upload uploads a file to S3 and returns the file ID
func (s *FileStorage) Upload(ctx context.Context, file io.Reader, contentType string) (string, error) {
	fileID := uuid.New().String()
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fileID),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return fileID, nil
}

// Get retrieves a file from S3 by file ID
func (s *FileStorage) Get(ctx context.Context, fileID string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ensureBucket creates the bucket if it doesn't exist
func (s *FileStorage) ensureBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	})
	if err == nil {
		// Bucket exists
		return nil
	}
	// Try to create bucket
	_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucketName),
	})
	return err
}
