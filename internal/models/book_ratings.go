package models

import (
	"github.com/google/uuid"
	"time"
)

type BookRating struct {
	UserID    uuid.UUID `gorm:"type:uuid;index;primaryKey" json:"user_id"`
	BookID    uuid.UUID `gorm:"type:uuid;index;primaryKey" json:"book_id"`
	Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 10" json:"rating"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
