package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/controllers"
)

func TaskRouter(r *gin.Engine) {
	tasks := r.Group("/tasks")
	tasks.POST("/", controllers.CreateTask)
	tasks.GET("/", controllers.GetTasks)
}
