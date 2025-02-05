package models

import "time"

// Reading Progress model
type ReadingProgress struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	BookID    uint   `gorm:"index"`
	Status    string `gorm:"not null"` // читаю, прочитал, перестал
	PagesRead int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
