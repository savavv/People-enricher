package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var config *Config
var once sync.Once

// LoadConfig инициализирует конфигурацию из .env
func LoadConfig() *Config {
	once.Do(func() {
		// Загружаем .env файл
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using system env variables")
		}

		config = &Config{
			Port:       getEnv("PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", "password"),
			DBName:     getEnv("DB_NAME", "people_db"),
		}
	})

	return config
}

// Вспомогательная функция для чтения переменных окружения с дефолтом
func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
