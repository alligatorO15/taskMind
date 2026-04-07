package repository

import (
	"context"
	"time"

	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
)

// ReminderPublisher — интерфейс для отправки отложенных напоминаний в очередь.
type ReminderPublisher interface {
	PublishReminder(ctx context.Context, msg models.ReminderMessage, delay time.Duration) error
}
