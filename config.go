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
		log.Printf("Error loading .env file: %s. Proceeding with system environment variables.\n", err)
	}

	config := AppConfig{
		DBConnectionString: getEnv("DB_CONNECTION_STRING", "user:password@tcp(localhost:3306)/dbname"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
	}

	return &config
}

func getOriginVar(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Environment variable %s is set: %s\n", key, value)
		return value
	}
	log.Printf("Environment variable %s is not set, using fallback: %s\n", key, fallback)
	return fallback
}