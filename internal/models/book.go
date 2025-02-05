package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Book model
type Book struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	CoverImage  string         `json:"cover_image"`
	Confirmed   bool           `gorm:"default:false" json:"confirmed"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}
