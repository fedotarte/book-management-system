package dto

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// BaseFeedbackRequest содержит общие поля для создания и обновления отзыва
type BaseFeedbackRequest struct {

	// Текст отзыва об аппе
	// Required: true
	// Example: "разрабы дауны ваше прилага говно!"
	Text string `json:"text" binding:"required"`

	// Оценка (от 1 до 10)
	// Required: true
	// Example: 1
	Rating int `json:"rating" binding:"required,min=1,max=10"`
}

// CreatedFeedbackResponse содержит id созданного оотзыва
type CreatedFeedbackResponse struct {
	CreatedFeedbackId string `json:"createdFeedbackId"`
}

// FeedbackResponse DTO for returning feedback data
type FeedbackResponse struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    *uuid.UUID         `json:"user_id,omitempty"`
	Text      string             `json:"text"`
	Rating    int                `json:"rating"`
	Checked   bool               `json:"checked"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

// PaginatedFeedbackResponse DTO for paginated feedback results
type PaginatedFeedbackResponse struct {
	Feedbacks []FeedbackResponse `json:"feedbacks"`
	LastID    *string            `json:"last_id,omitempty"`
}
