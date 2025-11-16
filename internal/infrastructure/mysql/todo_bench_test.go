package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/ar-agahian/ice-assignment/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func BenchmarkTodoRepository_Create(b *testing.B) {
	dsn := "root:password@tcp(localhost:3306)/todo_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		b.Skipf("Skipping benchmark: failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		b.Skipf("Skipping benchmark: failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		b.Skipf("Skipping benchmark: failed to ping database: %v", err)
	}

	// AutoMigrate to create tables
	db.AutoMigrate(&domain.TodoItem{})
	defer db.Exec("DROP TABLE IF EXISTS todo_items")

	repo := NewTodoRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		todo := domain.NewTodoItem(
			"Benchmark description",
			time.Now().Add(24*time.Hour),
			"file-123",
		)
		_ = repo.Create(ctx, todo)
	}
}

