package repositories

import (
	database2 "book-management-system/internal/database"
	"book-management-system/internal/models"
	"github.com/google/uuid"

	"book-management-system/pkg/logger"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAuthorRepository() *AuthorRepository {
	return &AuthorRepository{
		db:  database2.DB,
		log: logger.GetLogger(),
	}
}

func (r *AuthorRepository) CreateAuthor(author *models.Author, bookIds []uuid.UUID) error {
	tx := r.db.Begin()

	if err := tx.Create(author).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("create author error: %v", err)
	}

	for _, bookId := range bookIds {
		authorBooks := models.BookAuthor{BookID: &bookId, AuthorID: &author.ID}
		if err := tx.Create(&authorBooks).Error; err != nil {
			tx.Rollback()
			r.log.Warnf("Ошибка связывания книг с автороми: %v", err)
			return err
		}
	}

	return tx.Commit().Error
}

func (r *AuthorRepository) GetAuthorByID(authorID uuid.UUID) (*models.Author, error) {
	var author models.Author
	err := r.db.Preload("Books").
		Where("id = ?", authorID).
		First(&author).Error

	if err != nil {
		r.log.Warnf("Ошибка получения автора по ID: %v", err)
		return nil, err
	}

	return &author, nil
}

func (r *AuthorRepository) UpdateAuthor(author *models.Author, bookIDs []uuid.UUID) error {
	tx := r.db.Begin()

	if err := tx.Save(author).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("update author error: %v", err)
	}

	if err := tx.Where("author_id = ?", author.ID).Delete(&models.BookAuthor{}).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("delete old author-book relations error: %v", err)
		return err
	}

	for _, bookId := range bookIDs {
		booAuthor := models.BookAuthor{BookID: &bookId, AuthorID: &author.ID}
		if err := tx.Create(&booAuthor).Error; err != nil {
			tx.Rollback()
			r.log.Warnf("Ошиька связывания автора с книгами: %v", err)
			return err
		}
	}

	return tx.Commit().Error

}

func (r *AuthorRepository) DeleteAuthor(authorID uuid.UUID) error {
	tx := r.db.Begin()

	if err := tx.Model(&models.Author{}).
		Where("author_id = ?", authorID).
		Update("author_id", nil).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка обнуления связи книги с авторами: %v", err)
		return err
	}

	if err := tx.Where("id = ?", authorID).
		Delete(&models.Author{}).Error; err != nil {
		tx.Rollback()
		r.log.Warnf("Ошибка удаления автора: %v", err)
		return err
	}

	return tx.Commit().Error

}

func (r *AuthorRepository) GetAuthorsPaginated(limit int, afterID *uuid.UUID) ([]models.Author, error) {
	if limit <= 0 {
		limit = 10
	}
	var authors []models.Author

	query := r.db.Model(&models.Author{}).
		Order("created_at ASC").
		Limit(limit)

	if afterID != nil {
		query = query.
			Where("created_at > (?)", r.db.Model(&models.Author{}).
				Select("created_at").
				Where("id = ?", *afterID))
	}

	err := query.Find(&authors).Error

	if err != nil {
		r.log.Warnf("Ошибка получения списка авторов: %v", err)
		return nil, err
	}

	return authors, nil

}

func (r *AuthorRepository) GetAuthorsByIDs(authorIDs []uuid.UUID) ([]models.Author, error) {
	var authors []models.Author

	err := r.db.Model(&models.Author{}).
		Where("id IN (?)", authorIDs).
		Find(&authors).Error

	if err != nil {
		r.log.Warnf("Ошибка получения авторов: %v", err)
		return nil, err
	}

	return authors, nil
}
