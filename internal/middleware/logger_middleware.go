package middleware

import (
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

// LoggerMiddleware записывает информацию о запросах в лог
func LoggerMiddleware() gin.HandlerFunc {
	log := logger.GetLogger()
	return func(c *gin.Context) {
		start := time.Now()

		// Выполняем запрос
		c.Next()

		// Логируем информацию о запросе
		log.Infof("%s %s %d %s", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}
