package services

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/internal/repositories"
	"book-management-system/pkg/jwtutil"
	"book-management-system/pkg/logger"
	"book-management-system/pkg/utils"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository

	log *logger.Logger
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo *repositories.UserRepository, refreshTokenRepo *repositories.RefreshTokenRepository) *UserService {
	return &UserService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,

		log: logger.GetLogger(),
	}
}

// RegisterUser регистрирует нового пользователя
func (s *UserService) RegisterUser(req dto.UserRegisterRequest) error {
	// Проверяем, существует ли уже пользователь
	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return errors.New("пользователь уже зарегистрирован")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Warnf("Ошибка хеширования пароля: %v", err)
		return err
	}

	// Создаем пользователя
	user := models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.RoleUser, // По умолчанию обычный пользователь
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		s.log.Warnf("Ошибка сохранения пользователя: %v", err)
		return err
	}

	return nil
}

// LoginUser проверяет учетные данные и выдает JWT
func (s *UserService) LoginUser(req dto.UserLoginRequest) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", "", errors.New("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", errors.New("неверный email или пароль")
	}

	accessToken, err := jwtutil.GenerateToken(utils.ConvertUUIDToString(user.ID), user.Role, time.Minute*60)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwtutil.GenerateToken(utils.ConvertUUIDToString(user.ID), user.Role, time.Hour*24*7) // 7 дней
	if err != nil {
		return "", "", err
	}

	// Сохраняем refresh-токен в БД
	err = s.refreshTokenRepo.SaveToken(user.ID, refreshToken, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken обновляет `access_token` и выдаёт новый `refresh_token`
func (s *UserService) RefreshToken(req dto.TokenRefreshRequest) (string, string, error) {
	// Проверяем, есть ли `refresh_token` в БД
	existingToken, err := s.refreshTokenRepo.GetToken(req.Token)
	if err != nil {
		return "", "", errors.New("недействительный `refresh_token`")
	}

	// Проверяем, не истёк ли токен
	if time.Now().After(existingToken.ExpiresAt) {
		err := s.refreshTokenRepo.DeleteToken(req.Token)
		if err != nil {
			return "", "", err
		}

		return "", "", errors.New("`refresh_token` истёк")
	}

	// Проверяем валидность `refresh_token`
	claims, err := jwtutil.ParseAndValidateToken(req.Token)
	if err != nil {
		s.log.Warnf("ошибка валидации токена: %v", err)

		return "", "", errors.New("ошибка валидации токена")
	}

	// Генерируем новый `access_token`
	newAccessToken, err := jwtutil.GenerateToken(claims.UserID, claims.Role, time.Minute*15)
	if err != nil {
		return "", "", err
	}

	// Генерируем новый `refresh_token`
	newRefreshToken, err := jwtutil.GenerateToken(claims.UserID, claims.Role, time.Hour*24*7) // 7 дней
	if err != nil {
		return "", "", err
	}

	// Удаляем старый `refresh_token`
	err = s.refreshTokenRepo.DeleteToken(req.Token)

	if err != nil {
		return "", "", err
	}

	stringUserID, err := utils.ConvertStringToUUID(claims.UserID)
	if err != nil {
		return "", "", err
	}
	// Сохраняем новый `refresh_token`

	err = s.refreshTokenRepo.SaveToken(stringUserID, newRefreshToken, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *UserService) GetUserByID(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {

		return nil, err
	}
	userResponse := &dto.UserResponse{
		ID:    utils.ConvertUUIDToString(user.ID),
		Role:  user.Role,
		Email: user.Email,
	}

	return userResponse, nil
}
