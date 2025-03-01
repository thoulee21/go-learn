package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/thoulee21/go-learn/controllers"
)

func SetupChatRoutes(r *gin.Engine, cc *controllers.ChatController) {
    chatGroup := r.Group("/chat")
    {
        chatGroup.POST("", cc.Chat)
        chatGroup.GET("/history/:session_id", cc.GetChatHistory)
        chatGroup.GET("/test", cc.Test)
    }
}