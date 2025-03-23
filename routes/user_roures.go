package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thoulee21/go-learn/controllers/user"
)

func SetupUserRoutes(r *gin.Engine, uc *user.UserController) {
	u := r.Group("/user")
	{
		u.GET("/", uc.GetAllUsers)
		u.POST("/", uc.NewUser)
		u.GET("/:id", uc.GetUserByID)
		u.PUT("/:id", uc.UpdateUser)
		u.DELETE("/:id", uc.DeleteUser)
	}
}
