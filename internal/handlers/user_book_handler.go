package handlers

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/internal/services"
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var log = logger.GetLogger()

type UserBookHandler struct {
	service *services.UserBookService
}

// NewUserBookHandler создает новый обработчик
func NewUserBookHandler(service *services.UserBookService) *UserBookHandler {
	return &UserBookHandler{service: service}
}

// AddBookToUser добавляет книгу в список пользователя
//
//	@Summary		Добавить книгу в список пользователя
//	@Description	Добавляет указанную книгу в список пользователя
//	@Tags			UserBooks
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			book	body		dto.AddBookRequest	true	"ID книги"
//	@Success		201		{object}	map[string]string	"message: Книга добавлена в список"
//	@Failure		400		{object}	map[string]string	"Неверный формат запроса"
//	@Failure		500		{object}	map[string]string	"Ошибка сервера"
//	@Router			/users/me/books/ [post]
func (h *UserBookHandler) AddBookToUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.AddBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	err := h.service.AddBookToUser(userID.(uuid.UUID), req.BookID)
	if err != nil {
		log.Warnf("Ошибка добавления книги пользователю: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении книги"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Книга добавлена в список"})
}

// UpdateReadingProgress обновляет статус и прогресс чтения книги
//
//	@Summary		Обновить прогресс чтения
//	@Description	Обновляет статус и количество прочитанных страниц
//	@Tags			UserBooks
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			bookID		path		string								true	"UUID книги"
//	@Param			progress	body		dto.UpdateReadingProgressRequest	true	"Прогресс чтения"
//	@Success		200			{object}	map[string]string					"message: Прогресс чтения обновлен"
//	@Failure		400			{object}	map[string]string					"Неверный формат запроса"
//	@Failure		500			{object}	map[string]string					"Ошибка сервера"
//	@Router			/users/me/books/{bookID}/progress [put]
func (h *UserBookHandler) UpdateReadingProgress(c *gin.Context) {
	userID, _ := c.Get("userID")
	bookID, err := uuid.Parse(c.Param("bookID"))
	if err != nil {
		log.Warnf("Ошибка парсинга bookID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	var req dto.UpdateReadingProgressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	err = h.service.UpdateReadingProgress(userID.(uuid.UUID), bookID, models.ReadingStatus(req.Status), req.PagesRead)
	if err != nil {
		log.Warnf("Ошибка обновления прогресса чтения: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении прогресса"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Прогресс чтения обновлен"})
}

// RemoveBookFromUser удаляет книгу из списка пользователя
//
//	@Summary		Удалить книгу из списка пользователя
//	@Description	Удаляет указанную книгу из списка пользователя
//	@Tags			UserBooks
//	@Security		BearerAuth
//	@Produce		json
//	@Param			bookID	path		string				true	"UUID книги"
//	@Success		200		{object}	map[string]string	"message: Книга удалена из списка"
//	@Failure		400		{object}	map[string]string	"Неверный формат запроса"
//	@Failure		500		{object}	map[string]string	"Ошибка сервера"
//	@Router			/users/me/books/{bookID} [delete]
func (h *UserBookHandler) RemoveBookFromUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	bookID, err := uuid.Parse(c.Param("bookID"))
	if err != nil {
		log.Warnf("Ошибка парсинга bookID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	err = h.service.RemoveBookFromUser(userID.(uuid.UUID), bookID)
	if err != nil {
		log.Warnf("Ошибка удаления книги из списка пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении книги"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга удалена из списка"})
}

// GetUserBooks получает список книг пользователя
//
//	@Summary		Получить список книг пользователя
//	@Description	Возвращает список книг, добавленных пользователем
//	@Tags			UserBooks
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		dto.UserBookResponse	"Список книг пользователя"
//	@Failure		500	{object}	map[string]string		"Ошибка сервера"
//	@Router			/users/me/books [get]
func (h *UserBookHandler) GetUserBooks(c *gin.Context) {
	userID, _ := c.Get("userID")

	books, err := h.service.GetUserBooks(userID.(uuid.UUID))
	if err != nil {
		log.Warnf("Ошибка получения списка книг пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка книг"})
		return
	}

	c.JSON(http.StatusOK, books)
}
