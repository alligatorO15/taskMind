package handler

import (
	"net/http"
	"strconv"

	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/usecase"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationHandler — HTTP-хендлер для эндпоинтов управления уведомлениями.
type NotificationHandler struct {
	notificationUC *usecase.NotificationUseCase
}

// NewNotificationHandler — создает экземпляр NotificationHandler.
func NewNotificationHandler(notificationUC *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{notificationUC: notificationUC}
}

// GetAll — возвращает список уведомлений пользователя с пагинацией.
// Поддерживает query-параметры: limit (по умолчанию 20), offset (по умолчанию 0).
func (h *NotificationHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	notifications, err := h.notificationUC.GetByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}

	if notifications == nil {
		notifications = []*models.Notification{}
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkAsRead — помечает уведомление как прочитанное.
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	notificationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор уведомления"})
		return
	}

	err = h.notificationUC.MarkAsRead(c.Request.Context(), notificationID, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "уведомление помечено как прочитанное"})
}

// MarkAllAsRead — помечает все уведомления пользователя как прочитанные.
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	err := h.notificationUC.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "все уведомления помечены как прочитанные"})
}

// CountUnread — возвращает количество непрочитанных уведомлений.
func (h *NotificationHandler) CountUnread(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	count, err := h.notificationUC.CountUnread(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
