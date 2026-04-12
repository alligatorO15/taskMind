package middleware

import (
	"net/http"
	"strings"

	"github.com/alligatorO15/taskmind-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware — middleware для проверки JWT access-токена.
// Извлекает токен из заголовка Authorization (формат: "Bearer <token>"),
// валидирует его и помещает userID в контекст запроса.
func AuthMiddleware(authUC *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// Извлекаем токен из заголовка Authorization
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// Fallback: query-параметр ?token= (для WebSocket, где нельзя задать заголовки)
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "отсутствует токен авторизации"})
			c.Abort()
			return
		}

		userID, err := authUC.ParseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидный или просроченный токен"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
