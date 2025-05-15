# passon-utils

## 專案簡介

本專案為一套以 Go 語言開發，結合 Nacos 配置中心、PostgreSQL（GORM）、Redis（單點/Cluster）等常用服務的高可維護性工具模組，並支援 fx 依賴注入，方便快速建構現代化後端服務。

---

## 依賴安裝

```bash
go get github.com/passontw/passon-utils
```

---

## 設定說明

### 1. .env 基本參數（Nacos 連線）

```
NACOS_HOST=127.0.0.1
NACOS_PORT=8848
NACOS_NAMESPACE=public
NACOS_USER=nacos
NACOS_PASSWORD=nacos
NACOS_DATAID=your_config.json
NACOS_GROUP=DEFAULT_GROUP
NACOS_SERVICE=your-service
NACOS_IP=127.0.0.1
NACOS_SERVICE_PORT=8080
```

### 2. Nacos 設定中心 config 範例

```json
{
  "PORT": "8080",
  "DB_HOST": "127.0.0.1",
  "DB_PORT": 5432,
  "DB_NAME": "postgres",
  "DB_USER": "postgres",
  "DB_PASSWORD": "a12345678",
  "REDIS_HOST": "127.0.0.1",
  "REDIS_PORT": 6379,
  "REDIS_USERNAME": "",
  "REDIS_PASSWORD": "",
  "REDIS_DB": 0
}
```
- 若為 Redis Cluster，`REDIS_HOST` 可設為多個逗號分隔的 host，例如：`"REDIS_HOST": "10.0.0.1,10.0.0.2,10.0.0.3"`

---

## 啟動方式

```go
import (
    "github.com/passontw/passon-utils/utils"
    "go.uber.org/fx"
)

func main() {
    app := fx.New(
        utils.Module(),
        fx.Invoke(YourHandler), // 你可直接注入 *gorm.DB, redis.UniversalClient, utils.AppConfig
    )
    app.Run()
}
```

---

## PostgreSQL CRUD/分頁查詢範例

```go
import (
    "github.com/passontw/passon-utils/utils"
    "gorm.io/gorm"
)

func Example(db *gorm.DB) {
    // 新增
    user := &utils.User{Name: "Tom", Phone: "0912345678", Password: "pw"}
    err := utils.CreateUser(db, user)

    // 取一筆
    u, err := utils.GetUserByID(db, user.ID)

    // 取全部
    users, err := utils.GetAllUsers(db)

    // 分頁查詢
    page, pageSize := 1, 10
    users, total, err := utils.GetUsersPaginated(db, page, pageSize)

    // 更新
    user.Name = "Tommy"
    err = utils.UpdateUser(db, user)

    // 刪除
    err = utils.DeleteUser(db, user.ID)
}
```

---

## Redis 單點/Cluster 範例

```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/passontw/passon-utils/utils"
)

func Example(rdb redis.UniversalClient) {
    err := rdb.Set(ctx, "hello", "world", 0).Err()
    val, err := rdb.Get(ctx, "hello").Result()
}
```
- 若 Nacos config 設定多個 host，會自動使用 Cluster 模式。

---

## 擴充教學

### 1. PostgreSQL 自定義 Model 與 CRUD

你可以在自己的專案中自定義 struct，例如：

```go
type User struct {
    ID   int64  `gorm:"column:id"`
    Name string `gorm:"column:name"`
    // ...其他欄位
}
```

然後直接使用泛型 CRUD：

```go
user := &User{Name: "Tom"}
utils.Create(db, user)

var u User
utils.GetByID(db, user.ID, &u)

var users []User
utils.GetAll(db, &users)

total, err := utils.GetPaginated(db, 1, 10, &users)

user.Name = "Tommy"
utils.Update(db, &user)
utils.DeleteByID[User](db, user.ID)
```

### 2. Redis 擴充與 Cluster 支援

- 只要在 Nacos config 設定多個 host（逗號分隔），即可自動切換為 Cluster 模式。
- 你可直接用 `redis.UniversalClient` 進行所有 Redis 操作。

```go
func Example(rdb redis.UniversalClient) {
    rdb.Set(ctx, "foo", "bar", 0)
    val, _ := rdb.Get(ctx, "foo").Result()
}
```

### 3. Nacos config 動態擴充

- 只要在 Nacos 設定中心新增欄位，並於 `AppConfig` struct 增加對應欄位即可。
- 例如：
```json
{
  "NEW_FEATURE": true,
  "SOME_LIMIT": 100
}
```

```go
type AppConfig struct {
    // ...原有欄位
    NewFeature bool `json:"NEW_FEATURE"`
    SomeLimit  int  `json:"SOME_LIMIT"`
}
```

---

## 測試建議

- 可用 Go 標準 testing 套件撰寫單元測試
- 建議針對 CRUD、分頁、Redis cluster 切換等情境撰寫測試

---

## 聯絡/貢獻

如有建議或需求，歡迎 PR 或 issue！ 