package routes

import (
	"book-management-system/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterHealthCheckRoutes(r *gin.Engine) {
	healthCheckHandler := handlers.NewHealthCheckHandler()
	r.GET("/health", healthCheckHandler.HealthCheck)
}
