package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thoulee21/go-learn/controllers"
	_ "github.com/thoulee21/go-learn/docs"
	"github.com/thoulee21/go-learn/models"
	"github.com/thoulee21/go-learn/routes"
	"github.com/thoulee21/go-learn/services"
	"gorm.io/driver/sqlite"
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

	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.ChatMessage{})
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

	r.Run(":80")
}
