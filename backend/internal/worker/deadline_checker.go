package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/delivery/websocket"
	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/domain/repository"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/logger"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/rabbitmq"
)

// DeadlineChecker - фоновый воркер для периодической проверки просроченных задач
// Сканирует базуу данных на наличие задач с прошедшим дедлайном ( и меняет статус на overdue потом отправляет уведомление пользователю)
type DeadlineChecker struct {
	rabbitConn       *rabbitmq.Connection
	taskRepo         repository.TaskRepository
	notificationRepo repository.NotificationRepository
	wsHub            *websocket.Hub
	checkInterval    time.Duration
}

// NewDeadlineChecker - создает воркер проверки дедлайнов
func NewDeadlineChecker(
	rabbitConn *rabbitmq.Connection,
	taskRepo repository.TaskRepository,
	notificationRepo repository.NotificationRepository,
	wsHub *websocket.Hub,
	checkInterval time.Duration,
) *DeadlineChecker {
	return &DeadlineChecker{
		rabbitConn:       rabbitConn,
		taskRepo:         taskRepo,
		notificationRepo: notificationRepo,
		wsHub:            wsHub,
		checkInterval:    checkInterval,
	}
}

// Start - запускает периодическую проверку просрчоенных задач
// Ввыполняет первую проверу сразу
func (dc *DeadlineChecker) Start(ctx context.Context) {
	logger.Logger.Infof("Воркер проверки дедлайнов запущен (интервал: %v)", dc.checkInterval)

	// первая проверка псразу при запуске
	dc.check(ctx)

	ticker := time.NewTicker(dc.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Logger.Info("Воркер проверки дедлайнов остановлен")
			return
		case <-ticker.C:
			dc.check(ctx)
		}
	}
}

// check - выполняет одну итерацию проверки:
// находит просроченные задачи, обновляет их статус и отправляет уведомления
func (dc *DeadlineChecker) check(ctx context.Context) {
	now := time.Now()

	// находим задачи с просроченным дедлайном
	overdueTasks, err := dc.taskRepo.FindOverdue(ctx, now)
	if err != nil {
		logger.Logger.Errorf("Ошибка поиска просроченных задач: %v", err)
		return
	}

	for _, task := range overdueTasks {
		// обновляем статус
		if err := dc.taskRepo.UpdateStatus(ctx, task.ID, models.TaskStatusOverdue); err != nil {
			logger.Logger.Errorf("Ошибка обновления статуса задачи %s: %v", task.ID.Hex(), err)
			continue
		}

		// создаем уведомление о просрочке
		notification := &models.Notification{
			UserID:    task.UserID,
			TaskID:    task.ID,
			Type:      models.NotificationTypeOverdue,
			Title:     "Задача просрочена",
			Message:   fmt.Sprintf("Задача \"%s\" просрочена!", task.Title),
			Read:      false,
			CreatedAt: time.Now(),
		}

		if err := dc.notificationRepo.Create(ctx, notification); err != nil {
			logger.Logger.Errorf("Ошибка создания уведомления: %v", err)
			continue
		}

		// отправляем через websocket
		dc.wsHub.SendNotification(task.UserID, notification)
		logger.Logger.Infof("Задача %s помечена как просроченная", task.ID.Hex())
	}

	// ищем задачи для которых нужно отправить напоминание
	pendingReminders, err := dc.taskRepo.FindPendingReminders(ctx, now)
	if err != nil {
		logger.Logger.Errorf("Ошибка поиска задача для напоминаний: %v", err)
		return
	}

	for _, task := range pendingReminders {
		// отправляем напоминание через rabbitmq (будет обработано воркером ReminderWorker)
		msg := models.ReminderMessage{
			TaskID: task.ID.Hex(),
			UserID: task.UserID.Hex(),
			Type:   string(models.NotificationTypeReminder),
		}

		if err := dc.rabbitConn.PublishReminder(ctx, msg, 0); err != nil {
			logger.Logger.Errorf("Ошибка отправки напоминания в очередь: %v", err)
		}
	}

	if len(overdueTasks) > 0 || len(pendingReminders) > 0 {
		logger.Logger.Infof("Проверка дедлайнов: %d просрочено, %d напоминаний", len(overdueTasks), len(pendingReminders))
	}
}
