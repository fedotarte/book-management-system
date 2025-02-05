package dto

import "github.com/google/uuid"

// AuthorByBookResponse DTO с данными об авторе
// @Description Данные об авторе книги
type AuthorByBookResponse struct {
	// Уникальный идентификатор автора (UUID)
	// Example: "b6d46cd4-e89b-12d3-a456-426614174111"
	ID uuid.UUID `json:"id"`

	// Имя автора (обязательное поле)
	// Required: true
	// Example: "Дж. К. Роулинг"
	Name string `json:"name"`
}

type CreateAuthorRequest struct {
	// Имя автора (обязательное поле)
	// Required: true
	// Example: "Дж. К. Роулинг"
	Name string `json:"name"`

	// Имя автора (обязательное поле)
	// Required: true
	// Example: "Из Англии"
	Bio string `json:"bio"`

	// Список книг
	// Required: true
	// Example: "Дж. К. Роулинг"
	BookIDS []uuid.UUID `json:"book_ids"`
}

type UpdateAuthorRequest struct {
	CreateAuthorRequest
}

type DetailedAuthorResponse struct {
	ID uuid.UUID `json:"id"`

	// Имя автора (обязательное поле)
	// Required: true
	// Example: "Дж. К. Роулинг"
	Name string `json:"name"`

	// Имя автора (обязательное поле)
	// Required: true
	// Example: "Дж. К. Роулинг"
	Bio string `json:"bio"`
}

type AuthorWithBookResponse struct {
	// Уникальный идентификатор автора (UUID)
	// Example: "b6d46cd4-e89b-12d3-a456-426614174111"
	ID uuid.UUID `json:"id"`

	// Имя автора (обязательное поле)
	// Required: true
	// Example: "Дж. К. Роулинг"
	Name string `json:"name"`

	Books []BookByAuthorResponse `json:"books"`
}

type PaginatedAuthorsResponse struct {
	// Список авторов
	Authors []AuthorWithBookResponse `json:"authors"`

	// Следующий маркер для пагинации (если есть)
	// Example: "123e4567-e89b-12d3-a456-426614174002"
	NextCursor *uuid.UUID `json:"next_cursor,omitempty"`
}
