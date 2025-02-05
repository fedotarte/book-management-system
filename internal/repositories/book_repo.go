package repositories

import (
	database2 "book-management-system/internal/database"
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type BookRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewBookRepository xсоздает новый репозиторий для книг
func NewBookRepository() *BookRepository {
	return &BookRepository{
		db:  database2.DB,
		log: logger.GetLogger(),
	}
}

// CreateBook создает новую книгу и связывает с авторами
func (r *BookRepository) CreateBook(book *models.Book, authorIDs []uuid.UUID) error {
	tx := r.db.Begin()

	// Создаем книгу
	if err := tx.Create(book).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка создания книги: %v", err)
		return err
	}

	// Связываем с авторами через book_authors
	for _, authorID := range authorIDs {
		bookAuthor := models.BookAuthor{BookID: &book.ID, AuthorID: &authorID}
		if err := tx.Create(&bookAuthor).Error; err != nil {
			tx.Rollback()
			r.log.Warnf("Ошибка связывания книги с авторамм: %v", err)
			return err
		}
	}

	return tx.Commit().Error
}

// GetBookByID получает книгу по ID, с возможной фильтрацией по `confirmed`
func (r *BookRepository) GetBookByID(bookID uuid.UUID, onlyConfirmed bool) (*models.Book, error) {
	var book models.Book

	query := r.db.Model(&models.Book{}).Where("id = ?", bookID)

	// Если `onlyConfirmed == true`, добавляем фильтр `confirmed = true`
	if onlyConfirmed {
		query = query.Where("confirmed = ?", true)
	}

	err := query.First(&book).Error
	if err != nil {
		r.log.Warnf("Ошибка получения книги по ID: %v", err)
		return nil, err
	}

	return &book, nil
}

// UpdateBook обновляет данные книги и связи с авторами
func (r *BookRepository) UpdateBook(book *models.Book, authorIDs []uuid.UUID) error {
	tx := r.db.Begin()

	// Обновляем книгу
	if err := tx.Save(book).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка обновления книги: %v", err)
		return err
	}

	// Удаляем старые связи
	if err := tx.Where("book_id = ?", book.ID).Delete(&models.BookAuthor{}).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка удаления старых связей книги: %v", err)
		return err
	}

	// Добавляем новые связи
	for _, authorID := range authorIDs {
		bookAuthor := models.BookAuthor{BookID: &book.ID, AuthorID: &authorID}
		if err := tx.Create(&bookAuthor).Error; err != nil {
			tx.Rollback()
			r.log.Warnf("Ошибка связывания книги с автором: %v", err)
			return err
		}
	}

	return tx.Commit().Error
}

// DeleteBook удаляет книгу и связи с авторами (ставит NULL)
func (r *BookRepository) DeleteBook(bookID uuid.UUID) error {
	tx := r.db.Begin()

	// 1. Удаляем связи в book_authors (ставим NULL)
	if err := tx.Model(&models.BookAuthor{}).
		Where("book_id = ?", bookID).
		Update("book_id", nil).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка обнуления связи книги с авторами: %v", err)
		return err
	}

	// 2. Оставляем user_books, но book_id ставим NULL
	if err := tx.Model(&models.UserBook{}).
		Where("book_id = ?", bookID).
		Update("book_id", nil).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка обнуления связи книги с пользователями: %v", err)
		return err
	}

	// 3. Полностью удаляем оценки книги (book_ratings)
	if err := tx.Where("book_id = ?", bookID).Delete(&models.BookRating{}).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка удаления оценок книги: %v", err)
		return err
	}

	// 4. Удаляем саму книгу (soft delete)
	if err := tx.Where("id = ?", bookID).Delete(&models.Book{}).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка удаления книги: %v", err)
		return err
	}

	// 5. Удаляем отзывы из MongoDB
	_, err := database2.MongoDB.Database("bookstore").
		Collection("reviews").
		DeleteMany(context.TODO(), bson.M{"book_id": bookID})

	if err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка удаления отзывов в MongoDB: %v", err)
		return err
	}

	return tx.Commit().Error
}

// GetBooksPaginated получает книги с маркерной пагинацией
func (r *BookRepository) GetBooksPaginated(limit int, afterID *uuid.UUID) ([]models.Book, error) {
	if limit <= 0 {
		limit = 10
	}

	var books []models.Book
	query := r.db.Model(&models.Book{}).
		Where("confirmed = ?", true).
		Order("created_at ASC").
		Limit(limit)

	// Если есть afterID, загружаем книги после указанной
	if afterID != nil {
		query = query.Where("created_at > (?)", r.db.Model(&models.Book{}).
			Select("created_at").
			Where("id = ?", *afterID))
	}

	err := query.Find(&books).Error
	if err != nil {
		r.log.Warnf("Ошибка получения списка книг: %v", err)
		return nil, err
	}

	return books, nil
}

// UpdateBookRating обновляет средний рейтинг книги в PostgreSQL
func (r *BookRepository) UpdateBookRating(bookID uuid.UUID, averageRating float64) error {
	err := r.db.Model(&models.Book{}).
		Where("id = ?", bookID).
		Update("average_rating", averageRating).Error

	if err != nil {
		r.log.Warnf("Ошибка обновления среднего рейтинга книги %s: %v", bookID, err)
		return err
	}

	return nil
}

// UpdateBookCover обновляет путь к обложке в базе
func (r *BookRepository) UpdateBookCover(bookID string, coverPath string) error {
	err := r.db.Model(&models.Book{}).
		Where("id = ?", bookID).
		Update("cover_image", coverPath).Error

	if err != nil {
		r.log.Warnf("Ошибка обновления обложки книги: %v", err)
		return err
	}
	return nil
}

func (r *BookRepository) GetBooksByIds(bookIDs []uuid.UUID) ([]models.Book, error) {
	var books []models.Book

	err := r.db.Model(&models.Book{}).
		Where("id IN (?)", bookIDs).
		Find(&books).Error

	if err != nil {
		r.log.Warnf("Ошибка получения книг по множественным id: %v", err)
		return nil, err
	}

	return books, nil
}
