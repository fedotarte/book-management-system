package handlers

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/services"
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type AuthorHandler struct {
	service *services.AuthorService
	log     *logger.Logger
}

func NewAuthorHandler(service *services.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		service: service,
		log:     logger.GetLogger(),
	}
}

// CreateAuthor создаёт нового автора
//
//	@Summary		Создать автора
//	@Description	Добавляет нового автора
//	@Tags			Authors
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			author	body		dto.CreateAuthorRequest	true	"Данные для создания автора"
//	@Success		201		{object}	dto.BookResponse
//	@Failure		400		{object}	map[string]string	"Invalid data"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/authors [post]
func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var req dto.CreateAuthorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("Ошибка привязки JSON %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	role, _ := c.Get("role")

	stringifiedRole := role.(string)

	ableToCreateAuthor := stringifiedRole == "moderator" || stringifiedRole == "admin"

	if !ableToCreateAuthor {
		h.log.Warnf("у пользователя нет необходимой роли")
		c.JSON(http.StatusForbidden, gin.H{"error": "Не хватает роли для модерации контента"})
		return
	}

	author, err := h.service.CreateAuthor(req)

	if err != nil {
		h.log.Warnf("Ошибка создания автора в базе: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании автора"})
		return
	}

	c.JSON(http.StatusCreated, author)
}

// GetAuthor получение автора
//
//	@Summary		получить автора
//	@Description	получить автора по id
//	@Tags			Authors
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			authorID	path		string	true	"ID  автора"
//	@Success		200			{object}	dto.BookResponse
//	@Failure		400			{object}	map[string]string	"Invalid data"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/authors/{authorID} [get]
func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	authorID, err := uuid.Parse(c.Param("authorID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга authorID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор автора"})
		return
	}

	author, err := h.service.GetAuthorByID(authorID)
	if err != nil {
		h.log.Warnf("Ошибка получения Автора: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Автор не найдена"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// GetAuthorsPaginated возвращает список авторов с маркерной пагинацией
//
//	@Summary		маркерная пагинация авторов
//	@Description	Получить поджинированный список с курсором
//	@Tags			Authors
//	@Produce		json
//	@Param			after_id	query		string	false	"UUID последней автора (для пагинации)"
//	@Param			limit		query		int		false	"Количество авторов на страницу (по умолчанию 10)"
//	@Success		200			{array}		dto.PaginatedAuthorsResponse
//	@Failure		400			{object}	map[string]string	"Invalid data"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/books [get]
func (h *AuthorHandler) GetAuthorsPaginated(c *gin.Context) {
	queryAfterId := c.Query("after_id")
	queryLimit := c.Query("limit")
	limitInt, err := strconv.Atoi(queryLimit)

	if err != nil {
		h.log.Warnf("ошибка конвертации query limi=%s : %v", queryLimit, err)
	}

	var afterUUID *uuid.UUID

	if queryAfterId != "" {
		parsedID, err := uuid.Parse(queryAfterId)
		if err != nil {
			h.log.Warnf("Ошибка парсинга after_id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный параметр after_id"})
			return
		}
		afterUUID = &parsedID
	}

	authors, err := h.service.GetAuthorList(limitInt, afterUUID)

	if err != nil {
		h.log.Warnf("Ошибка получения списка авторов: %v")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка авторов"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

// UpdateAuthorByID UpdateAuthor обновляет автора
//
//	@Summary		обновление автора
//	@Description	обновление автора
//	@Tags			Authors
//	@Produce		json
//	@Param			authorID	path		string					true	"UUID автора"
//	@Param			author		body		dto.UpdateAuthorRequest	true	"Данные для создания автора"	false
//	@Success		200			{object}	dto.AuthorByBookResponse
//	@Failure		400			{object}	map[string]string	"Invalid data"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/authors/{authorID} [put]
func (h *AuthorHandler) UpdateAuthorByID(c *gin.Context) {
	authorID, err := uuid.Parse(c.Param("authorID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга authorID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный authorID"})
		return
	}

	var req dto.UpdateAuthorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if err := h.service.UpdateAuthor(authorID, req); err != nil {
		h.log.Warnf("Ошибка обновления автора %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении автора"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Автор успешно обновлен"})
}

// DeleteAuthor удаляет автора
//
//	@Summary		удалить автора по ID
//	@Description	удалить автора в базе по id
//	@Tags			Authors
//	@Produce		json
//	@Param			authorID	path		string	true	"UUID автора"
//	@Success		204			{object}	dto.BookDeletionResponse
//	@Failure		400			{object}	map[string]string	"Invalid data"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/books/{authorID} [delete]
func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	authorID, err := uuid.Parse(c.Param("authorID"))
	if err != nil {
		h.log.Warnf("Ошибка получения authorID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения authorID"})
		return
	}

	if err := h.service.DeleteAuthor(authorID); err != nil {
		h.log.Warnf("Ошибка удаления автора: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении автора"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "автор успешно удален"})
}
