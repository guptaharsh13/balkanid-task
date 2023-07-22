package initializers

import (
	"fmt"

	"github.com/guptaharsh13/balkanid-task/models"
)

func SyncDatabase() {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintf("Couldn't sync users table: %s", err))
	}
	if err := DB.AutoMigrate(&models.VerifyEmail{}); err != nil {
		panic(fmt.Sprintf("Couldn't sync verify_emails table: %s", err))
	}
	if err := DB.AutoMigrate(&models.Role{}); err != nil {
		panic(fmt.Sprintf("Couldn't sync roles table: %s", err))
	}
	if err := DB.AutoMigrate(&models.Permission{}); err != nil {
		panic(fmt.Sprintf("Couldn't sync permissions table: %s", err))
	}
	if err := DB.AutoMigrate(&models.Task{}); err != nil {
		panic(fmt.Sprintf("Couldn't sync tasks table: %s", err))
	}
	fmt.Println("✅ Synced Database")
}
