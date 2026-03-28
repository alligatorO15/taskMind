package handler

import (
	"fmt"
	"net/http"

	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthHandler - HTTP-хендлер для эндпоинтов аутентификации
type AuthHandler struct {
	authUC *usecase.AuthUseCase
}

// NewAuthHandler -  создает экземпляр AuthHandler
func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// Register - регистрирует нового пользователя
// Принимает UserRegisterRequest, вызывает authUC.Register, возвращает 201 с парой токенов
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("некорректные данные запроса: %s", err.Error())})
		return
	}

	tokenPair, err := h.authUC.Register(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, tokenPair)
}

// Login - авторищует пользователя по email и паролю
// Принимает UserLoginRequest, вызывает authUc.Login, возвращает 200 с парой токенов
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("некорректные данные запроса:", err.Error())})
		return
	}

	tokenPair, err := h.authUC.Login(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// RefreshToken - обновляет access-токен по refresh-токену
// Принимает RefreshTokenRequest, вызывает authUC.RefreshToken, возвращает 200 с парой токенов
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("некорректные данные запроса: %s", err.Error())})
		return
	}

	tokenPair, err := h.authUC.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// GetProfile - возвращает профиль текущего пользователя
// Извлекает userID из контекста (установлен AuthMiddleware), вызывает authUC.GetUserByID, возвращает 200
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	userID, ok := userIDVal.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения идентификатора пользователя"})
		return
	}

	user, err := h.authUC.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
