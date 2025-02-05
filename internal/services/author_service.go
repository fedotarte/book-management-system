package services

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/internal/repositories"
	"book-management-system/pkg/logger"
	"github.com/google/uuid"
)

type AuthorService struct {
	authorRepository            *repositories.AuthorRepository
	bookAuthorMappingRepository *repositories.BookAuthorRepository
	bookRepository              *repositories.BookRepository
	log                         *logger.Logger
}

func NewAuthorService(
	bookRepository *repositories.BookRepository,
	bookAuthorMappingRepository *repositories.BookAuthorRepository,
	authorRepository *repositories.AuthorRepository,
) *AuthorService {
	return &AuthorService{
		authorRepository:            authorRepository,
		bookRepository:              bookRepository,
		bookAuthorMappingRepository: bookAuthorMappingRepository,
		log:                         logger.GetLogger(),
	}
}

func (s *AuthorService) CreateAuthor(author dto.CreateAuthorRequest) (*models.Author, error) {
	authorToSave := &models.Author{
		//ID:   uuid.New(),
		Name: author.Name,
		Bio:  author.Bio,
	}

	err := s.authorRepository.CreateAuthor(authorToSave, author.BookIDS)
	if err != nil {
		s.log.Warnf("create author error: %v", err)
		return nil, err
	}

	return authorToSave, nil
}

func (s *AuthorService) GetAuthorByID(id uuid.UUID) (*dto.DetailedAuthorResponse, error) {
	author, err := s.authorRepository.GetAuthorByID(id)
	if err != nil {
		s.log.Warnf("get author error: %v", err)
		return nil, err
	}

	resAuthor := dto.DetailedAuthorResponse{
		Name: author.Name,
		Bio:  author.Bio,
		ID:   author.ID,
	}

	return &resAuthor, nil
}

func (s *AuthorService) UpdateAuthor(authorId uuid.UUID, author dto.UpdateAuthorRequest) error {
	authorByID, err := s.authorRepository.GetAuthorByID(authorId)
	if err != nil {
		s.log.Warnf("get author by id error: %v", err)
	}

	authorByID.Name = author.Name
	authorByID.Bio = author.Bio

	err = s.authorRepository.UpdateAuthor(authorByID, author.BookIDS)

	if err != nil {
		s.log.Warnf("update author error: %v", err)
		return err
	}

	return nil
}

func (s *AuthorService) DeleteAuthor(id uuid.UUID) error {
	err := s.authorRepository.DeleteAuthor(id)
	if err != nil {
		s.log.Warnf("Ошибка удаления книги: %v", err)
		return err
	}

	return nil
}

func (s *AuthorService) GetAuthorList(limit int, afterID *uuid.UUID) (*dto.PaginatedAuthorsResponse, error) {
	authors, err := s.authorRepository.GetAuthorsPaginated(limit, afterID)
	if err != nil {
		s.log.Warnf("Ошибка получения списка авторов %v", err)
		return nil, err
	}

	if len(authors) == 0 {
		return &dto.PaginatedAuthorsResponse{Authors: []dto.AuthorWithBookResponse{}, NextCursor: nil}, nil
	}

	// 1️⃣ Collect author IDs
	authorIds := make([]uuid.UUID, len(authors))
	for i, author := range authors {
		authorIds[i] = author.ID
	}

	// 2️⃣ Get Book-Author Mappings
	authorBooks, err := s.bookAuthorMappingRepository.GetBooksForAuthors(authorIds)
	if err != nil {
		s.log.Warnf("Ошибка получения книг по авторам: %v", err)
		return nil, err
	}

	// 3️⃣ Extract unique book IDs
	bookIDMap := make(map[uuid.UUID]struct{})
	for _, ab := range authorBooks {
		if ab.BookID != nil { // Ensure BookID is not nil
			bookIDMap[*ab.BookID] = struct{}{}
		}
	}

	bookIds := make([]uuid.UUID, 0, len(bookIDMap))
	for bookID := range bookIDMap {
		bookIds = append(bookIds, bookID)
	}

	// 4️⃣ Fetch Books by IDs
	books, err := s.bookRepository.GetBooksByIds(bookIds)
	if err != nil {
		s.log.Warnf("Ошибка загрузки книг: %v", err)
		return nil, err
	}

	// Формируем карту книг {bookID -> BookByAuthorResponse}
	booksMap := make(map[uuid.UUID]dto.BookByAuthorResponse, len(books))
	for _, book := range books {
		booksMap[book.ID] = dto.BookByAuthorResponse{
			ID:    book.ID,
			Title: book.Title,
		}
	}

	// Формируем карту {authorID -> []Books}
	authorBookMap := make(map[uuid.UUID][]dto.BookByAuthorResponse, len(authors))
	for _, ab := range authorBooks {
		if book, exists := booksMap[*ab.BookID]; exists {
			authorBookMap[*ab.AuthorID] = append(authorBookMap[*ab.AuthorID], book)
		}
	}

	// Собираем финальный список авторов с их книгами
	authorsResponse := make([]dto.AuthorWithBookResponse, len(authors))
	for i, author := range authors {
		authorsResponse[i] = dto.AuthorWithBookResponse{
			ID:    author.ID,
			Name:  author.Name,
			Books: authorBookMap[author.ID], // Связываем книги с авторами
		}
	}

	// 8️⃣ Define Next Cursor
	var nextAfterID *uuid.UUID
	if len(authors) > 0 {
		nextAfterID = &authors[len(authors)-1].ID
	}

	return &dto.PaginatedAuthorsResponse{
		Authors:    authorsResponse,
		NextCursor: nextAfterID,
	}, nil
}
