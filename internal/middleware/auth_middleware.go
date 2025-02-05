package middleware

import (
	jwtutil "book-management-system/pkg/jwtutil"
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware проверяет JWT-токен и извлекает информацию о пользователе
func AuthMiddleware() gin.HandlerFunc {
	log := logger.GetLogger()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			log.Warn("Отсутствует токен авторизации")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			log.Warn("Некорректный формат заголовка авторизации")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Некорректный формат заголовка префикса токена"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := jwtutil.ParseAndValidateToken(tokenString)
		if err != nil {
			log.Warnf("Ошибка валидации токена: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
