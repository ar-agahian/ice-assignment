package redis

import (
	"context"
	"os"
	"testing"
)

func BenchmarkStreamPublisher_Publish(b *testing.B) {
	ctx := context.Background()

	// Set environment variables
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")
	defer func() {
		os.Unsetenv("REDIS_ADDR")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("REDIS_DB")
	}()

	publisher, err := NewStreamPublisher(ctx)
	if err != nil {
		b.Skipf("Skipping benchmark: Redis not available: %v", err)
	}
	defer publisher.Close()

	data := map[string]interface{}{
		"id":          "benchmark-id",
		"description": "benchmark description",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = publisher.Publish(ctx, "benchmark-stream", data)
	}
}

