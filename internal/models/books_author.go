package models

import (
	"github.com/google/uuid"
)

// BookAuthor — стыковочная таблица для связи книг и авторов (многие ко многим)
type BookAuthor struct {
	BookID   *uuid.UUID `gorm:"type:uuid;index;onDelete:SET NULL" json:"book_id"`
	AuthorID *uuid.UUID `gorm:"type:uuid;index;onDelete:SET NULL" json:"author_id"`
}
