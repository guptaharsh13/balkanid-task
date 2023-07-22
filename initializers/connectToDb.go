package initializers

import (
	"fmt"

	"github.com/guptaharsh13/balkanid-task/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb(config config.DBConfig) {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.Host, config.User, config.Password, config.Name, config.Port, config.SSLMode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Couldn't connect to the database")
	}
	fmt.Println("✅ Database Connected!")
}
