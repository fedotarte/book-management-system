package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

// User model
type User struct {
	ID        uuid.UUID  `gorm:"primaryKey"`
	Username  string     `gorm:"uniqueIndex;not null"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Password  string     `gorm:"not null"`
	Role      string     `gorm:"not null;default:user"`
	DeletedAt *time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
