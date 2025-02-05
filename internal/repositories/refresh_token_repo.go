package repositories

import (
	"book-management-system/internal/database"
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewRefreshTokenRepository создаёт новый репозиторий
func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db:  database.DB,
		log: logger.GetLogger(),
	}
}

// SaveToken сохраняет refresh-токен
func (r *RefreshTokenRepository) SaveToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	refreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	err := r.db.Create(&refreshToken).Error
	if err != nil {
		r.log.Warnf("Ошибка сохранения refresh-токена: %v", err)
		return err
	}
	return nil
}

// GetToken получает refresh-токен из БД
func (r *RefreshTokenRepository) GetToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		r.log.Warnf("Ошибка поиска refresh-токена: %v", err)
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteToken удаляет refresh-токен при logout
func (r *RefreshTokenRepository) DeleteToken(token string) error {
	err := r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
	if err != nil {
		r.log.Warnf("Ошибка удаления refresh-токена: %v", err)
		return err
	}
	return nil
}

// CleanupExpiredTokens удаляет просроченные refresh-токены (CRON)
func (r *RefreshTokenRepository) CleanupExpiredTokens() {
	r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{})
	r.log.Info("Удалены устаревшие refresh-токены")
}

// Запуск задачи удаления токенов раз в день
func StartTokenCleanupTask(repo *RefreshTokenRepository) {
	go func() {
		for {
			repo.CleanupExpiredTokens()
			time.Sleep(24 * time.Hour) // Удаляем раз в день
		}
	}()
}
