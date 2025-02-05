package handlers

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/services"
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	service *services.UserService
	log     *logger.Logger
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service, log: logger.GetLogger()}
}

// RegisterUser регистрирует нового пользователя
//
//	@Summary		Регистрация пользователя
//	@Description	Создает нового пользователя в системе
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.UserRegisterRequest	true	"Данные для регистрации"
//	@Success		201		{object}	map[string]string		"message: Пользователь зарегистрирован"
//	@Failure		400		{object}	map[string]string		"Неверный формат запроса"
//	@Failure		500		{object}	map[string]string		"Ошибка сервера"
//	@Router			/users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	err := h.service.RegisterUser(req)
	if err != nil {
		h.log.Warnf("Ошибка регистрации пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь зарегистрирован"})
}

// LoginUser выполняет вход пользователя
//
//	@Summary		Аутентификация пользователя
//	@Description	Выполняет вход и возвращает токены
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			login	body		dto.UserLoginRequest	true	"Данные для входа"
//	@Success		200		{object}	dto.AuthResponse		"Токены доступа"
//	@Failure		400		{object}	map[string]string		"Неверный формат запроса"
//	@Failure		401		{object}	map[string]string		"Ошибка авторизации"
//	@Router			/users/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(req)
	if err != nil {
		h.log.Warnf("Ошибка авторизации пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"token": token})
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken обновляет `access_token` и `refresh_token`
// RefreshToken обновляет `access_token` и `refresh_token`
//
//	@Summary		Обновление токенов
//	@Description	Обновляет `access_token` и `refresh_token` по действующему `refresh_token`
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			refresh	body		dto.TokenRefreshRequest	true	"Refresh токен"
//	@Success		200		{object}	dto.AuthResponse		"Новые токены"
//	@Failure		400		{object}	map[string]string		"Неверный формат запроса"
//	@Failure		401		{object}	map[string]string		"Ошибка авторизации"
//	@Router			/users/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.TokenRefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("Ошибка привязки JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(req)
	if err != nil {
		log.Warnf("Ошибка обновления токена: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка обновления токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetCurrentUser получает информацию о текущем пользователе
//
//	@Summary		Информация о текущем пользователе
//	@Description	Возвращает информацию о пользователе по его `userID`
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	dto.UserResponse	"Данные о пользователе"
//	@Failure		401	{object}	map[string]string	"Пользователь не аутентифицирован"
//	@Failure		404	{object}	map[string]string	"Пользователь не найден"
//	@Router			/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	h.log.Info("userId is: ", userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	stringUserID := userID.(string)
	UUIDuserId, uuidParseErr := uuid.Parse(stringUserID)
	if uuidParseErr != nil {
		log.Warnf("parsing uuid %s error: %v", stringUserID, uuidParseErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка идентификации пользователя"})
		return
	}
	user, err := h.service.GetUserByID(UUIDuserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}
