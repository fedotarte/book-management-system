package dto

import (
	"github.com/google/uuid"
	"time"
)

// AddBookRequest DTO для добавления книги в список пользователя
type AddBookRequest struct {
	BookID uuid.UUID `json:"book_id"`
}

type UpdateReadingProgressRequest struct {
	Status    string `json:"status"`
	PagesRead int    `json:"pages_read"`
}

type UserBookResponse struct {
	BookID    *uuid.UUID `json:"book_id"`
	Status    string     `json:"status"`
	PagesRead int        `json:"pages_read"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
