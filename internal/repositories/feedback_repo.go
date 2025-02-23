package repositories

import (
	"book-management-system/internal/database"
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FeedbackRepository struct {
	collection string
	log        *logger.Logger
}

func NewFeedbackRepository() *FeedbackRepository {
	return &FeedbackRepository{
		collection: "feedbacks",
		log:        logger.GetLogger(),
	}
}

func (r *FeedbackRepository) CreateFeedback(feedbackDto dto.BaseFeedbackRequest, reviewer uuid.UUID) (primitive.ObjectID, error) {
	feedback := models.Feedback{
		Checked:   false,
		Text:      feedbackDto.Text,
		Rating:    feedbackDto.Rating,
		UserID:    &reviewer,
		CreatedAt: time.Now().UTC(),
	}
	insertedResult, err := database.MongoDB.Database("bookstore").
		Collection(r.collection).
		InsertOne(context.TODO(), feedback)

	if err != nil {
		r.log.Warnf("Ошибка добавления отзыва о приложении: %v", err)
		return primitive.NilObjectID, err
	}

	insertedID, ok := insertedResult.InsertedID.(primitive.ObjectID)

	if !ok {
		r.log.Warnf("Ошибка приведения InsertedID к ObjectID")
		return primitive.NilObjectID, fmt.Errorf("не удалось привести InsertedID к ObjectID")
	}

	return insertedID, nil
}

func (r *FeedbackRepository) FindPaginatedFeedbacks(onlyUnchecked bool, limit int, afterID *primitive.ObjectID) ([]models.Feedback, *primitive.ObjectID, error) {
	filter := bson.M{}

	if onlyUnchecked {
		filter["checked"] = false
	}

	if afterID != nil {
		var lastFeedback models.Feedback
		err := database.MongoDB.Database("bookstore").
			Collection(r.collection).
			FindOne(context.TODO(), bson.M{"_id": *afterID}).
			Decode(&lastFeedback)
		if err != nil {
			r.log.Warnf("Ошибка получения курсора отзыва: %v", err)
			return nil, nil, err
		}

		filter["created_at"] = bson.M{"$gt": lastFeedback.CreatedAt}
	}

	options := options.Find()
	options.SetSort(bson.D{{"created_at", 1}}) // Сортировка по возрастанию
	options.SetLimit(int64(limit))             // Ограничение по количеству

	cursor, err := database.MongoDB.Database("bookstore").
		Collection(r.collection).
		Find(context.TODO(), filter, options)
	if err != nil {
		r.log.Warnf("Ошибка получения отзывов: %v", err)
		return nil, nil, err
	}
	defer cursor.Close(context.TODO())

	var feedbacks []models.Feedback
	if err = cursor.All(context.TODO(), &feedbacks); err != nil {
		r.log.Warnf("Ошибка обработки отзывов: %v", err)
		return nil, nil, err
	}

	var nextCursor *primitive.ObjectID
	if len(feedbacks) > 0 {
		lastID := feedbacks[len(feedbacks)-1].ID // Последний `ObjectID`
		nextCursor = &lastID
	}

	return feedbacks, nextCursor, nil
}

func (r *FeedbackRepository) CheckReview(feedbackID primitive.ObjectID) error {
	collection := database.MongoDB.Database("bookstore").Collection(r.collection)
	var existingFeedback models.Feedback
	err := collection.FindOne(context.TODO(), bson.M{"_id": feedbackID}).Decode(&existingFeedback)
	if err != nil {
		r.log.Warnf("Ошибка получения отзыва о приложеении перед обновлением: %v", err)
		return err
	}

	updatedReview := bson.M{
		"$set": bson.M{
			"checked":    true,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": feedbackID}, updatedReview)
	if err != nil {
		r.log.Warnf("Ошибка обновления отзыва: %v", err)
		return err
	}
	return nil
}
