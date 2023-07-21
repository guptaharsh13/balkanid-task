package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     uint
	SSLMode  string
}

func LoadConfig() *Config {

	config := Config{
		DB: DBConfig{
			Host:     genEnv("DB_HOST", "localhost"),
			User:     genEnv("DB_USER", "postgres"),
			Password: genEnv("DB_PASS", "postgres"),
			Name:     genEnv("DB_NAME", "balkanid-task"),
			Port:     getEnvAsUint("DB_PORT", 3000),
			SSLMode:  genEnv("DB_SSLMODE", "disable"),
		},
	}
	return &config
}

func genEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsUint(key string, defaultValue uint) uint {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	uintValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return uint(uintValue)
}
