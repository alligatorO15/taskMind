package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alligatorO15/taskmind-backend/internal/delivery/websocket"
	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"github.com/alligatorO15/taskmind-backend/internal/domain/repository"
	"github.com/alligatorO15/taskmind-backend/internal/infrastructure/logger"
	"github.com/alligatorO15/taskmind-backend/internal/infrastructure/rabbitmq"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReminderWorker - воркер обработки напоминаний из очереди RabitMQ
// Читает отложенные сообщения, создаем уведомления в БД
// и отправляет их пользователю через websocket
type ReminderWorker struct {
	rabbitConn       *rabbitmq.Connection
	taskRepo         repository.TaskRepository
	notificationRepo repository.NotificationRepository
	wsHub            *websocket.Hub
}

// NewReminderWorker - создает воркер обработки напоминаний (pull-модель)
func NewReminderWorker(
	rabbitConn *rabbitmq.Connection,
	taskRepo repository.TaskRepository,
	notificationRepo repository.NotificationRepository,
	wsHub *websocket.Hub,
) *ReminderWorker {
	return &ReminderWorker{
		rabbitConn:       rabbitConn,
		taskRepo:         taskRepo,
		notificationRepo: notificationRepo,
		wsHub:            wsHub,
	}
}

// Start - запускает воркера для потребелния сообщений из очереди наопминаний
// Работает в бесконечном цикле, обрабатывая каждое входящее сообщение
func (w *ReminderWorker) Start(ctx context.Context) error {
	msgs, err := w.rabbitConn.Consume()
	if err != nil {
		return fmt.Errorf("не удалось запустить потребление сообщений: %w", err)
	}

	logger.Logger.Infof("Воркер напоминаний запущен")

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Logger.Infof("Воркер напоминаний остановлен")
				return
			case msg, ok := <-msgs:
				if !ok {
					logger.Logger.Warn("Канал сообщений RabbitMQ закрыт")
					return
				}

				var reminderMsg models.ReminderMessage
				if err := json.Unmarshal(msg.Body, &reminderMsg); err != nil {
					logger.Logger.Errorf("Ошибка десериализации сообщения: %v", err)
					msg.Nack(false, false)
					continue
				}

				if err := w.processReminder(ctx, reminderMsg); err != nil {
					logger.Logger.Errorf("Ошибки обработки напоминания: %v", err)
					msg.Nack(false, true) // requeue для повторной попытки
					continue
				}

				msg.Ack(false)
			}

		}
	}()

	return nil
}

// processReminder - обрабатывает одно сообщение-напоминание
// Проверяет актулаьность задачи, создает уведомление и отправвляет через websocket
func (w *ReminderWorker) processReminder(ctx context.Context, msg models.ReminderMessage) error {
	taskID, err := primitive.ObjectIDFromHex(msg.TaskID)
	if err != nil {
		return fmt.Errorf("невалидный ID задачи: %w", err)
	}

	userID, err := primitive.ObjectIDFromHex(msg.UserID)
	if err != nil {
		return fmt.Errorf("невалидный ID пользвоателя: %w", err)
	}

	// проверяем что заадча еще актулаьна (не удалена и не завершена)
	task, err := w.taskRepo.FindByID(ctx, taskID, userID)
	if err != nil {
		logger.Logger.Warnf("Задача %s не найдена, пропускаем напоминание", msg.TaskID)
		return nil
	}

	if task.Status == models.TaskStatusDone {
		return nil
	}

	// Формируем текст уведомления в зависимости от типа
	var title, message string
	notifType := models.NotificationType(msg.Type)

	switch notifType {
	case models.NotificationTypeReminder:
		title = "Напоминание о дедлайне"
		if task.Deadline != nil {
			message = fmt.Sprintf("Задача \"%s\" - дедлайн через %s", task.Title, time.Until(*task.Deadline).Round(time.Minute))
		} else {
			message = fmt.Sprintf("Напоминание о задаче \"%s\"", task.Title)
		}
		// Помечаем напоминание как отправленное
		if err := w.taskRepo.MarkReminderSent(ctx, taskID); err != nil {
			logger.Logger.Warnf("Не удалось пометить напоминание как отправленное: %v", err)
		}
	case models.NotificationTypeOverdue:
		title = "Задача просрочена"
		message = fmt.Sprintf("Задача \"%s\" просрочена!", task.Title)

		// Обновляем статус
		if err := w.taskRepo.UpdateStatus(ctx, taskID, models.TaskStatusOverdue); err != nil {
			logger.Logger.Warnf("Не удалось обновить статус задачи: %v", err)
		}
	default:
		title = "Уведомление"
		message = fmt.Sprintf("Не удалось обновить статус задачи: %v", task.Title)
	}

	// Создаем уведомление в базе данных
	notification := &models.Notification{
		UserID:    userID,
		TaskID:    taskID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Read:      false,
		CreatedAt: time.Now(),
	}

	if err := w.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("не удалось сохранить уведомление: %w", err)
	}

	// Отправляем уведомление через websocket
	w.wsHub.SendNotification(userID, notification)
	logger.Logger.Infof("Напоминание отправлено пользователю %s для задачи %s", msg.UserID, msg.TaskID)

	return nil
}
