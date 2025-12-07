package main

import (
	"fmt"
	"log"
	"os"
	docs "what-went-wrong-api/cmd/docs"
	"what-went-wrong-api/internal/handlers"
	"what-went-wrong-api/internal/middleware"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/seed"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title What Went Wrong API
// @version 1.0
// @description API for what went wrong
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 環境変数から接続情報を取得
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	// DSNを構築
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", dbHost, dbUser, dbPassword, dbName, dbPort)
	fmt.Println("DSN:", dsn) // Prevent unused variable error

	// GORMでデータベースに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// データベースにテーブルを作成
	db.AutoMigrate(
		&models.User{},
		&models.Goal{},
		&models.ExcuseEntry{},
		&models.ExcuseTemplate{},
		&models.UserPlan{},
	)

	// 開発環境でのみ初期データをシード
	if os.Getenv("APP_ENV") != "production" {
		if err := seed.Run(db); err != nil {
			log.Printf("Warning: Failed to seed database: %v", err)
		}
	}

	// Auth Middlewareの初期化
	authMiddleware, err := middleware.NewAuthMiddleware()
	if err != nil {
		log.Fatalf("Failed to initialize auth middleware: %v", err)
	}

	// Serviceの初期化
	entitlementService := services.NewEntitlementService(db)
	planHandler := handlers.NewPlanHandler(entitlementService)

	// AIHandler も復活させる (別Issueの実装だが、依存関係として EntitlementService が必要だったため)
	// ただし、AIHandler のソースコードは別ブランチだが、main にマージされている前提で書くか？
	// ユーザーは main ブランチに戻したと言っていた。
	// ここでは Issue #6 の範囲である「プラン管理」に集中するが、
	// ユーザーが「Issue #6 を実装」と言った時、Issue #7 のコードは消えている状態 (main ブランチ)。
	// なので、AI関連のコードは参照できない可能性がある。
	// しかし、先ほどのステップで AI関連のファイルは削除されたと出ている。
	// したがって、AIの実装は入れない。
	// Planの実装だけを入れる。

	// Ginエンジンのインスタンスを作成
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	v1.Use(authMiddleware)
	{
		v1.GET("/users", handlers.GetUsers(db))
		v1.GET("/me/plan", planHandler.GetMePlan)
		v1.POST("/me/plan", planHandler.PostMePlan)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
