package services

import (
	"book-management-system/internal/dto"
	"book-management-system/internal/repositories"
	"book-management-system/pkg/logger"
	"book-management-system/pkg/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackService struct {
	feedbackRepo *repositories.FeedbackRepository
	log          *logger.Logger
}

func NewFeedbackService(feedbackRepo *repositories.FeedbackRepository) *FeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
		log:          logger.GetLogger(),
	}
}

func (s *FeedbackService) CreateFeedback(feedback dto.BaseFeedbackRequest, creator uuid.UUID) (primitive.ObjectID, error) {
	createdId, err := s.feedbackRepo.CreateFeedback(feedback, creator)
	if err != nil {
		s.log.Warnf("Failed to create feedback %v", err)
		return primitive.NilObjectID, err
	}

	return createdId, err
}

func (s *FeedbackService) FindPaginatedFeedbacks(onlyUnchecked bool, limit int, afterID *primitive.ObjectID) (*dto.PaginatedFeedbackResponse, error) {
	afterIdObject, err := utils.ConvertStringToObjectID(afterID.String())
	if err != nil {
		return nil, err
	}

	// Fetch data from repository
	feedbacks, lastObjectID, err := s.feedbackRepo.FindPaginatedFeedbacks(onlyUnchecked, limit, &afterIdObject)
	if err != nil {
		s.log.Warnf("Failed to find feedbacks: %v", err)
		return nil, err
	}

	feedbackResponses := make([]dto.FeedbackResponse, len(feedbacks))

	for i, feedback := range feedbacks {
		feedbackResponses[i] = dto.FeedbackResponse{
			ID:        feedback.ID,
			UserID:    feedback.UserID,
			Text:      feedback.Text,
			Rating:    feedback.Rating,
			Checked:   feedback.Checked,
			CreatedAt: feedback.CreatedAt,
			UpdatedAt: feedback.UpdatedAt,
		}
	}

	var lastIDString *string
	if lastObjectID != nil {
		idString := lastObjectID.Hex()
		lastIDString = &idString
	}

	return &dto.PaginatedFeedbackResponse{
		Feedbacks: feedbackResponses,
		LastID:    lastIDString,
	}, nil

}

func (s *FeedbackService) CheckFeedback(feedbackID primitive.ObjectID) error {
	return s.feedbackRepo.CheckReview(feedbackID)
}
