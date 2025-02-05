package middleware

import (
	"book-management-system/internal/constants"
	"book-management-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RoleMiddleware проверяет, имеет ли пользователь нужную роль
func RoleMiddleware(allowedRoles ...constants.Role) gin.HandlerFunc {
	log := logger.GetLogger()
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			log.Warn("Роль пользователя не найдена в контексте")
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
			c.Abort()
			return
		}

		userRole, ok := role.(constants.Role)
		if !ok {
			log.Warn("Ошибка преобразования роли пользователя")
			c.JSON(http.StatusForbidden, gin.H{"error": "Ошибка авторизации"})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		log.Warnf("Доступ запрещен для роли %s", userRole)
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		c.Abort()
	}
}
