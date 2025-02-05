package dto

// BaseReviewRequest содержит общие поля для создания и обновления отзыва
type BaseReviewRequest struct {
	// ID книги, к которой относится отзыв (UUID)
	// Required: true
	// Example: "123e4567-e89b-12d3-a456-426614174000"
	BookID string `json:"book_id" binding:"required"`

	// Текст отзыва
	// Required: true
	// Example: "Отличная книга, советую!"
	Text string `json:"text" binding:"required"`

	// Оценка (от 1 до 10)
	// Required: true
	// Example: 9
	Rating int `json:"rating" binding:"required,min=1,max=10"`
}

// ReviewResponse DTO для ответа API с информацией об отзыве
// @Description Ответ API с информацией об отзыве
type ReviewResponse struct {

	// Уникальный идентификатор отзыва (ObjectID в MongoDB)
	// Example: "60c72b2f5f1b2c001f6f1b20"
	ID string `json:"id"`

	BaseReviewRequest

	// ID автора отзыва (UUID)
	// Example: "550e8400-e29b-41d4-a716-446655440000"
	UserID string `json:"user_id"`

	// Лайки
	// Example: 10
	Likes int `json:"likes"`

	// Дизлайки
	// Example: 2
	Dislikes int `json:"dislikes"`

	// Дата создания
	// Example: "2024-02-01T12:00:00Z"
	CreatedAt string `json:"created_at"`

	// Дата обновления (если есть)
	// Example: "2024-02-02T14:30:00Z"
	UpdatedAt *string `json:"updated_at,omitempty"`
}

// ReviewListResponse DTO для списка отзывов
// @Description Ответ API со списком отзывов
type ReviewListResponse struct {
	// Массив отзывов
	Reviews []ReviewResponse `json:"reviews"`
}

// VoteReviewRequest DTO для голосования ща отзыв
// @Description Запрос API с голосованием (1 -1 0)
type VoteReviewRequest struct {
	Vote int `json:"vote"` // 1 - лайк, -1 - дизлайк, 0 - удалить голос
}
