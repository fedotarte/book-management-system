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
	authorRepo := repositories.NewAuthorRepository()
	reviewRepo := repositories.NewReviewRepository()

	r := gin.Default()

	// мидлвари
	r.Use(middleware.LoggerMiddleware())
	//Сваггер
	RegisterSwaggerRoutes(r)
	// хэлс-чек
	RegisterHealthCheckRoutes(r)

	apiV1 := r.Group("/api/v1")

	RegisterUserRoutes(apiV1, userRepo, refreshTokenRepo)
	RegisterBookRoutes(apiV1, bookRepo, booksAuthorMappingRepo, authorRepo)
	RegisterAuthorRoutes(apiV1, bookRepo, booksAuthorMappingRepo, authorRepo)
	RegisterReviewRoutes(apiV1, reviewRepo, bookRepo)

	return r
}
