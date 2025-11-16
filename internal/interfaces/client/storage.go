package client

import (
	"context"
	"io"
)

// IFileStorage defines the interface for file storage operations
type IFileStorage interface {
	Upload(ctx context.Context, file io.Reader, contentType string) (fileID string, err error)
	Get(ctx context.Context, fileID string) ([]byte, error)
}
