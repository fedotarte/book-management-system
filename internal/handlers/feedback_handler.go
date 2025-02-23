package handlers

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/services"
	"book-management-system/pkg/logger"
	"book-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FeedbackHandler struct {
	service *services.FeedbackService
	log     *logger.Logger
}

// NewFeedbackHandler создает новый экземпляр FeedbackHandler
func NewFeedbackHandler(service *services.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		service: service,
		log:     logger.NewLogger(),
	}
}

// CreateFeedback создаёт новый отзыв о приложении
//
//	@Summary		Создать отзыв о приложении
//	@Description	Добавляет новый отзыв к книге
//	@Tags			Feedbacks
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			review	body		dto.BaseFeedbackRequest	true	"Данные для создания отзыва"
//	@Success		201		{object}	dto.CreatedFeedbackResponse
//	@Failure		400		{object}	map[string]string	"Неверный формат запроса"
//	@Failure		500		{object}	map[string]string	"Ошибка сервера"
//	@Router			/feedbacks [post]
func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	var feedback dto.BaseFeedbackRequest
	if err := c.ShouldBindJSON(&feedback); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	userID, userIdExists := c.Get("userID")
	if !userIdExists {
		h.log.Warnf("Отсутствует создатель фидбэка")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует создатель фидбэка"})
	}
	stringUserID, err := utils.ConvertStringToUUID(userID.(string))
	if err != nil {
		h.log.Warnf("Ошибка привязки создателя отзыва: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка привязки создателя отзыва"})
		return
	}
	createdFeedbackId, err := h.service.CreateFeedback(feedback, stringUserID)
	if err != nil {
		h.log.Warnf("Ошибка создания отзыва: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать отзыв"})
		return
	}

	createdFeedbackIdResponse := &dto.CreatedFeedbackResponse{
		CreatedFeedbackId: createdFeedbackId.String(),
	}

	c.JSON(http.StatusCreated, createdFeedbackIdResponse)

}

// GetFeedbacks создаёт новый отзыв о приложении
//
//	@Summary		получить отзывы о приложении
//	@Description	Добавляет новый отзыв к книге
//	@Tags			Feedbacks
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			after_id	query		string	false	"UUID последней автора (для пагинации)"
//	@Param			limit		query		int		false	"Количество авторов на страницу (по умолчанию 10)"
//	@Param			checked		query		bool	false	"только проверенные"
//	@Success		200			{object}	dto.PaginatedFeedbackResponse
//	@Failure		400			{object}	map[string]string	"Неверный формат запроса"
//	@Failure		500			{object}	map[string]string	"Ошибка сервера"
//	@Router			/feedbacks [get]
func (h *FeedbackHandler) GetFeedbacks(c *gin.Context) {

	role, _ := c.Get("role")

	stringifiedRole := role.(string)

	ableToCreateAuthor := stringifiedRole == "moderator" || stringifiedRole == "admin"

	if !ableToCreateAuthor {
		h.log.Warnf("у пользователя нет необходимой роли")
		c.JSON(http.StatusForbidden, gin.H{"error": "Не хватает роли для работы с отзывом"})
		return
	}

	checkedParam := c.Query("checked")
	limitParam := c.Query("limit")
	afterIDParam := c.Query("afterID")

	checked := checkedParam == "true"
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		h.log.Warn("Ошибка конвертации limit")
		limit = 10
	}

	afterID, err := utils.ConvertStringToObjectID(afterIDParam)
	if err != nil {
		h.log.Warnf("Ошибка конвертации afterID: %v", err)
	}

	paginatedFeedbacks, err := h.service.FindPaginatedFeedbacks(checked, limit, &afterID)
	if err != nil {
		h.log.Warn(err)
	}

	c.JSON(http.StatusOK, paginatedFeedbacks)

}

// CheckFeedback обновляет отзыв о приложении
//
//	@Summary		получить отзывы о приложении
//	@Description	Добавляет новый отзыв к книге
//	@Tags			Feedbacks
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			feedbackID	path		string	false	"ID отзыва"
//	@Success		200			{object}	dto.PaginatedFeedbackResponse
//	@Failure		400			{object}	map[string]string	"Неверный формат запроса"
//	@Failure		500			{object}	map[string]string	"Ошибка сервера"
//	@Router			/feedbacks/{feedbackID} [put]
func (h *FeedbackHandler) CheckFeedback(c *gin.Context) {
	role, _ := c.Get("role")

	stringifiedRole := role.(string)

	ableToCreateAuthor := stringifiedRole == "moderator" || stringifiedRole == "admin"

	if !ableToCreateAuthor {
		h.log.Warnf("у пользователя нет необходимой роли")
		c.JSON(http.StatusForbidden, gin.H{"error": "Не хватает роли для работы с отзывом"})
		return
	}

	feedbackIdParam := c.Param("feedbackID")

	feedbackId, convErr := utils.ConvertStringToObjectID(feedbackIdParam)
	if convErr != nil {
		h.log.Warnf("Ошибка конвертации feedbackId: %v", convErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	err := h.service.CheckFeedback(feedbackId)
	if err != nil {
		h.log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обработки отзыва"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
