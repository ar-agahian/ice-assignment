package repository

import (
	"context"

	"github.com/ar-agahian/ice-assignment/internal/domain"
)

// ITodoRepository defines the interface for todo item persistence
type ITodoRepository interface {
	Create(ctx context.Context, item *domain.TodoItem) error
	GetByID(ctx context.Context, id string) (*domain.TodoItem, error)
}
