package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ar-agahian/ice-assignment/mocks"
	apperrors "github.com/ar-agahian/ice-assignment/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTodoItem(t *testing.T) {
	tests := []struct {
		name          string
		req           CreateTodoItemRequest
		setupMocks    func(*mocks.MockITodoRepository, *mocks.MockIStreamPublisher)
		expectedError error
	}{
		{
			name: "successful creation",
			req: CreateTodoItemRequest{
				Description: "Test todo",
				DueDate:     time.Now().Add(24 * time.Hour),
				FileID:      "file-123",
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				todoRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.TodoItem")).Return(nil)
				streamRepo.On("Publish", mock.Anything, "todo-items", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "empty description",
			req: CreateTodoItemRequest{
				Description: "",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				// No mocks needed, validation fails early
			},
			expectedError: apperrors.NewAppError("INVALID_DESCRIPTION", "description cannot be empty", http.StatusBadRequest, nil),
		},
		{
			name: "description too long",
			req: CreateTodoItemRequest{
				Description: strings.Repeat("a", 501),
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				// No mocks needed, validation fails early
			},
			expectedError: apperrors.NewAppError("INVALID_DESCRIPTION", "description must be at most 500 characters", http.StatusBadRequest, nil),
		},
		{
			name: "past due date",
			req: CreateTodoItemRequest{
				Description: "Test todo",
				DueDate:     time.Now().Add(-24 * time.Hour),
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				// No mocks needed, validation fails early
			},
			expectedError: apperrors.NewAppError("INVALID_DUE_DATE", "due date must be in the future", http.StatusBadRequest, nil),
		},
		{
			name: "database error",
			req: CreateTodoItemRequest{
				Description: "Test todo",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				todoRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.TodoItem")).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
		{
			name: "stream publish error",
			req: CreateTodoItemRequest{
				Description: "Test todo",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				todoRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.TodoItem")).Return(nil)
				streamRepo.On("Publish", mock.Anything, "todo-items", mock.Anything).Return(errors.New("stream error"))
			},
			expectedError: errors.New("stream error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoRepo := mocks.NewMockITodoRepository(t)
			streamRepo := mocks.NewMockIStreamPublisher(t)
			tt.setupMocks(todoRepo, streamRepo)

			uc := NewTodoUseCase(todoRepo, streamRepo)
			result, err := uc.CreateTodoItem(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if appErr, ok := apperrors.AsAppError(tt.expectedError); ok {
					actualErr, ok := apperrors.AsAppError(err)
					assert.True(t, ok, "expected AppError")
					assert.Equal(t, appErr.Code, actualErr.Code)
				} else {
					assert.NotNil(t, err)
				}
				if tt.name == "stream publish error" {
					// Stream errors don't fail the operation, item is still returned
					assert.NotNil(t, result)
				} else {
					assert.Nil(t, result)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Description, result.Description)
				todoRepo.AssertExpectations(t)
				streamRepo.AssertExpectations(t)
			}
		})
	}
}
