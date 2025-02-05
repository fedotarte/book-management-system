package repositories

import (
	"book-management-system/internal/database"
	"book-management-system/internal/dto"
	"book-management-system/internal/models"
	"book-management-system/pkg/logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ReviewRepository struct {
	collection      string
	votesCollection string
	log             *logger.Logger
}

// NewReviewRepository создает новый экземпляр ReviewRepository
func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		collection:      "reviews",
		votesCollection: "review_votes",
		log:             logger.GetLogger(),
	}
}

// CreateReview добавляет новый отзыв в MongoDB
func (r *ReviewRepository) CreateReview(reviewDto dto.BaseReviewRequest, creator string) error {
	review := models.Review{
		BookID:    reviewDto.BookID,
		UserID:    creator,
		Text:      reviewDto.Text,
		Rating:    reviewDto.Rating,
		Likes:     0,
		Dislikes:  0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err := database.MongoDB.Database("bookstore").Collection(r.collection).InsertOne(context.TODO(), review)
	if err != nil {
		r.log.Warnf("Ошибка добавления отзыва: %v", err)
		return err
	}
	return nil
}

// UpdateReviewText обновляет существующий отзыв и сохраняет его версию
func (r *ReviewRepository) UpdateReviewText(reviewID primitive.ObjectID, updatedText string, editor string) error {
	collection := database.MongoDB.Database("bookstore").Collection(r.collection)

	// Получаем текущий отзыв
	var existingReview models.Review
	err := collection.FindOne(context.TODO(), bson.M{"_id": reviewID}).Decode(&existingReview)
	if err != nil {
		r.log.Warnf("Ошибка получения отзыва перед обновлением: %v", err)
		return err
	}

	// Добавляем версию
	updatedReview := bson.M{
		"$set": bson.M{
			"text":       updatedText,
			"updated_at": primitive.NewDateTimeFromTime(existingReview.UpdatedAt),
		},
		"$push": bson.M{
			"versions": models.ReviewVersion{
				Text:     existingReview.Text,
				EditedAt: existingReview.UpdatedAt,
				EditedBy: editor,
			},
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": reviewID}, updatedReview)
	if err != nil {
		r.log.Warnf("Ошибка обновления отзыва: %v", err)
		return err
	}
	return nil
}

// GetReviewById получает review по ID
func (r *ReviewRepository) GetReviewById(reviewID primitive.ObjectID) (*models.Review, error) {
	var review models.Review
	err := database.MongoDB.Database("bookstore").Collection(r.collection).FindOne(context.TODO(), bson.M{"_id": reviewID}).Decode(&review)
	if err != nil {
		r.log.Warnf("Ошибка получения отзыва по ID=%s %v", reviewID.Hex(), err)
		return nil, err
	}
	return &review, nil
}

// GetReviewsByBookID получает отзывы по ID книги
func (r *ReviewRepository) GetReviewsByBookID(bookID primitive.ObjectID) ([]models.Review, error) {
	var reviews []models.Review
	cursor, err := database.MongoDB.Database("bookstore").Collection(r.collection).Find(context.TODO(), bson.M{"book_id": bookID})
	if err != nil {
		r.log.Warnf("Ошибка получения отзывов: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &reviews); err != nil {
		r.log.Warnf("Ошибка обработки отзывов: %v", err)
		return nil, err
	}
	return reviews, nil
}

// DeleteReviewByID удаляет отзыв
func (r *ReviewRepository) DeleteReviewByID(reviewID primitive.ObjectID) error {
	_, err := database.MongoDB.Database("bookstore").Collection(r.collection).DeleteOne(context.TODO(), bson.M{"_id": reviewID})
	if err != nil {
		r.log.Warnf("Ошибка удаления отзыва: %v", err)
		return err
	}
	return nil
}

// VoteReview добавляет/обновляет голос за отзыв
// VoteReview добавляет/обновляет голос пользователя за отзыв
func (r *ReviewRepository) VoteReview(reviewID primitive.ObjectID, userID string, vote int) error {
	collection := database.MongoDB.Database("bookstore").Collection(r.votesCollection)

	// Если голос 0, удаляем его
	if vote == 0 {
		_, err := collection.DeleteOne(context.TODO(), bson.M{"review_id": reviewID, "user_id": userID})
		if err != nil {
			r.log.Warnf("Ошибка удаления голоса: %v", err)
			return err
		}
	} else {
		// Обновляем или вставляем голос
		_, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"review_id": reviewID, "user_id": userID},
			bson.M{"$set": bson.M{"vote": vote}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			r.log.Warnf("Ошибка голосования за отзыв: %v", err)
			return err
		}
	}

	return r.recalculateVotes(reviewID)
}

// recalculateVotes пересчитывает лайки/дизлайки для отзыва
func (r *ReviewRepository) recalculateVotes(reviewID primitive.ObjectID) error {
	collection := database.MongoDB.Database("bookstore").Collection(r.votesCollection)
	aggCursor, err := collection.Aggregate(context.TODO(), bson.A{
		bson.M{"$match": bson.M{"review_id": reviewID}},
		bson.M{"$group": bson.M{
			"_id":      "$review_id",
			"likes":    bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$vote", 1}}, 1, 0}}},
			"dislikes": bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$vote", -1}}, 1, 0}}},
		}},
	})
	if err != nil {
		r.log.Warnf("Ошибка пересчета голосов: %v", err)
		return err
	}
	defer aggCursor.Close(context.TODO())

	var result struct {
		Likes    int `bson:"likes"`
		Dislikes int `bson:"dislikes"`
	}
	if aggCursor.Next(context.TODO()) {
		if err := aggCursor.Decode(&result); err != nil {
			r.log.Warnf("Ошибка декодирования пересчитанных голосов: %v", err)
			return err
		}
	}

	// Обновляем отзыв с новыми значениями
	_, err = database.MongoDB.Database("bookstore").Collection(r.collection).UpdateOne(
		context.TODO(),
		bson.M{"_id": reviewID},
		bson.M{"$set": bson.M{"likes": result.Likes, "dislikes": result.Dislikes}},
	)
	if err != nil {
		r.log.Warnf("Ошибка обновления счетчиков голосов: %v", err)
		return err
	}

	return nil
}

// CalculateAverageRating вычисляет средний рейтинг книги на основе отзывов
func (r *ReviewRepository) CalculateAverageRating(bookID string) (float64, error) {
	collection := database.MongoDB.Database("bookstore").Collection(r.collection)

	// Агрегируем средний рейтинг книги
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"book_id", bookID}}}},
		bson.D{{"$group", bson.D{
			{"_id", nil},
			{"average_rating", bson.D{{"$avg", "$rating"}}},
		}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		r.log.Warnf("Ошибка агрегации среднего рейтинга для книги %s: %v", bookID, err)
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		AverageRating float64 `bson:"average_rating"`
	}

	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			r.log.Warnf("Ошибка декодирования среднего рейтинга: %v", err)
			return 0, err
		}
	} else {
		// Если отзывов нет, устанавливаем рейтинг в 0
		result.AverageRating = 0
	}

	return result.AverageRating, nil
}
