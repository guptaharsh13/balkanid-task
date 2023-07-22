package initializers

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("❌ Couldn't load .env file")
	}
	fmt.Println("✅ .env Loaded")
}
