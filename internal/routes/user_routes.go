package routes

import (
	"book-management-system/internal/handlers"
	"book-management-system/internal/middleware"
	"book-management-system/internal/repositories"
	"book-management-system/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(
	r *gin.RouterGroup,
	userRepo *repositories.UserRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
) {

	userService := services.NewUserService(userRepo, refreshTokenRepo)
	userHandler := handlers.NewUserHandler(userService)

	repositories.StartTokenCleanupTask(refreshTokenRepo)

	authRoutes := r.Group("/users")
	{
		authRoutes.POST("/register", userHandler.RegisterUser)
		authRoutes.POST("/login", userHandler.LoginUser)
		authRoutes.POST("/refresh", userHandler.RefreshToken)
		authRoutes.GET("/me", middleware.AuthMiddleware(), userHandler.GetCurrentUser)
	}
}
