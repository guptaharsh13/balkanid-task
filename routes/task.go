package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/controllers"
	"github.com/guptaharsh13/balkanid-task/middleware"
)

func TaskRouter(r *gin.Engine) {
	tasks := r.Group("/tasks")
	{
		tasks.POST("/", middleware.RequireAuth, controllers.CreateTask)
		tasks.GET("/", middleware.IsAdmin, controllers.GetTasks)
		tasks.GET("/:id", middleware.IsAdmin, controllers.GetTaskByID)
		tasks.DELETE("/:id", middleware.RequireAuth, controllers.DeleteTask)
		tasks.POST("/upload", middleware.IsAdmin, controllers.BulkUploadTasks)
		tasks.POST("/:id/asignees", middleware.RequireAuth, controllers.AssignTaskToUsers)
	}
}
