package usecase

import (
	"context"

	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"github.com/alligatorO15/taskmind-backend/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationUseCase — бизнес-логика управления уведомлениями
type NotificationUseCase struct {
	notificationRepo repository.NotificationRepository
}

// NewNotificationUseCase — создает экземпляр NotificationUseCase
func NewNotificationUseCase(notificationRepo repository.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo: notificationRepo,
	}
}

// Create — создает новое уведомление
func (uc *NotificationUseCase) Create(ctx context.Context, notification *models.Notification) error {
	return uc.notificationRepo.Create(ctx, notification)
}

// GetByUserID — возвращает уведомления пользователя с пагинацией
func (uc *NotificationUseCase) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*models.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.notificationRepo.FindByUserID(ctx, userID, limit, offset)
}

// MarkAsRead — помечает уведомление как прочитанное
func (uc *NotificationUseCase) MarkAsRead(ctx context.Context, notificationID, userID primitive.ObjectID) error {
	return uc.notificationRepo.MarkAsRead(ctx, notificationID, userID)
}

// MarkAllAsRead — помечает все уведомления пользователя как прочитанные
func (uc *NotificationUseCase) MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error {
	return uc.notificationRepo.MarkAllAsRead(ctx, userID)
}

// CountUnread — возвращает количество непрочитанных уведомлений
func (uc *NotificationUseCase) CountUnread(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return uc.notificationRepo.CountUnread(ctx, userID)
}
