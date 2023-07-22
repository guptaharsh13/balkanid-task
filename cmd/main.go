package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/config"
	"github.com/guptaharsh13/balkanid-task/controllers"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/routes"
	"github.com/guptaharsh13/balkanid-task/utils"
	"github.com/spf13/cobra"
)

var configuration *config.Config

func init() {
	initializers.LoadEnvVariables()

	configuration = config.LoadConfig()
	initializers.ConnectToDb(configuration.DB)
	initializers.SyncDatabase()
	initializers.SyncPermissions()
	err := utils.SetupValidator()
	if err != nil {
		fmt.Println("❌ Couldn't setup validator")
	}
	fmt.Println("✅ Validator Setup")
}

func StartApp() error {
	gin.SetMode(configuration.Environment)
	fmt.Printf("🚧 Starting %s Environment\n", configuration.Environment)

	r := gin.Default()

	r.Use(cors.Default())
	if err := r.SetTrustedProxies(configuration.TrustedProxies); err != nil {
		panic(err)
	}

	r.GET("/health", controllers.Health)
	routes.UserRouter(r)
	routes.TaskRouter(r)
	routes.GroupRouter(r)
	routes.RoleRouter(r)

	if err := r.Run(); err != nil {
		return fmt.Errorf("couldn't start the server: %s", err.Error())
	} else {
		fmt.Println("Server started")
	}
	return nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "start",
		Short: "Start App",
		Long:  `This start command will start the app.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := StartApp(); err != nil {
				panic(err)
			}
		},
	}

	rootCmd.AddCommand(adminCommand)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
