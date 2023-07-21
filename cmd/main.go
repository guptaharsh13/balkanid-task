package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/config"
	"github.com/guptaharsh13/balkanid-task/controllers"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/routes"
	"github.com/guptaharsh13/balkanid-task/utils"
)

func init() {
	initializers.LoadEnvVariables()

	config := config.LoadConfig()
	initializers.ConnectToDb(config.DB)
	initializers.SyncDatabase()
	initializers.SyncPermissions()
	err := utils.SetupValidator()
	if err != nil {
		fmt.Println("Couldn't setup validator")
	}
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	routes.TaskRouter(r)

	r.Run()
}
