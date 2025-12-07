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
	aiService := services.NewMockAIService()
	aiHandler := handlers.NewAIHandler(aiService)
	goalHandler := handlers.NewGoalHandler(db)
	excuseHandler := handlers.NewExcuseHandler(db)
	excuseTemplateHandler := handlers.NewExcuseTemplateHandler(db)

	// Middleware の初期化
	entitlementMiddleware := middleware.NewEntitlementMiddleware(entitlementService)

	// Ginエンジンのインスタンスを作成
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	v1.Use(authMiddleware)
	v1.Use(entitlementMiddleware)
	{
		v1.GET("/me/plan", planHandler.GetMePlan)
		v1.POST("/me/plan", planHandler.PostMePlan)
		v1.POST("/ai-excuse", aiHandler.PostAiExcuse)
		v1.GET("/goals", goalHandler.GetGoals)
		v1.POST("/goals", goalHandler.PostGoals)
		v1.PATCH("/goals/:id", goalHandler.PatchGoal)
		v1.DELETE("/goals/:id", goalHandler.DeleteGoal)
		v1.GET("/excuse-templates", excuseTemplateHandler.GetExcuseTemplates)
		v1.GET("/excuse-templates/:id", excuseTemplateHandler.GetExcuseTemplate)

		v1.GET("/goals/:goal_id/excuses", excuseHandler.GetExcuses)
		v1.GET("/goals/:goal_id/excuses/today", excuseHandler.GetExcuseToday)
		v1.POST("/goals/:goal_id/excuses", excuseHandler.PostExcuse)
		v1.PATCH("/excuses/:id", excuseHandler.PatchExcuse)
		v1.DELETE("/excuses/:id", excuseHandler.DeleteExcuse)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
