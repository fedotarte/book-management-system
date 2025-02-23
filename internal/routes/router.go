package routes

import (
	"book-management-system/internal/middleware"
	"book-management-system/internal/repositories"
	"github.com/gin-gonic/gin"
)

// InitRouter собирает все роуты проекта
func InitRouter() *gin.Engine {
	userRepo := repositories.NewUserRepository()
	refreshTokenRepo := repositories.NewRefreshTokenRepository()
	bookRepo := repositories.NewBookRepository()
	booksAuthorMappingRepo := repositories.NewBookAuthorRepository()
	userBookRepo := repositories.NewUserBookRepository()
	authorRepo := repositories.NewAuthorRepository()
	reviewRepo := repositories.NewReviewRepository()
	feedbackRepo := repositories.NewFeedbackRepository()

	r := gin.Default()

	// мидлвари
	r.Use(middleware.LoggerMiddleware())
	// Сваггер
	RegisterSwaggerRoutes(r)
	// хэлс-чек
	RegisterHealthCheckRoutes(r)

	apiV1 := r.Group("/api/v1")

	RegisterUserRoutes(apiV1, userRepo, refreshTokenRepo)
	RegisterBookRoutes(apiV1, bookRepo, booksAuthorMappingRepo, authorRepo)
	RegisterUserBookRoutes(apiV1, userBookRepo)
	RegisterAuthorRoutes(apiV1, bookRepo, booksAuthorMappingRepo, authorRepo)
	RegisterReviewRoutes(apiV1, reviewRepo, bookRepo)
	RegisterFeedbackRoutes(apiV1, feedbackRepo)

	return r
}
