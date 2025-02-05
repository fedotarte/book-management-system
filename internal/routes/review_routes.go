package routes

import (
	"book-management-system/internal/handlers"
	"book-management-system/internal/middleware"
	"book-management-system/internal/repositories"
	"book-management-system/internal/services"
	"github.com/gin-gonic/gin"
)

// RegisterReviewRoutes регистрирует роуты отзывов
func RegisterReviewRoutes(
	r *gin.RouterGroup,
	reviewRepo *repositories.ReviewRepository,
	bookRepo *repositories.BookRepository,
) {
	reviewService := services.NewReviewService(reviewRepo, bookRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	reviewRoutes := r.Group("/reviews")
	{
		reviewRoutes.POST("/", middleware.AuthMiddleware(), reviewHandler.CreateReview)
		reviewRoutes.GET("/:reviewID", reviewHandler.GetReviewByID)
		reviewRoutes.PUT("/:reviewID", middleware.AuthMiddleware(), reviewHandler.UpdateReview)
		reviewRoutes.DELETE("/:reviewID", middleware.AuthMiddleware(), reviewHandler.DeleteReview)
		reviewRoutes.POST("/:reviewID/vote", middleware.AuthMiddleware(), reviewHandler.VoteReview)
	}
}
