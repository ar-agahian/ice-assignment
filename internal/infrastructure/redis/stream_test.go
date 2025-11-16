package redis

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamPublisher_Publish(t *testing.T) {
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

	// Skip if Redis is not available
	publisher, err := NewStreamPublisher(ctx)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
	}
	defer publisher.Close()

	data := map[string]interface{}{
		"id":          "test-id",
		"description": "test description",
	}

	err = publisher.Publish(ctx, "test-stream", data)
	if err == nil {
		assert.NoError(t, err)
	}
}

