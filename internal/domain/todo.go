package domain

import (
	"time"

	"github.com/google/uuid"
)

// TodoItem represents a todo item in the domain
type TodoItem struct {
	ID          uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Description string    `gorm:"type:varchar(500);not null"`
	DueDate     time.Time `gorm:"type:timestamp;not null;index"`
	FileID      string    `gorm:"type:varchar(255)"` // Reference to file stored in S3
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (TodoItem) TableName() string {
	return "todo_items"
}

// NewTodoItem creates a new TodoItem with a generated UUID
func NewTodoItem(description string, dueDate time.Time, fileID string) *TodoItem {
	return &TodoItem{
		ID:          uuid.New(),
		Description: description,
		DueDate:     dueDate,
		FileID:      fileID,
	}
}

