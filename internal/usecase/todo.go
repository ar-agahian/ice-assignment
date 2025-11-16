package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/ar-agahian/ice-assignment/internal/domain"
	"github.com/ar-agahian/ice-assignment/internal/interfaces/client"
	"github.com/ar-agahian/ice-assignment/internal/interfaces/repository"
	apperrors "github.com/ar-agahian/ice-assignment/pkg/errors"
)

// TodoUseCase handles todo item business logic
type TodoUseCase struct {
	todoRepo   repository.ITodoRepository
	streamRepo client.IStreamPublisher
}

// NewTodoUseCase creates a new TodoUseCase
func NewTodoUseCase(todoRepo repository.ITodoRepository, streamRepo client.IStreamPublisher) *TodoUseCase {
	return &TodoUseCase{
		todoRepo:   todoRepo,
		streamRepo: streamRepo,
	}
}

// CreateTodoItemRequest represents the request to create a todo item
type CreateTodoItemRequest struct {
	Description string
	DueDate     time.Time
	FileID      string
}

// CreateTodoItem creates a new todo item and publishes it to the stream
func (uc *TodoUseCase) CreateTodoItem(ctx context.Context, req CreateTodoItemRequest) (*domain.TodoItem, error) {
	if req.Description == "" {
		return nil, apperrors.NewAppError("INVALID_DESCRIPTION", "description cannot be empty", http.StatusBadRequest, nil)
	}
	if len(req.Description) > 500 {
		return nil, apperrors.NewAppError("INVALID_DESCRIPTION", "description must be at most 500 characters", http.StatusBadRequest, nil)
	}
	if req.DueDate.Before(time.Now()) {
		return nil, apperrors.NewAppError("INVALID_DUE_DATE", "due date must be in the future", http.StatusBadRequest, nil)
	}
	todoItem := domain.NewTodoItem(req.Description, req.DueDate, req.FileID)
	if err := uc.todoRepo.Create(ctx, todoItem); err != nil {
		return nil, err
	}
	streamData := map[string]interface{}{
		"id":          todoItem.ID.String(),
		"description": todoItem.Description,
		"dueDate":     todoItem.DueDate.Format(time.RFC3339),
		"fileId":      todoItem.FileID,
	}
	if err := uc.streamRepo.Publish(ctx, "todo-items", streamData); err != nil {
		return todoItem, err
	}
	return todoItem, nil
}
