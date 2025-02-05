package models

import "time"

// Moderator Action model
type ModeratorAction struct {
	ID          uint   `gorm:"primaryKey"`
	ModeratorID uint   `gorm:"index"`
	Action      string `gorm:"not null"`
	TargetID    uint   `gorm:"index"`
	TargetType  string `gorm:"not null"` // book, author, review
	CreatedAt   time.Time
}
