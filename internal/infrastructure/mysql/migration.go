package mysql

import (
	"github.com/ar-agahian/ice-assignment/internal/domain"
	"gorm.io/gorm"
)

// RunMigrations runs GORM AutoMigrate to create/update database schema
func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.TodoItem{}); err != nil {
		return err
	}
	return nil
}

