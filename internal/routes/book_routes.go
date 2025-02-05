package routes

import (
	"book-management-system/internal/constants"
	"book-management-system/internal/handlers"
	"book-management-system/internal/middleware"
	"book-management-system/internal/repositories"
	"book-management-system/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(
	r *gin.RouterGroup,
	bookRepo *repositories.BookRepository,
	booksAuthorMappingRepo *repositories.BookAuthorRepository,
	authorRepo *repositories.AuthorRepository,
) {

	bookService := services.NewBookService(bookRepo, booksAuthorMappingRepo, authorRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	bookRoutes := r.Group("/books")
	{
		bookRoutes.GET("/", bookHandler.GetBooksPaginated)
		bookRoutes.GET("/:bookID", bookHandler.GetBookByID)
		bookRoutes.POST("/", middleware.AuthMiddleware(), bookHandler.CreateBook)
		bookRoutes.PUT("/:bookID", middleware.AuthMiddleware(), bookHandler.UpdateBook)
		bookRoutes.DELETE("/:bookID", middleware.AuthMiddleware(), middleware.RoleMiddleware(constants.Roles.Moderator, constants.Roles.Admin), bookHandler.DeleteBook)
		bookRoutes.POST("/:bookID/confirm", middleware.AuthMiddleware(), middleware.RoleMiddleware(constants.Roles.Moderator, constants.Roles.Admin), bookHandler.ConfirmBook)
	}
}
