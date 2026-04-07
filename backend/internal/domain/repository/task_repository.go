package repository

import (
	"context"
	"time"

	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository - интерфейс для работы с харнилищем задач
// поддерживает CRUD-операции, фильтрацию и поиск просроченных задач
type TaskRepository interface {
	// Create - создает новую задачу
	Create(ctx context.Context, task *models.Task) error
	// FindByID - находит задачу по id задачи и userID (для проверки прав доступа)
	FindByID(ctx context.Context, id, userID primitive.ObjectID) (*models.Task, error)
	// FindByUserID - возвращает все задачи пользователя с применением фильтров
	FindByUserID(ctx context.Context, userID primitive.ObjectID, filter models.TaskFilter) ([]*models.Task, error)
	// Update - обновляет существующую задачу
	Update(ctx context.Context, task *models.Task) error
	// Delete - удаляет задачу по id и userID
	Delete(ctx context.Context, is, userID primitive.ObjectID) error
	// FindOverdue - находит задачи с просроченным дедлайном, еще не отмеченные как overdue
	FindOverdue(ctx context.Context, now time.Time) ([]*models.Task, error)
	// FindPendingReminders - находит задачи, по которым нужно отправить напоминание
	FindPendingReminders(ctx context.Context, time time.Time) ([]*models.Task, error)
	// MarkReminderSent - помечает, что напоминание для задачи было отправлено
	MarkReminderSent(ctx context.Context, id primitive.ObjectID) error
	// UpdateStatus - обновляет статус задачи
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.TaskStatus) error
}
