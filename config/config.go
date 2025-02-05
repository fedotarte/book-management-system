package config

import (
	"book-management-system/pkg/logger"
	"github.com/joho/godotenv"
	"os"
)

var log = logger.GetLogger()

// Load загружает переменные окружения из .env файла в зависимости от среды
func Load() {
	env := os.Getenv("APP_ENV")
	var envFile string

	switch env {
	case "production":
		envFile = ".env.production"
	case "staging":
		envFile = ".env.staging"
	default:
		envFile = ".env.development"
	}

	log.Infof("Загрузка .env файла: %s", envFile)

	if err := godotenv.Load(envFile); err != nil {
		log.Warnf("Не удалось загрузить %s, используем переменные окружения по умолчанию", envFile)
	}

	log.Infof("JWT_SECRET после загрузки: %s", os.Getenv("JWT_SECRET"))

}

// GetEnv возвращает значение переменной окружения или значение по умолчанию
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {

		log.Infof("required value %s by key: %s exists", value, key)

		return value
	}

	log.Infof("required value by key: %s not exists", key)

	return defaultValue
}
