package mysql

import (
	"context"
	"errors"
	"net/http"

	apperrors "github.com/ar-agahian/ice-assignment/pkg/errors"
	"github.com/ar-agahian/ice-assignment/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TodoRepository implements the TodoRepository interface using MySQL with GORM
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new MySQL TodoRepository
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// Create inserts a new todo item into the database
func (r *TodoRepository) Create(ctx context.Context, item *domain.TodoItem) error {
	result := r.db.WithContext(ctx).Create(item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID retrieves a todo item by its ID
func (r *TodoRepository) GetByID(ctx context.Context, id string) (*domain.TodoItem, error) {
	var item domain.TodoItem
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.NewAppError("INVALID_ID", "invalid todo item id", http.StatusBadRequest, nil)
	}
	result := r.db.WithContext(ctx).Where("id = ?", parsedID).First(&item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewAppError("TODO_NOT_FOUND", "todo item not found", http.StatusNotFound, nil)
		}
		return nil, result.Error
	}
	return &item, nil
}
