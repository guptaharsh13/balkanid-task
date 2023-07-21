package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/controllers"
)

func UserRouter(r *gin.Engine) {
	tasks := r.Group("/users")
	tasks.POST("/signup", controllers.Signup)
	tasks.GET("/login", controllers.Login)
}
