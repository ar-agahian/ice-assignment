package redis

import (
	"context"
	"encoding/json"
	"os"

	"github.com/redis/go-redis/v9"
)

// StreamPublisher implements the StreamPublisher interface using Redis Streams
type StreamPublisher struct {
	client *redis.Client
}

// NewStreamPublisher creates a new Redis StreamPublisher
func NewStreamPublisher(ctx context.Context) (*StreamPublisher, error) {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &StreamPublisher{client: rdb}, nil
}

// Publish publishes a message to a Redis stream
func (p *StreamPublisher) Publish(ctx context.Context, stream string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	args := redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"data": string(jsonData),
		},
	}
	if err := p.client.XAdd(ctx, &args).Err(); err != nil {
		return err
	}
	return nil
}

// Close closes the Redis connection
func (p *StreamPublisher) Close() error {
	return p.client.Close()
}
