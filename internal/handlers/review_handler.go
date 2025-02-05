package handlers

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/internal/services"
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"slices"
)

type ReviewHandler struct {
	service *services.ReviewService
	log     *logger.Logger
}

// NewReviewHandler создает новый экземпляр ReviewHandler
func NewReviewHandler(service *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		service: service,
		log:     logger.GetLogger(),
	}
}

// GetReviewByID получает отзыв по ID
//
//	@Summary		Получить отзыв по ID
//	@Description	Возвращает карточку отзыва
//	@Tags			Reviews
//	@Produce		json
//	@Param			reviewID	path		string	true	"ObjectID отзыва"
//	@Success		200			{object}	dto.ReviewResponse
//	@Failure		400			{object}	map[string]string	"Неверный ID"
//	@Failure		404			{object}	map[string]string	"Отзыв не найден"
//	@Router			/review/{reviewID} [get]
func (h *ReviewHandler) GetReviewByID(c *gin.Context) {
	reviewID, err := primitive.ObjectIDFromHex(c.Param("reviewID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга reviewID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор отзыва"})
		return
	}

	review, err := h.service.GetReviewById(reviewID)
	if err != nil {
		h.log.Warnf("Ошибка получения отзыва: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Отзыв не найден"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// CreateReview создаёт новый отзыв
//
//	@Summary		Создать отзыв
//	@Description	Добавляет новый отзыв к книге
//	@Tags			Reviews
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			review	body		dto.BaseReviewRequest	true	"Данные для создания отзыва"
//	@Success		201		{object}	dto.ReviewResponse
//	@Failure		400		{object}	map[string]string	"Неверный формат запроса"
//	@Failure		500		{object}	map[string]string	"Ошибка сервера"
//	@Router			/review [post]

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var review dto.BaseReviewRequest
	if err := c.ShouldBindJSON(&review); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	userID, userIdExists := c.Get("userID")
	if !userIdExists {
		h.log.Warnf("Отсутствует создатель ревью")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует создатель реаью"})
	}
	stringUserID := userID.(string)

	if err := h.service.CreateReview(review, stringUserID); err != nil {
		h.log.Warnf("Ошибка создания отзыва: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать отзыв"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Отзыв успешно создан"})
}

// UpdateReview обновляет существующий отзыв
//
//	@Summary		Обновить отзыв
//	@Description	Редактирует текст и/или оценку отзыва
//	@Tags			Reviews
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			reviewID	path		string					true	"ObjectID отзыва"
//	@Param			review		body		dto.BaseReviewRequest	true	"Данные для обновления отзыва"
//	@Success		200			{object}	map[string]string		"message: Отзыв успешно обновлён"
//	@Failure		400			{object}	map[string]string		"Неверные параметры"
//	@Failure		403			{object}	map[string]string		"Нет прав на редактирование"
//	@Failure		404			{object}	map[string]string		"Отзыв не найден"
//	@Router			/reviews/{reviewID} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	reviewID, err := primitive.ObjectIDFromHex(c.Param("reviewID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга reviewID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор отзыва"})
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	stringifiedUserId := userID.(string)

	var body struct {
		Text   string `json:"text"`
		Rating int    `json:"rating"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	review, err := h.service.GetReviewById(reviewID)
	if err != nil {
		h.log.Warnf("Ошибка получения отзыва: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Отзыв не найден"})
		return
	}

	if review.UserID != userID && role != models.RoleModerator && role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав на редактирование"})
		return
	}

	if err := h.service.UpdateReview(reviewID, body.Text, body.Rating, stringifiedUserId); err != nil {
		h.log.Warnf("Ошибка обновления отзыва: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления отзыва"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отзыв успешно обновлен"})
}

// VoteReview обрабатывает запрос на добавление отзыва
//
//	@Summary		Проголосовать за отзыв
//	@Description	Проголосовать за отзыв +1 -1 0
//	@Tags			Reviews
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			reviewID	path		string				true	"ObjectID отзыва"
//	@Param			book		body	dto.VoteReviewRequest	true	"Данные для голосования"
//	@Success		202
//	@Router			/review/{reviewID}/vote [post]
func (h *ReviewHandler) VoteReview(c *gin.Context) {
	reviewID, err := primitive.ObjectIDFromHex(c.Param("reviewID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга reviewID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор отзыва"})
		return
	}

	userID, _ := c.Get("userID")

	stringifiedUserId := userID.(string)

	var body dto.VoteReviewRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	allowedVotes := []int{-1, 0, 1}

	if !slices.Contains(allowedVotes, body.Vote) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое значение голоса"})
		return
	}

	if err := h.service.VoteReview(reviewID, stringifiedUserId, body.Vote); err != nil {
		h.log.Warnf("Ошибка голосования за отзыв: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка голосования"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Голос успешно учтен"})
}

// DeleteReview удаляет отзыв
//
//	@Summary		Удалить отзыв
//	@Description	Удаляет отзыв по ID (только автор, модератор или админ)
//	@Tags			Reviews
//	@Security		BearerAuth
//	@Produce		json
//	@Param			reviewID	path		string				true	"ObjectID отзыва"
//	@Success		200			{object}	map[string]string	"message: Отзыв успешно удалён"
//	@Failure		400			{object}	map[string]string	"Неверные параметры"
//	@Failure		403			{object}	map[string]string	"Нет прав на удаление"
//	@Failure		404			{object}	map[string]string	"Отзыв не найден"
//	@Router			/reviews/{reviewID} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID, err := primitive.ObjectIDFromHex(c.Param("reviewID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга reviewID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор отзыва"})
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	review, err := h.service.GetReviewById(reviewID)
	if err != nil {
		h.log.Warnf("Ошибка получения отзыва: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Отзыв не найден"})
		return
	}

	if review.UserID != userID && role != models.RoleModerator && role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав на удаление"})
		return
	}

	if err := h.service.DeleteReviewByID(reviewID); err != nil {
		h.log.Warnf("Ошибка удаления отзыва: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления отзыва"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отзыв успешно удален"})
}
