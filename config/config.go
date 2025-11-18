package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"sigs.k8s.io/yaml"
)

type Config struct {
	DB  DB  `yaml:"db"`
	Log Log `yaml:"log"`
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL"`
}

type DB struct {
	User     string `yaml:"user" env:"DB_USERNAME"`
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	DBName   string `yaml:"dbname" env:"DB_NAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE"`
}

func LoadConfig() (Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Читаем YAML конфиг
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	// Переопределяем значения из .env файла
	config.DB.User = getEnv("DB_USER", config.DB.User)
	config.DB.Host = getEnv("DB_HOST", config.DB.Host)
	config.DB.Port = getEnv("DB_PORT", config.DB.Port)
	config.DB.DBName = getEnv("DB_NAME", config.DB.DBName)
	config.DB.Password = getEnv("DB_PASSWORD", config.DB.Password)
	config.DB.SSLMode = getEnv("DB_SSLMODE", config.DB.SSLMode)
	config.Log.Level = getEnv("LOG_LEVEL", config.Log.Level)

	// Устанавливаем значения по умолчанию если пустые
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}

	return config, nil
}

// getEnv — вспомогательная функция для получения переменной окружения или значения по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
