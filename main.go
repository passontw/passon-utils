package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/passontw/passon-utils/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// User 對應 users 資料表
type User struct {
	ID               int64      `gorm:"column:id"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at"`
	Name             string     `gorm:"column:name"`
	Phone            string     `gorm:"column:phone"`
	Password         string     `gorm:"column:password"`
	AvailableBalance float64    `gorm:"column:available_balance"`
	FrozenBalance    float64    `gorm:"column:frozen_balance"`
}

// 查詢並印出 admin 資料
func PrintAdmin(db *gorm.DB) {
	var user User
	if err := db.Where("phone = ?", "0987654321").First(&user).Error; err != nil {
		fmt.Println("查詢失敗:", err)
		return
	}
	fmt.Printf("Admin 資料: %+v\n", user)
}

func RunCLI(lc fx.Lifecycle, client *redis.Client) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := client.Set(ctx, "hello", "world", 0).Err()
			if err != nil {
				log.Fatalf("寫入失敗: %v", err)
			}
			log.Println("成功寫入 key: hello, value: world")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})
}

func main() {
	app := fx.New(
		utils.Module(),
		fx.Invoke(RunCLI),
		fx.Invoke(PrintAdmin),
	)
	app.Run()
}
