package dto

import "github.com/google/uuid"

// CreateBookRequest тело запроса на создание книги
type CreateBookRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CoverImage  string      `json:"cover_image"`
	AuthorIDs   []uuid.UUID `json:"author_ids"`
}

// UpdateBookRequest тело запроса на обновление книги
// @Description запрос API на обновление книги
type UpdateBookRequest struct {
	CreateBookRequest
}

// BookResponse тело ответа на создание/обновление книги
// @Description Ответ API на создание/обновление/получение книги
type BookResponse struct {
	// Уникальный идентификатор книги (UUID)
	// Example: "123e4567-e89b-12d3-a456-426614174000"
	ID uuid.UUID `json:"id"`

	// Название книги (обязательное поле)
	// Required: true
	// Example: "Гарри Поттер и философский камень"
	Title string `json:"title"`

	// Описание книги
	// Example: "Первая книга о приключениях Гарри Поттера"
	Description string `json:"description"`

	// Обложка книги (URL)
	// Example: "/uploads/123e4567-e89b-12d3-a456-426614174000.jpg"
	CoverImage string `json:"cover_image"`

	// Средний рейтинг книги (из 10)
	// Example: 8.5
	AverageRating float64 `json:"average_rating"`

	// Авторы книги (массив объектов)
	Authors []AuthorByBookResponse `json:"authors"`
}

// BookListResponse DTO для списка книг с пагинацией
// @Description Ответ API со списком книг
type PaginatedBooksResponse struct {
	// Массив книг
	Books []BookResponse `json:"books"`

	// Следующий маркер для пагинации (если есть)
	// Example: "123e4567-e89b-12d3-a456-426614174002"
	NextCursor *uuid.UUID `json:"next_cursor,omitempty"`
}

// BookDeletionResponse DTO для списка книг с пагинацией
// @Description Ответ API со списком книг
type BookDeletionResponse struct {
	// Массив книг
	Message string `json:"message"`
}

type BookByAuthorResponse struct {
	// Уникальный идентификатор книги (UUID)
	// Example: "123e4567-e89b-12d3-a456-426614174000"
	ID uuid.UUID `json:"id"`

	// Название книги (обязательное поле)
	// Required: true
	// Example: "Гарри Поттер и философский камень"
	Title string `json:"title"`
}
