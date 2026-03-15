package mongo

import (
	"context"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/domain/repository"
	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const notificationsCollection = "notifications"

// NotificationRepository — MongoDB-реализация интерфейса repository.NotificationRepository
type NotificationRepository struct {
	coll *mongo.Collection
}

// NewNotificationRepository — конструктор репозитория уведомлений
func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		coll: db.Collection(notificationsCollection),
	}
}

// Create — сохраняет новое уведомление
func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	now := time.Now()
	notification.CreatedAt = now

	if notification.ID.IsZero() {
		notification.ID = primitive.NewObjectID()
	}

	_, err := r.coll.InsertOne(ctx, notification)
	return err
}

// FindByUserID — возвращает уведомления пользователя, отсортированные по created_at по убыванию, с limit и offset
func (r *NotificationRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*models.Notification, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

// MarkAsRead — помечает уведомление как прочитанное по идентификатору и идентификатору пользователя
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id, "user_id": userID},
		bson.D{{Key: "$set", Value: bson.M{"read": true}}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// MarkAllAsRead — помечает все уведомления пользователя как прочитанные
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.coll.UpdateMany(ctx,
		bson.M{"user_id": userID, "read": false},
		bson.D{{Key: "$set", Value: bson.M{"read": true}}},
	)
	return err
}

// CountUnread — возвращает количество непрочитанных уведомлений пользователя
func (r *NotificationRepository) CountUnread(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{"user_id": userID, "read": false})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Проверка, что NotificationRepository реализует интерфейс repository.NotificationRepository
var _ repository.NotificationRepository = (*NotificationRepository)(nil)
