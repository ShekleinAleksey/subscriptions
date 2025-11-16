package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	User     string
	Host     string
	Port     string
	DBName   string
	Password string
	SSLMode  string
}

func LoadConfig() Config {
	err := godotenv.Load("example.env")
	if err != nil {
		log.Println("error loading env variables")
	}

	return Config{
		User:     getEnv("DB_USER", "admin"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		DBName:   getEnv("DB_NAME", "subscriptiondb"),
		Password: getEnv("DB_PASSWORD", "root123"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv — вспомогательная функция для получения переменной окружения или значения по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
