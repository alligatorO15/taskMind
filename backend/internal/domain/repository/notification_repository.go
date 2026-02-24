package repository

import (
	"context"

	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationRepository — интерфейс для работы с хранилищем уведомлений
type NotificationRepository interface {
	// Create — сохраняет новое уведомление
	Create(ctx context.Context, notification *models.Notification) error
	// FindByUserID — возвращает уведомления пользователя, отсортированные по времени
	FindByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*models.Notification, error)
	// MarkAsRead — помечает уведомление как прочитанное
	MarkAsRead(ctx context.Context, id, userID primitive.ObjectID) error
	// MarkAllAsRead — помечает все уведомления пользователя как прочитанные
	MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error
	// CountUnread — возвращает количество непрочитанных уведомлений
	CountUnread(ctx context.Context, userID primitive.ObjectID) (int64, error)
}
