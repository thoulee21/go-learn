package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thoulee21/go-learn/controllers"
	_ "github.com/thoulee21/go-learn/docs"
	"github.com/thoulee21/go-learn/models"
	"github.com/thoulee21/go-learn/routes"
	"github.com/thoulee21/go-learn/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading from environment")
	} else {
		log.Println("Loading .env file")
	}

	// 连接数据库
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	const database = "aichatbot"

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	dbErr := db.AutoMigrate(&models.ChatMessage{})
	if dbErr != nil {
		panic("failed to migrate database")
	}
}

func main() {
	r := gin.Default()

	aiService, err := services.NewAIService()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize AI service: %v", err))
	}

	chatController := &controllers.ChatController{DB: db, AIService: aiService}
	routes.SetupChatRoutes(r, chatController)

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 默认在8080端口启动服务
	if err := r.Run(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
