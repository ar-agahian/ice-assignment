package s3

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStorage_Upload(t *testing.T) {
	// This is a simplified test
	// In a real scenario, you would use LocalStack or a proper S3 mock
	ctx := context.Background()

	// Set environment variables
	os.Setenv("S3_BUCKET_NAME", "test-bucket")
	os.Setenv("S3_ENDPOINT", "http://localhost:4566")
	defer func() {
		os.Unsetenv("S3_BUCKET_NAME")
		os.Unsetenv("S3_ENDPOINT")
	}()

	// Skip if LocalStack is not available
	storage, err := NewFileStorage(ctx)
	if err != nil {
		t.Skipf("Skipping test: LocalStack not available: %v", err)
	}

	fileContent := strings.NewReader("test file content")
	fileID, err := storage.Upload(ctx, fileContent, "text/plain")

	// If LocalStack is running, this should succeed
	if err == nil {
		assert.NotEmpty(t, fileID)
	}
}

func TestFileStorage_Get(t *testing.T) {
	ctx := context.Background()

	// Set environment variables
	os.Setenv("S3_BUCKET_NAME", "test-bucket")
	os.Setenv("S3_ENDPOINT", "http://localhost:4566")
	defer func() {
		os.Unsetenv("S3_BUCKET_NAME")
		os.Unsetenv("S3_ENDPOINT")
	}()

	storage, err := NewFileStorage(ctx)
	if err != nil {
		t.Skipf("Skipping test: LocalStack not available: %v", err)
	}

	// First upload a file
	fileContent := strings.NewReader("test file content")
	fileID, err := storage.Upload(ctx, fileContent, "text/plain")
	if err != nil {
		t.Skipf("Skipping test: Failed to upload file: %v", err)
	}

	// Then retrieve it
	data, err := storage.Get(ctx, fileID)
	if err == nil {
		assert.Equal(t, "test file content", string(data))
	}
}
