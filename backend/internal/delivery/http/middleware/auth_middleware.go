package middleware

import (
	"net/http"
	"strings"

	"github.com/alligatorO15/taskMind/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware - middleware для проверки JWT access-токена
// Извлекает токен из заголвока Authorzation (формат: "Bearer <token>")
// валилирует и помкщает userID в контекст запроса
func AuthMiddleware(authUC *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "отсутствует заголовок авторизации"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		userID, err := authUC.ParseAccessToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидный или просроченный токен"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()

	}
}
