package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBConnectionString string
	ServerPort         string
}

func NewConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading configuration from environment")
	}

	config := AppConfig{
		DBConnectionString: getEnv("DB_CONNECTION_STRING", "user:password@tcp(localhost:3306)/dbname"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
	}

	return &config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}