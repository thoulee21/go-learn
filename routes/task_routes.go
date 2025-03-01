package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thoulee21/go-learn/controllers"
)

func SetupTaskRoutes(r *gin.Engine, tc *controllers.TaskController) {
	r.GET("/tasks", tc.GetTasks)
	r.POST("/tasks", tc.CreateTask)
	r.GET("/tasks/:id", tc.GetTask)
	r.PUT("/tasks/:id", tc.UpdateTask)
	r.DELETE("/tasks/:id", tc.DeleteTask)
}
