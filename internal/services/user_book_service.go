package services

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/internal/repositories"
	"book-management-system/pkg/logger"
	"github.com/google/uuid"
)

type UserBookService struct {
	repo *repositories.UserBookRepository
	log  *logger.Logger
}

// NewUserBookService создает новый сервис
func NewUserBookService(repo *repositories.UserBookRepository) *UserBookService {
	return &UserBookService{
		repo: repo,
		log:  logger.GetLogger(),
	}
}

// AddBookToUser добавляет книгу в список пользователя
func (s *UserBookService) AddBookToUser(userID, bookID uuid.UUID) error {
	err := s.repo.AddUserBook(userID, bookID)
	if err != nil {
		s.log.Warnf("Ошибка добавления книги пользователю: %v", err)
		return err
	}
	return nil
}

// UpdateReadingProgress обновляет статус и прогресс чтения
func (s *UserBookService) UpdateReadingProgress(userID, bookID uuid.UUID, status models.ReadingStatus, pagesRead int) error {
	err := s.repo.UpdateReadingProgress(userID, bookID, status, pagesRead)
	if err != nil {
		s.log.Warnf("Ошибка обновления прогресса чтения: %v", err)
		return err
	}
	return nil
}

// RemoveBookFromUser удаляет книгу из списка пользователя
func (s *UserBookService) RemoveBookFromUser(userID, bookID uuid.UUID) error {
	err := s.repo.RemoveUserBook(userID, bookID)
	if err != nil {
		s.log.Warnf("Ошибка удаления книги из списка пользователя: %v", err)
		return err
	}
	return nil
}

// GetUserBooks получает список книг пользователя
func (s *UserBookService) GetUserBooks(userID uuid.UUID) ([]dto.UserBookResponse, error) {
	userBooks, err := s.repo.GetUserBooks(userID)
	if err != nil {
		s.log.Warnf("Ошибка получения списка книг пользователя: %v", err)
		return nil, err
	}

	bookResponses := make([]dto.UserBookResponse, 0, len(userBooks))

	for _, book := range userBooks {
		bookResponses = append(bookResponses, dto.UserBookResponse{
			BookID:    book.BookID,
			Status:    string(book.Status),
			PagesRead: book.PagesRead,
			CreatedAt: book.CreatedAt,
			UpdatedAt: book.UpdatedAt,
		})
	}

	return bookResponses, nil
}
