package database

import (
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdminUser создает тестового пользователя-админа
func SeedAdminUser(db *gorm.DB) {
	log := logger.GetLogger()

	// Проверяем, есть ли админ
	var existingAdmin models.User
	err := db.Where("role = ?", models.RoleAdmin).First(&existingAdmin).Error

	if err == gorm.ErrRecordNotFound {
		// Создаём нового
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		admin := models.User{
			ID:       uuid.New(),
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			Role:     models.RoleAdmin,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Ошибка при создании администратора: %v", err)
		} else {
			log.Info("Администратор успешно создан: admin@example.com / admin123")
		}
	} else if err != nil {
		log.Warnf("Ошибка при проверке администратора: %v", err)
	} else {
		log.Info("Администратор уже существует, пропускаем создание.")
	}
}
