package handlers

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/services"
	"book-management-system/pkg/logger"
	"book-management-system/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"strconv"
)

type BookHandler struct {
	service *services.BookService
	log     *logger.Logger
}

// NewBookHandler создает новый обработчик книг
func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{
		service: service,
		log:     logger.GetLogger(),
	}
}

// CreateBook создаёт новую книгу
//
//	@Summary		Создать книгу
//	@Description	Добавляет новую книгу в базу
//	@Tags			Books
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			book	body		dto.CreateBookRequest	true	"Данные для создания книги"
//	@Success		201		{object}	dto.BookResponse
//	@Failure		400		{object}	map[string]string	"Invalid data"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req dto.CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Получаем роль пользователя из миддлвари
	role, _ := c.Get("role")

	// Создаем книгу через сервис
	book, err := h.service.CreateBook(req.Title, req.Description, req.CoverImage, req.AuthorIDs, role.(string))
	if err != nil {
		h.log.Warnf("Ошибка создания книги: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании книги"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetBookByID получает книгу по ID
//
//	@Summary		получить книгу по ID
//	@Description	получает книгу из базы по id
//	@Tags			Books
//	@Produce		json
//	@Param			bookID		path		string	true	"UUID книги"
//	@Param			withAuthors	query		bool	false	"получить привязку авторов"
//	@Param			confirmed	query		bool	false	"получить привязку авторов"
//	@Success		200			{object}	dto.BookResponse
//	@Failure		400			{object}	map[string]string	"Invalid data"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/books/{bookID} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("bookID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга bookID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	queryParams := map[string]bool{}
	keys := []string{"withAuthors", "confirmed"}

	for _, key := range keys {
		val, err := strconv.ParseBool(c.DefaultQuery(key, "false"))
		if err != nil {
			h.log.Warnf("Ошибка парсинга %s: %v", key, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неверный параметр %s", key)})
			return
		}
		queryParams[key] = val
	}

	withAuthors := queryParams["withAuthors"]
	confirmed := queryParams["confirmed"]

	var book interface{}
	if withAuthors && confirmed {
		book, err = h.service.GetConfirmedBookByIdWithAuthors(bookID)
	} else {
		book, err = h.service.GetBookByID(bookID)
	}

	if err != nil {
		h.log.Warnf("Ошибка получения книги: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Книга не найдена"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// UpdateBook обновляет информацию о книге
// UpdateBookByID получает книгу по ID
//
//	@Summary		обновить книгу по ID
//	@Description	обновить книгу в базе по id
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Param			bookID	path		string					true	"UUID книги"
//	@Param			book	body		dto.UpdateBookRequest	true	"Данные для обновления книги"
//	@Success		202		{object}	dto.BookResponse
//	@Failure		400		{object}	map[string]string	"Invalid data"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/books/{bookID} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("bookID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга bookID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	var req dto.UpdateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if err := h.service.UpdateBook(bookID, req.Title, req.Description, req.CoverImage, req.AuthorIDs); err != nil {
		h.log.Warnf("Ошибка обновления книги: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении книги"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Книга успешно обновлена"})
}

// DeleteBook удаляет книгу
//
//	@Summary		удалить книгу по ID
//	@Description	удалить книгу в базе по id
//	@Tags			Books
//	@Produce		json
//	@Param			bookID	path		string	true	"UUID книги"
//	@Success		204		{object}	dto.BookDeletionResponse
//	@Failure		400		{object}	map[string]string	"Invalid data"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/books/{bookID} [delete]

func (h *BookHandler) DeleteBook(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("bookID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга bookID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	if err := h.service.DeleteBook(bookID); err != nil {
		h.log.Warnf("Ошибка удаления книги: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении книги"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Книга успешно удалена"})
}

// GetBooksPaginated возвращает список книг с маркерной пагинацией
//
//	@Summary		маркерная пагинация книг
//	@Description	Получить поджинированный список с курсором
//	@Tags			Books
//	@Produce		json
//	@Param			after_id	query		string	false	"UUID последней книги (для пагинации)"
//	@Param			limit		query		int		false	"Количество книг на страницу (по умолчанию 10)"
//	@Success		200			{array}		dto.PaginatedBooksResponse
//	@Failure		400			{object}	map[string]string	"Invalid data"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/books [get]
func (h *BookHandler) GetBooksPaginated(c *gin.Context) {
	queryAfterId := c.Query("after_id")
	queryLimit := c.Query("limit")
	limitInt, err := strconv.Atoi(queryLimit)
	if err != nil {
		h.log.Warnf("ошибка конвертации query limi=%s : %v", queryLimit, err)
		limitInt = 10
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

	books, err := h.service.GetBooksPaginated(limitInt, afterUUID)
	if err != nil {
		h.log.Warnf("Ошибка получения списка книг: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка книг"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// ConfirmBook подтверждает книгу (только для модераторов и админов)
//
//	@Summary		Подтвердить книгу
//	@Description	Модератор или админ может подтвердить книгу перед публикацией
//	@Tags			Books
//	@Security		BearerAuth
//	@Param			bookID	path		string				true	"UUID книги"
//	@Success		200		{object}	map[string]string	"Book confirmed"
//	@Failure		400		{object}	map[string]string	"Invalid data"
//	@Failure		404		{object}	map[string]string	"Book not found"
//	@Router			/books/{bookID}/confirm [put]
func (h *BookHandler) ConfirmBook(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("bookID"))
	if err != nil {
		h.log.Warnf("Ошибка парсинга bookID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор книги"})
		return
	}

	if err := h.service.ConfirmBook(bookID); err != nil {
		h.log.Warnf("Ошибка подтверждения книги: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подтверждении книги"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга подтверждена"})
}

// UploadBookCover загружает изображение обложки книги
//
//	@Summary		Загрузить обложку книги
//	@Description	Загрузка файла с обложкой книги (JPG, PNG)
//	@Tags			Books
//	@Security		BearerAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			bookID	path		string	true	"UUID книги"
//	@Param			cover	formData	file	true	"Файл изображения (JPG, PNG, макс. 50 KB)"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string	"Invalid file format"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/books/{bookID}/upload [post]
func (h *BookHandler) UploadBookCover(c *gin.Context) {
	bookID := c.Param("id")
	if bookID == "" {
		log.Warn("Отсутствует ID книги")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID книги обязателен"})
		return
	}

	file, err := c.FormFile("cover")
	if err != nil {
		log.Warnf("Ошибка загрузки файла: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл обязателен"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		log.Warnf("Недопустимый формат файла: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Допустимые форматы: JPG, PNG"})
		return
	}

	src, err := file.Open()
	if err != nil {
		log.Warnf("Ошибка открытия файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки файла"})
		return
	}
	defer src.Close()

	const maxSize = 50 * 1024
	const safeSize = 5 * 1024

	fileSize := file.Size
	if fileSize > maxSize {
		log.Warnf("Файл слишком большой: %d KB", fileSize/1024)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл слишком большой, максимум 50 KB"})
		return
	}

	filePath := fmt.Sprintf("uploads/%s%s", bookID, ext)

	// Если файл > 5 KB, сжимаем его перед сохранением
	if fileSize > safeSize {
		img, err := utils.DecodeImage(src, ext)
		if err != nil {
			log.Warnf("Ошибка декодирования изображения: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки изображения"})
			return
		}

		err = utils.CompressAndSaveImage(img, filePath, 75)
		if err != nil {
			log.Warnf("Ошибка сжатия файла: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки изображения"})
			return
		}
	} else {
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			log.Warnf("Ошибка сохранения файла: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки файла"})
			return
		}
	}

	err = h.service.UpdateBookCover(bookID, "/"+filePath)
	if err != nil {
		log.Warnf("Ошибка обновления обложки книги: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cover_url": "/" + filePath})
}
