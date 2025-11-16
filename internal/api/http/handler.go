package http

import (
	"github.com/ar-agahian/ice-assignment/internal/usecase"
	apperrors "github.com/ar-agahian/ice-assignment/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Handler sets up HTTP routes and middleware
type Handler struct {
	todoHandler *TodoHandler
	fileHandler *FileHandler
}

// NewHandler creates a new HTTP handler
func NewHandler(todoUseCase *usecase.TodoUseCase, fileUseCase *usecase.FileUseCase) *Handler {
	return &Handler{
		todoHandler: NewTodoHandler(todoUseCase),
		fileHandler: NewFileHandler(fileUseCase),
	}
}

// SetupRoutes configures all HTTP routes
func (h *Handler) SetupRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(errorHandler())
	api := r.Group("/api")
	{
		h.todoHandler.RegisterRoutes(api)
		h.fileHandler.RegisterRoutes(api)
	}
	return r
}

// errorHandler is a middleware that handles errors consistently
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			statusCode := apperrors.GetHTTPStatus(err)
			response := apperrors.GetErrorResponse(err)
			c.JSON(statusCode, response)
			return
		}
	}
}
