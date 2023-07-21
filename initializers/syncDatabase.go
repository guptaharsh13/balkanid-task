package initializers

import "github.com/guptaharsh13/balkanid-task/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Task{})
}
