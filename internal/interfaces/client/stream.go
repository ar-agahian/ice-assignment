package client

import (
	"context"
)

// IStreamPublisher defines the interface for publishing messages to streams
type IStreamPublisher interface {
	Publish(ctx context.Context, stream string, data map[string]interface{}) error
}

