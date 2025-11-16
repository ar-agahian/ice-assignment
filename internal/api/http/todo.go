package http

import (
	"net/http"
	"time"

	"github.com/ar-agahian/ice-assignment/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TodoHandler handles todo-related HTTP requests
type TodoHandler struct {
	todoUseCase *usecase.TodoUseCase
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(todoUseCase *usecase.TodoUseCase) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
	}
}

// CreateTodoRequest represents the request body for creating a todo item
type CreateTodoRequest struct {
	Description string    `json:"description" binding:"required"`
	DueDate     time.Time `json:"dueDate" binding:"required"`
	FileID      string    `json:"fileId,omitempty" binding:"omitempty,uuid"`
}

// TodoResponse represents the response for a todo item
type TodoResponse struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId,omitempty"`
}

// CreateTodo handles POST /todo requests
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	todoItem, err := h.todoUseCase.CreateTodoItem(c.Request.Context(), usecase.CreateTodoItemRequest{
		Description: req.Description,
		DueDate:     req.DueDate,
		FileID:      req.FileID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, TodoResponse{
		ID:          todoItem.ID.String(),
		Description: todoItem.Description,
		DueDate:     todoItem.DueDate,
		FileID:      todoItem.FileID,
	})
}

// RegisterRoutes registers todo routes
func (h *TodoHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/todo", h.CreateTodo)
}

