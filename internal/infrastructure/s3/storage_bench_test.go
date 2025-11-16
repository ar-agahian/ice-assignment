package s3

import (
	"context"
	"os"
	"strings"
	"testing"
)

func BenchmarkFileStorage_Upload(b *testing.B) {
	ctx := context.Background()
	os.Setenv("S3_BUCKET_NAME", "test-bucket")
	os.Setenv("S3_ENDPOINT", "http://localhost:4566")
	defer func() {
		os.Unsetenv("S3_BUCKET_NAME")
		os.Unsetenv("S3_ENDPOINT")
	}()

	storage, err := NewFileStorage(ctx)
	if err != nil {
		b.Skipf("Skipping benchmark: LocalStack not available: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = storage.Upload(ctx, strings.NewReader("benchmark test content"), "text/plain")
	}
}

