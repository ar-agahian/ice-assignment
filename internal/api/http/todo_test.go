package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ar-agahian/ice-assignment/internal/usecase"
	"github.com/ar-agahian/ice-assignment/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTodoHandler_CreateTodo(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mocks.MockITodoRepository, *mocks.MockIStreamPublisher)
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: CreateTodoRequest{
				Description: "Test todo",
				DueDate:     time.Now().Add(24 * time.Hour),
				FileID:      "file-123",
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				todoRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				streamRepo.On("Publish", mock.Anything, "todo-items", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				// No mocks needed
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty description",
			requestBody: CreateTodoRequest{
				Description: "",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			setupMocks: func(todoRepo *mocks.MockITodoRepository, streamRepo *mocks.MockIStreamPublisher) {
				// No mocks needed, validation fails early
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			todoRepo := mocks.NewMockITodoRepository(t)
			streamRepo := mocks.NewMockIStreamPublisher(t)
			tt.setupMocks(todoRepo, streamRepo)

			todoUseCase := usecase.NewTodoUseCase(todoRepo, streamRepo)
			handler := NewTodoHandler(todoUseCase)

			router := gin.New()
			router.POST("/todo", handler.CreateTodo)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/todo", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
