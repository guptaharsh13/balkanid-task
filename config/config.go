package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Environment    string
	TrustedProxies []string
	DB             DBConfig
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     uint
	SSLMode  string
}

func findEnvironment() string {
	if flag.Lookup("test.v") == nil {
		env := os.Getenv("GO_ENV")
		if env == "production" || env == "prod" {
			return gin.ReleaseMode
		}
		return gin.DebugMode
	}
	return gin.TestMode
}

func LoadConfig() *Config {

	trustedProxies := getEnv("TRUSTED_PROXIES", "localhost,127.0.0.1")
	config := Config{
		Environment:    findEnvironment(),
		TrustedProxies: strings.Split(trustedProxies, ","),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "postgres"),
			Name:     getEnv("DB_NAME", "balkanid-task"),
			Port:     getEnvAsUint("DB_PORT", 3000),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
	fmt.Println("✅ Config Loaded")
	return &config
}

func getEnv(key string, defaultValue string) string {
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
