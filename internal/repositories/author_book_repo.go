package repositories

import (
	database2 "book-management-system/internal/database"
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookAuthorRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewBookAuthorRepository() *BookAuthorRepository {
	return &BookAuthorRepository{
		db:  database2.DB,
		log: logger.GetLogger(),
	}
}

// GetAuthorsByBookID получает авторов книги через book_authors
func (r *BookAuthorRepository) GetAuthorIDsByBookID(bookID uuid.UUID) ([]uuid.UUID, error) {
	var authorIDs []uuid.UUID

	err := r.db.Model(&models.BookAuthor{}).
		Where("book_id = ?", bookID).
		Pluck("author_id", &authorIDs).Error

	if err != nil {
		r.log.Warnf("Ошибка получения author_id для книги %v: %v", bookID, err)
		return nil, err
	}

	return authorIDs, nil
}
func (r *BookAuthorRepository) GetBookIDsByAuthorID(authorID uuid.UUID) ([]uuid.UUID, error) {
	var bookIDs []uuid.UUID

	err := r.db.Model(&models.BookAuthor{}).
		Where("author_id = ?", authorID).
		Pluck("book_id", &bookIDs).Error

	if err != nil {
		r.log.Warnf("Ошибка получения author_id для книги %v: %v", authorID, err)
		return nil, err

	}

	return bookIDs, nil
}

func (r *BookAuthorRepository) GetAuthorsForBooks(bookIDs []uuid.UUID) ([]models.BookAuthor, error) {
	var bookAuthors []models.BookAuthor

	err := r.db.Model(&models.BookAuthor{}).
		Where("book_id IN (?)", bookIDs).
		Find(&bookAuthors).Error

	if err != nil {
		r.log.Warnf("Ошибка получения связей книга-автор: %v", err)
		return nil, err
	}

	return bookAuthors, nil
}

func (r *BookAuthorRepository) GetBooksForAuthors(authorIds []uuid.UUID) ([]models.BookAuthor, error) {
	var bookAuthors []models.BookAuthor

	err := r.db.Model(&models.BookAuthor{}).
		Where("author_id IN (?)", authorIds).
		Find(&bookAuthors).Error

	if err != nil {
		r.log.Warnf("Ошибка получения связей автор-книга: %v", err)
		return nil, err
	}

	return bookAuthors, nil
}
