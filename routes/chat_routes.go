package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thoulee21/go-learn/controllers"
	"github.com/thoulee21/go-learn/middlewares"
)

func SetupChatRoutes(r *gin.Engine, cc *controllers.ChatController) {
	chatGroup := r.Group("/chat")
	chatGroup.Use(middlewares.AuthJWTMiddleware())
	{
		chatGroup.POST("", cc.Chat)
		chatGroup.POST("/stream", cc.StreamChat)
		chatGroup.GET("/history/:session_id", cc.GetChatHistory)
	}

	r.GET("/test", cc.Test)
}
