package main

import (
	"fmt"
	"strings"

	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/models"
	"github.com/guptaharsh13/balkanid-task/utils"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var adminCommand = &cobra.Command{
	Use:   "create admin",
	Short: "This command can be used to create an admin user.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("Invalid number of arguments")
			fmt.Println("Command: create admin <username> <email> <password>")
			return
		}

		username := args[0]
		email := args[1]
		password := args[2]

		fmt.Println("Creating admin user...")
		if len(strings.TrimSpace(username)) == 0 {
			fmt.Println("Username cannot be empty")
			return
		}
		if !utils.IsValidEmail(email) {
			fmt.Println("Invalid email address")
			return
		}
		if len(strings.TrimSpace(password)) == 0 {
			fmt.Println("Password cannot be empty")
			return
		}

		var user models.User
		if result := initializers.DB.Take(&user, "username = ?", username); result.RowsAffected > 0 {
			fmt.Println("Username already taken")
			return
		}
		if result := initializers.DB.Take(&user, "email = ?", username); result.RowsAffected > 0 {
			fmt.Println("Email already taken")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			fmt.Printf("Couldn't hash password: %s", err.Error())
			return
		}

		admin := models.User{
			Username: username,
			Email:    email,
			Password: string(hash),
			IsActive: true,
			IsAdmin:  true,
		}
		if result := initializers.DB.Create(&admin); result.Error != nil {
			fmt.Printf("Couldn't create admin user: %s", result.Error.Error())
			return
		}
		fmt.Println("Successfully created admin user")
	},
}
