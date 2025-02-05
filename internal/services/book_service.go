package services

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/internal/repositories"
	"book-management-system/pkg/logger"
	"github.com/google/uuid"
)

type BookService struct {
	authorRepository            *repositories.AuthorRepository
	bookRepository              *repositories.BookRepository
	bookAuthorMappingRepository *repositories.BookAuthorRepository
	log                         *logger.Logger
}

// NewBookService создает новый сервис для работы с книгами
func NewBookService(
	bookRepository *repositories.BookRepository,
	bookAuthorMappingRepository *repositories.BookAuthorRepository,
	authorRepository *repositories.AuthorRepository,
) *BookService {
	return &BookService{
		authorRepository:            authorRepository,
		bookRepository:              bookRepository,
		bookAuthorMappingRepository: bookAuthorMappingRepository,
		log:                         logger.GetLogger(),
	}
}

// CreateBook создает новую книгу и связывает с авторами
func (s *BookService) CreateBook(title, description, coverImage string, authorIDs []uuid.UUID, userRole string) (*models.Book, error) {
	book := &models.Book{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		CoverImage:  coverImage,
		Confirmed:   userRole == "moderator" || userRole == "admin", // Если создает модератор или админ, сразу подтверждаем
	}

	err := s.bookRepository.CreateBook(book, authorIDs)
	if err != nil {
		s.log.Warnf("Ошибка создания книги: %v", err)
		return nil, err
	}

	return book, nil
}

func (s *BookService) ConfirmBook(bookID uuid.UUID) error {
	book, err := s.bookRepository.GetBookByID(bookID, false)
	if err != nil {
		s.log.Warnf("Ошибка получения книги перед подтверждением: %v", err)
		return err
	}

	book.Confirmed = true

	err = s.bookRepository.UpdateBook(book, nil)
	if err != nil {
		s.log.Warnf("Ошибка подтверждения книги: %v", err)
		return err
	}
	return nil
}

// GetBookByID получает книгу по ID с авторами
func (s *BookService) GetBookByID(bookID uuid.UUID) (*models.Book, error) {
	book, err := s.bookRepository.GetBookByID(bookID, false)
	if err != nil {
		s.log.Warnf("Ошибка получения книги по ID: %v", err)
		return nil, err
	}
	return book, nil
}

// UpdateBook обновляет данные книги и связи с авторами
func (s *BookService) UpdateBook(bookID uuid.UUID, title, description, coverImage string, authorIDs []uuid.UUID) error {
	book, err := s.bookRepository.GetBookByID(bookID, false)
	if err != nil {
		s.log.Warnf("Ошибка получения книги перед обновлением: %v", err)
		return err
	}

	book.Title = title
	book.Description = description
	book.CoverImage = coverImage

	err = s.bookRepository.UpdateBook(book, authorIDs)
	if err != nil {
		s.log.Warnf("Ошибка обновления книги: %v", err)
		return err
	}
	return nil
}

// DeleteBook удаляет книгу, обнуляет связи и удаляет отзывы
func (s *BookService) DeleteBook(bookID uuid.UUID) error {
	err := s.bookRepository.DeleteBook(bookID)
	if err != nil {
		s.log.Warnf("Ошибка удаления книги: %v", err)
		return err
	}
	return nil
}

// GetBooksPaginated получает список книг с маркерной пагинацией
func (s *BookService) GetBooksPaginated(limit int, afterID *uuid.UUID) (*dto.PaginatedBooksResponse, error) {
	books, err := s.bookRepository.GetBooksPaginated(limit, afterID)
	if err != nil {
		s.log.Warnf("Ошибка получения списка книг: %v", err)
		return nil, err
	}

	if len(books) == 0 {
		return &dto.PaginatedBooksResponse{Books: []dto.BookResponse{}, NextCursor: nil}, nil
	}

	bookIDs := make([]uuid.UUID, len(books))
	for i, book := range books {
		bookIDs[i] = book.ID
	}

	bookAuthors, err := s.bookAuthorMappingRepository.GetAuthorsForBooks(bookIDs)
	if err != nil {
		s.log.Warnf("Ошибка получения авторов для книг: %v", err)
		return nil, err
	}

	authorIDMap := make(map[uuid.UUID]struct{})
	for _, ba := range bookAuthors {
		authorIDMap[*ba.AuthorID] = struct{}{}
	}

	authorIDs := make([]uuid.UUID, 0, len(authorIDMap))
	for authorID := range authorIDMap {
		authorIDs = append(authorIDs, authorID)
	}

	authors, err := s.authorRepository.GetAuthorsByIDs(authorIDs)
	if err != nil {
		s.log.Warnf("Ошибка загрузки авторов: %v", err)
		return nil, err
	}

	authorMap := make(map[uuid.UUID]dto.AuthorByBookResponse, len(authors))
	for _, author := range authors {
		authorMap[author.ID] = dto.AuthorByBookResponse{
			ID:   author.ID,
			Name: author.Name,
		}
	}

	bookAuthorMap := make(map[uuid.UUID][]dto.AuthorByBookResponse, len(books))
	for _, ba := range bookAuthors {
		bookAuthorMap[*ba.BookID] = append(bookAuthorMap[*ba.BookID], authorMap[*ba.AuthorID])
	}

	bookResponses := make([]dto.BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = dto.BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			CoverImage:  book.CoverImage,
			Authors:     bookAuthorMap[book.ID], // Авторы привязываются из мапы
		}
	}

	var nextAfterID *uuid.UUID
	if len(books) > 0 {
		nextAfterID = &books[len(books)-1].ID
	}

	return &dto.PaginatedBooksResponse{
		Books:      bookResponses,
		NextCursor: nextAfterID,
	}, nil

}

// UpdateBookCover обновляет путь к обложке книги
func (s *BookService) UpdateBookCover(bookID string, coverPath string) error {
	return s.bookRepository.UpdateBookCover(bookID, coverPath)
}

func (s *BookService) GetConfirmedBookByIdWithAuthors(bookID uuid.UUID) (*dto.BookResponse, error) {
	book, err := s.bookRepository.GetBookByID(bookID, true)
	if err != nil {
		s.log.Warnf("Ошибка получения книги по id=%s: %v", bookID, err)
		return nil, err
	}

	bookAuthors, err := s.bookAuthorMappingRepository.GetAuthorIDsByBookID(bookID)

	if err != nil {
		s.log.Warnf("Ошибка получения авторов для книги: %v", err)
		return nil, err
	}

	authors, err := s.authorRepository.GetAuthorsByIDs(bookAuthors)
	if err != nil {
		s.log.Warnf("Ошибка загрузки авторов: %v", err)
		return nil, err
	}

	authorResponses := make([]dto.AuthorByBookResponse, len(authors))
	for i, author := range authors {
		authorResponses[i] = dto.AuthorByBookResponse{
			ID:   author.ID,
			Name: author.Name,
		}
	}

	bookResponse := &dto.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		CoverImage:  book.CoverImage,
		Authors:     authorResponses,
	}

	return bookResponse, nil

}
