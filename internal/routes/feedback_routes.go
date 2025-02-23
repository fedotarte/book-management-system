package routes

import (
	"book-management-system/internal/handlers"
	"book-management-system/internal/middleware"
	"book-management-system/internal/repositories"
	"book-management-system/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterFeedbackRoutes(
	r *gin.RouterGroup,
	feedbackRepo *repositories.FeedbackRepository,
) {
	feedbackService := services.NewFeedbackService(feedbackRepo)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)

	feedbackRoutes := r.Group("/feedbacks")
	{
		feedbackRoutes.POST("/", middleware.AuthMiddleware(), feedbackHandler.CreateFeedback)
		feedbackRoutes.GET("/", middleware.AuthMiddleware(), feedbackHandler.GetFeedbacks)
		feedbackRoutes.PUT("/:feedbackID", middleware.AuthMiddleware(), feedbackHandler.CheckFeedback)
	}

}
