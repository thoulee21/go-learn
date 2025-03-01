package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thoulee21/go-learn/controllers"
	_ "github.com/thoulee21/go-learn/docs"
	"github.com/thoulee21/go-learn/models"
	"github.com/thoulee21/go-learn/routes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})
}

func main() {
	r := gin.Default()

	taskController := &controllers.TaskController{DB: db}

	routes.SetupTaskRoutes(r, taskController)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":80")
}
