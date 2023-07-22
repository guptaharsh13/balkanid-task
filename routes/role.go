package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/controllers"
	"github.com/guptaharsh13/balkanid-task/middleware"
)

func RoleRouter(r *gin.Engine) {
	roles := r.Group("/roles")
	roles.Use(middleware.IsAdmin)
	{
		roles.POST("/", controllers.CreateRole)
		roles.GET("/", controllers.GetRoles)
		roles.GET("/:name", controllers.GetRoleByName)
		roles.PUT("/:name", controllers.UpdateRolePut)
		roles.PATCH("/:name", controllers.UpdateRolePatch)
		roles.DELETE("/:name", controllers.DeleteRole)
		roles.POST("/:name/permissions", controllers.AddPermissionsToRole)
		roles.POST("/:name/users", controllers.AssignRoleToUser)
		roles.GET("/:name/users", controllers.GetUsersByRole)
	}
}
