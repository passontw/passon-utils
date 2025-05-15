package utils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGormDB 建立並回傳一個 GORM PostgreSQL 連線
func NewGormDB(cfg AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("GORM 連線失敗: %w", err)
	}
	return db, nil
}

// Create 新增一筆資料
func Create[T any](db *gorm.DB, obj *T) error {
	return db.Create(obj).Error
}

// GetByID 取一筆資料
func GetByID[T any](db *gorm.DB, id any, obj *T) error {
	return db.First(obj, id).Error
}

// GetAll 取全部資料
func GetAll[T any](db *gorm.DB, out *[]T) error {
	return db.Find(out).Error
}

// GetPaginated 分頁查詢
func GetPaginated[T any](db *gorm.DB, page, pageSize int, out *[]T) (int64, error) {
	var total int64
	db.Model(out).Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(out).Error
	return total, err
}

// Update 更新資料
func Update[T any](db *gorm.DB, obj *T) error {
	return db.Save(obj).Error
}

// DeleteByID 刪除資料
func DeleteByID[T any](db *gorm.DB, id any) error {
	var obj T
	return db.Delete(&obj, id).Error
}
