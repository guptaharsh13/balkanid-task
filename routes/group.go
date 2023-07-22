package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/controllers"
	"github.com/guptaharsh13/balkanid-task/middleware"
)

func GroupRouter(r *gin.Engine) {
	groups := r.Group("/groups")
	groups.Use(middleware.IsAdmin)
	{
		groups.POST("/", controllers.CreateGroup)
		groups.GET("/", controllers.GetGroups)
		groups.GET("/:name", controllers.GetGroupByName)
		groups.PUT("/:name", controllers.UpdateGroupPut)
		groups.PATCH("/:name", controllers.UpdateGroupPatch)
		groups.DELETE("/:name", controllers.DeleteGroup)
		groups.POST("/:name/permissions", controllers.AddPermissionsToGroup)
		groups.POST("/:name/users", controllers.AddUsersToGroup)
		groups.GET("/:name/users", controllers.GetUsersByGroup)
	}
}
