package repositories

import (
	"book-management-system/internal/database"
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"errors"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db:  database.DB,
		log: logger.GetLogger(),
	}
}

// CreateUser сохраняет нового пользователя в базе данных
func (r *UserRepository) CreateUser(user models.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		r.log.Warnf("Ошибка создания пользователя: %v", err)
		return err
	}
	return nil
}

// GetUserByEmail ищет пользователя по email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("пользователь не найден")
	} else if err != nil {
		r.log.Warnf("Ошибка поиска пользователя: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
