package repositories

import (
	"book-management-system/internal/database"
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBookRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUserBookRepository создает новый репозиторий
func NewUserBookRepository() *UserBookRepository {
	return &UserBookRepository{
		db:  database.DB,
		log: logger.GetLogger(),
	}
}

// AddUserBook добавляет книгу в список пользователя
func (r *UserBookRepository) AddUserBook(userID, bookID uuid.UUID) error {
	userBook := models.UserBook{
		UserID:    userID,
		BookID:    &bookID,
		Status:    models.StatusReading, // По умолчанию статус "Читаю"
		PagesRead: 0,
	}

	err := r.db.Create(&userBook).Error
	if err != nil {
		r.log.Warnf("Ошибка добавления книги в список пользователя: %v", err)
		return err
	}

	return nil
}

// UpdateReadingProgress обновляет статус чтения и количество прочитанных страниц
func (r *UserBookRepository) UpdateReadingProgress(userID, bookID uuid.UUID, status models.ReadingStatus, pagesRead int) error {
	err := r.db.Model(&models.UserBook{}).
		Where("user_id = ? AND book_id = ?", userID, bookID).
		Updates(models.UserBook{
			Status:    status,
			PagesRead: pagesRead,
		}).Error

	if err != nil {
		r.log.Warnf("Ошибка обновления прогресса чтения: %v", err)
		return err
	}

	return nil
}

// RemoveUserBook удаляет книгу из списка пользователя (soft delete)
func (r *UserBookRepository) RemoveUserBook(userID, bookID uuid.UUID) error {
	err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).
		Delete(&models.UserBook{}).Error

	if err != nil {
		r.log.Warnf("Ошибка удаления книги из списка пользователя: %v", err)
		return err
	}

	return nil
}

// GetUserBooks получает список книг пользователя
func (r *UserBookRepository) GetUserBooks(userID uuid.UUID) ([]models.UserBook, error) {
	var userBooks []models.UserBook

	err := r.db.Where("user_id = ?", userID).
		Find(&userBooks).Error

	if err != nil {
		r.log.Warnf("Ошибка получения списка книг пользователя: %v", err)
		return nil, err
	}

	return userBooks, nil
}
