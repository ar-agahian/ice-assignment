package app

import (
	"context"

	httphandler "github.com/ar-agahian/ice-assignment/internal/api/http"
	"github.com/ar-agahian/ice-assignment/internal/infrastructure/mysql"
	"github.com/ar-agahian/ice-assignment/internal/infrastructure/redis"
	"github.com/ar-agahian/ice-assignment/internal/infrastructure/s3"
	"github.com/ar-agahian/ice-assignment/internal/usecase"
	"gorm.io/gorm"
)

// App holds all application dependencies
type App struct {
	DB              *gorm.DB
	TodoUseCase     *usecase.TodoUseCase
	FileUseCase     *usecase.FileUseCase
	Handler         *httphandler.Handler
	StreamPublisher *redis.StreamPublisher
}

// NewApp initializes all application dependencies
func NewApp() (*App, error) {
	// database
	db, err := mysql.NewDatabase()
	if err != nil {
		return nil, err
	}
	if err := mysql.RunMigrations(db); err != nil {
		return nil, err
	}

	// repositories
	todoRepo := mysql.NewTodoRepository(db)

	// infrastructure clients
	ctx := context.Background()
	s3Storage, err := s3.NewFileStorage(ctx)
	if err != nil {
		return nil, err
	}

	streamPublisher, err := redis.NewStreamPublisher(ctx)
	if err != nil {
		return nil, err
	}

	// usecases
	todoUseCase := usecase.NewTodoUseCase(todoRepo, streamPublisher)
	fileUseCase := usecase.NewFileUseCase(s3Storage)

	// http-handler
	handler := httphandler.NewHandler(todoUseCase, fileUseCase)

	return &App{
		DB:              db,
		TodoUseCase:     todoUseCase,
		FileUseCase:     fileUseCase,
		Handler:         handler,
		StreamPublisher: streamPublisher,
	}, nil
}

// Close closes all application resources
func (a *App) Close() error {
	return nil
}
