package utils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGormDB 建立並回傳一個 GORM PostgreSQL 連線
func NewGormDB() (*gorm.DB, error) {
	dsn := "host=127.0.0.1 user=postgres password=a12345678 dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("GORM 連線失敗: %w", err)
	}
	return db, nil
}
