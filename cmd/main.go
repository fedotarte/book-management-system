package main

import (
	"book-management-system/config"
	_ "book-management-system/docs"
	"book-management-system/internal/database"
	"book-management-system/internal/routes"
	"book-management-system/pkg/logger"
)

// @title						Book Management API
// @version					1.0
// @description				API для управления книгами, пользователями и отзывами.
// @host						localhost:8080
// @BasePath					/api/v1/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	logger.InitLogger()
	log := logger.GetLogger()

	log.Info("Загрузка конфигурации...")
	config.Load()

	log.Info("Подключение постргри...")
	database.InitPostgres()

	log.Info("Подключение монги...")
	database.InitMongoDB()

	port := config.GetEnv("SERVER_PORT", "8080")

	log.Info("Инициализация маршрутизатора...")
	router := routes.InitRouter()

	log.Infof("Запуск сервера на порту %s...", port)
	err := router.Run(":" + port)

	if err != nil {
		log.Fatal(err)
	}

}
