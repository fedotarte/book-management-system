package services

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/internal/repositories"
	"book-management-system/pkg/logger"
	"book-management-system/pkg/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService struct {
	reviewRepo *repositories.ReviewRepository
	bookRepo   *repositories.BookRepository
	log        *logger.Logger
}

// NewReviewService создает новый сервис
func NewReviewService(reviewRepo *repositories.ReviewRepository, bookRepo *repositories.BookRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		bookRepo:   bookRepo,
		log:        logger.GetLogger(),
	}
}

// CreateReview добавляет новый отзыв и пересчитывает рейтинг
func (s *ReviewService) CreateReview(review dto.BaseReviewRequest, creator string) error {
	err := s.reviewRepo.CreateReview(review, creator)
	if err != nil {
		s.log.Warnf("Ошибка создания отзыва: %v", err)
		return err
	}

	// Пересчитываем рейтинг книги
	bookIdAsUUID, err := utils.ConvertStringToUUID(review.BookID)
	if err != nil {
		s.log.Warnf("Ошибка конвертации строки bookID %s s в UUID: %v", review.BookID, err)
	}
	return s.RecalculateBookRating(bookIdAsUUID)
}

func (s *ReviewService) GetReviewById(reviewID primitive.ObjectID) (*models.Review, error) {
	review, err := s.reviewRepo.GetReviewById(reviewID)
	if err != nil {
		s.log.Warnf("Ошибка получения review по ID %s : %v", reviewID.Hex(), err)

	}

	return review, nil
}

func (s *ReviewService) VoteReview(reviewID primitive.ObjectID, userID string, vote int) error {
	return s.reviewRepo.VoteReview(reviewID, userID, vote)
}

// UpdateReview обновляет отзыв, если изменился рейтинг — пересчитываем средний
func (s *ReviewService) UpdateReview(reviewID primitive.ObjectID, updatedText string, updatedRating int, editor string) error {
	existingReview, err := s.reviewRepo.GetReviewById(reviewID)
	if err != nil {
		s.log.Warnf("Ошибка получения отзыва: %v", err)
		return err
	}

	// Проверяем, изменился ли рейтинг
	shouldRecalculate := existingReview.Rating != updatedRating

	err = s.reviewRepo.UpdateReviewText(reviewID, updatedText, editor)
	if err != nil {
		s.log.Warnf("Ошибка обновления отзыва: %v", err)
		return err
	}

	// Если рейтинг изменился, пересчитываем
	if shouldRecalculate {
		bookIdAsUUID, err := utils.ConvertStringToUUID(existingReview.BookID)
		if err != nil {
			s.log.Warnf("Ошибка конвертации строки bookID %s s в UUID: %v", existingReview.BookID, err)
		}

		return s.RecalculateBookRating(bookIdAsUUID)
	}

	return nil
}

// DeleteReviewByID удаляет отзыв и пересчитывает рейтинг
func (s *ReviewService) DeleteReviewByID(reviewID primitive.ObjectID) error {
	review, err := s.reviewRepo.GetReviewById(reviewID)
	if err != nil {
		s.log.Warnf("Ошибка получения отзыва перед удалением: %v", err)
		return err
	}

	err = s.reviewRepo.DeleteReviewByID(reviewID)
	if err != nil {
		s.log.Warnf("Ошибка удаления отзыва: %v", err)
		return err
	}

	bookIdAsUUID, err := utils.ConvertStringToUUID(review.BookID)
	if err != nil {
		s.log.Warnf("Ошибка конвертации строки bookID %s s в UUID: %v", review.BookID, err)
	}

	// Пересчитываем рейтинг книги
	return s.RecalculateBookRating(bookIdAsUUID)
}

// RecalculateBookRating пересчитывает средний рейтинг книги
func (s *ReviewService) RecalculateBookRating(bookID uuid.UUID) error {
	bookIdStringified := utils.ConvertUUIDToString(bookID)

	averageRating, err := s.reviewRepo.CalculateAverageRating(bookIdStringified)
	if err != nil {
		s.log.Warnf("Ошибка агрегации рейтинга: %v", err)
		return err
	}

	// Обновляем средний рейтинг в PostgreSQL
	err = s.bookRepo.UpdateBookRating(bookID, averageRating)
	if err != nil {
		s.log.Warnf("Ошибка обновления среднего рейтинга: %v", err)
		return err
	}

	return nil
}
