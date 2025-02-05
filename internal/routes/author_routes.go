package routes

import (
	"book-management-system/internal/constants"
	"book-management-system/internal/handlers"
	"book-management-system/internal/middleware"
	"book-management-system/internal/repositories"
	"book-management-system/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterAuthorRoutes(r *gin.RouterGroup,
	bookRepo *repositories.BookRepository,
	booksAuthorMappingRepo *repositories.BookAuthorRepository,
	authorRepo *repositories.AuthorRepository) {

	authorService := services.NewAuthorService(bookRepo, booksAuthorMappingRepo, authorRepo)
	authorHandler := handlers.NewAuthorHandler(authorService)
	authorRoutes := r.Group("/authors")
	{
		authorRoutes.GET("/", authorHandler.GetAuthorsPaginated)
		authorRoutes.POST("/", middleware.AuthMiddleware(), middleware.RoleMiddleware(constants.Roles.Moderator, constants.Roles.Admin), authorHandler.CreateAuthor)
		authorRoutes.GET("/:authorID", authorHandler.GetAuthor)
		authorRoutes.PUT("/:authorID", middleware.AuthMiddleware(), middleware.RoleMiddleware(constants.Roles.Moderator, constants.Roles.Admin), authorHandler.UpdateAuthorByID)
		authorRoutes.DELETE("/:authorID", middleware.AuthMiddleware(), middleware.RoleMiddleware(constants.Roles.Moderator, constants.Roles.Admin), authorHandler.DeleteAuthor)

	}
}
