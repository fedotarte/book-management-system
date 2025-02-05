package models

import (
	"github.com/google/uuid"
	"time"
)

type ReadingStatus string

const (
	StatusReading   ReadingStatus = "reading"
	StatusCompleted ReadingStatus = "completed"
	StatusDropped   ReadingStatus = "dropped"
)

type UserBook struct {
	UserID    uuid.UUID     `gorm:"type:uuid;index;primaryKey" json:"user_id"`
	BookID    *uuid.UUID    `gorm:"type:uuid;index;primaryKey" json:"book_id"`
	Status    ReadingStatus `gorm:"type:varchar(20);not null" json:"status"`
	PagesRead int           `json:"pages_read"`
	CreatedAt time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
