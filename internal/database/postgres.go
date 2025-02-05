package database

import (
	"book-management-system/config"
	"book-management-system/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitDatabase инициализирует подключение к PostgreSQL
func InitPostgres() {
	dsn := config.GetEnv("DATABASE_URL", "host=localhost user=postgres password=secret dbname=books port=5432 sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	log.Println("База данных успешно подключена")

	// Автомиграция моделей
	db.AutoMigrate(
		&models.Author{},
		&models.Book{},
		&models.BookRating{},
		&models.BookAuthor{},
		&models.ModeratorAction{},
		&models.ReadingProgress{},
		&models.RefreshToken{},
		&models.User{},
		&models.UserBook{},
	)

	DB = db
}
