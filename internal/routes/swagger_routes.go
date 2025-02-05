package routes

import (
	"github.com/gin-gonic/gin"
	swaggofiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerRoutes(r *gin.Engine) {
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggofiles.Handler))
}
