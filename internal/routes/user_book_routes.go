package routes

import (
	"book-management-system/internal/handlers"
	"book-management-system/internal/middleware"
	"book-management-system/internal/services"
	"github.com/gin-gonic/gin"
)

// RegisterUserBookRoutes регистрирует роуты для управления книгами пользователя
func RegisterUserBookRoutes(r *gin.RouterGroup, userBookService *services.UserBookService) {
	userBookHandler := handlers.NewUserBookHandler(userBookService)

	userBookRoutes := r.Group("/users/me/books")
	userBookRoutes.Use(middleware.AuthMiddleware())
	{
		userBookRoutes.POST("/", userBookHandler.AddBookToUser)
		userBookRoutes.GET("/", userBookHandler.GetUserBooks)
		userBookRoutes.PUT("/:bookID/progress", userBookHandler.UpdateReadingProgress)
		userBookRoutes.DELETE("/:bookID", userBookHandler.RemoveBookFromUser)
	}
}
