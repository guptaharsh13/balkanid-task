package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/controllers"
	"github.com/guptaharsh13/balkanid-task/middleware"
)

func UserRouter(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("/signup", controllers.Signup)
		users.POST("/login", controllers.Login)
		users.POST("/verify/:username", controllers.VerifyEmail)
		users.GET("/activate/:username/:code", controllers.ActivateUser)
		users.POST("/deactivate/:username", middleware.IsAdmin, controllers.DeactivateUser)
		users.GET("/me", middleware.RequireAuth, controllers.GetCurrentUser)
		users.GET("/", middleware.IsAdmin, controllers.GetUsers)
		users.GET("/:username", middleware.IsAdmin, controllers.GetUserByUsername)
		users.DELETE("/:username", middleware.IsAdmin, controllers.DeleteUser)
		users.POST("/upload", middleware.IsAdmin, controllers.BulkUploadUsers)
	}
}
