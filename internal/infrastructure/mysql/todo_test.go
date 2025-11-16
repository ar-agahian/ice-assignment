package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/ar-agahian/ice-assignment/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database connection
	dsn := "root:password@tcp(localhost:3306)/todo_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: failed to connect to test database: %v", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		t.Skipf("Skipping test: failed to get underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Skipf("Skipping test: failed to ping test database: %v", err)
	}

	// AutoMigrate to create tables
	err = db.AutoMigrate(&domain.TodoItem{})
	require.NoError(t, err)

	// Clean up
	t.Cleanup(func() {
		db.Exec("DROP TABLE IF EXISTS todo_items")
		sqlDB.Close()
	})

	return db
}

func TestTodoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTodoRepository(db)

	todo := domain.NewTodoItem("Test description", time.Now().Add(24*time.Hour), "file-123")

	err := repo.Create(context.Background(), todo)
	assert.NoError(t, err)

	// Verify it was created
	retrieved, err := repo.GetByID(context.Background(), todo.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, todo.Description, retrieved.Description)
	assert.Equal(t, todo.FileID, retrieved.FileID)
}

func TestTodoRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTodoRepository(db)

	todo := domain.NewTodoItem("Test description", time.Now().Add(24*time.Hour), "file-123")
	err := repo.Create(context.Background(), todo)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(context.Background(), todo.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, todo.ID, retrieved.ID)
	assert.Equal(t, todo.Description, retrieved.Description)
	assert.Equal(t, todo.FileID, retrieved.FileID)

	// Test not found
	_, err = repo.GetByID(context.Background(), uuid.New().String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

