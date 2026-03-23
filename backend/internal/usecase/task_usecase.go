package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/config"
	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/domain/repository"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/logger"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/rabbitmq"
	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskUseCase - бизнес-логика управления задачами
// Обрабатывает создание, обновление, удаление задач и планирование напоминаний
type TaskUseCase struct {
	taskRepo    repository.TaskRepository
	projectRepo repository.ProjectRepository
	rabbitConn  *rabbitmq.Connection
	reminderCfg config.ReminderConfig
}

// NewTaskUseCase - создает экземпляр TaskUseCase
func NewTaskUseCase(
	taskRepo repository.TaskRepository,
	projectRepo repository.ProjectRepository,
	rabbitConn *rabbitmq.Connection,
	reminderCfg config.ReminderConfig,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		rabbitConn:  rabbitConn,
		reminderCfg: reminderCfg,
	}
}

// Create - создает новую задачу и планирует напоминание при наличие дедлайна
// Если указан ProjectID, проверяет приинадлежность проекта пользователю
func (uc *TaskUseCase) Create(ctx context.Context, userID primitive.ObjectID, req models.CreateTaskRequest) (*models.Task, error) {
	task := &models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      models.TaskStatusNew,
		Priority:    models.TaskPriority(req.Priority),
		Tags:        req.Tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if task.Tags == nil {
		task.Tags = []string{}
	}

	// Привязка к проекту, если указан
	if req.ProjectID != "" {
		projectID, err := primitive.ObjectIDFromHex(req.ProjectID)
		if err != nil {
			return nil, apperror.ErrInvalidInput
		}
		// проверяем, что проект сууществвует и принадлежит пользователю
		if _, err := uc.projectRepo.FindByID(ctx, projectID, userID); err != nil {
			return nil, fmt.Errorf("проект не найден: %w", err)
		}
		task.ProjectID = projectID

		// увеличиваем счетчик задач проекта
		if err := uc.projectRepo.IncrementTaskCount(ctx, projectID, 1); err != nil {
			logger.Logger.Warnf("Не удалось обновить счётчик задач проекта: %v", err)
		}
	}

	// Парсинг дедлайна
	if req.Deadline != nil {
		deadline, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			return nil, apperror.ErrInvalidInput
		}
		task.Deadline = &deadline
	}

	// Парсинг интервала напоминания
	if req.ReminderBefore != nil {
		duration, err := time.ParseDuration(*req.ReminderBefore)
		if err != nil {
			return nil, apperror.ErrInvalidInput
		}
		task.ReminderBefore = &duration
	} else if task.Deadline != nil {
		// Используем значение по умолчанию из конфига
		defaultBefore := uc.reminderCfg.DefaultBefore
		task.ReminderBefore = &defaultBefore
	}

	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Планирует напоминание через RabbitMQ, если есть дедлайн
	if task.Deadline != nil && uc.rabbitConn != nil {
		uc.scheduleReminder(ctx, task)
	}

	return task, nil
}

// GetByID - возвращает задачу по идентификатору с проверкой принадлежности пользователю
func (uc *TaskUseCase) GetByID(ctx context.Context, taskID, userID primitive.ObjectID) (*models.Task, error) {
	return uc.taskRepo.FindByID(ctx, taskID, userID)
}

// GetByUserID - возвращает список задач пользователя с приминением фильтров
func (uc *TaskUseCase) GetByUserID(ctx context.Context, userID primitive.ObjectID, filter models.TaskFilter) ([]*models.Task, error) {
	return uc.taskRepo.FindByUserID(ctx, userID, filter)
}

// Update - обновляет задачу. Перепланирует напоминание при изменении дедлайна
func (uc *TaskUseCase) Update(ctx context.Context, taskID, userID primitive.ObjectID, req models.UpdateTaskRequest) (*models.Task, error) {
	task, err := uc.taskRepo.FindByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = models.TaskStatus(*req.Status)
	}
	if req.Priority != nil {
		task.Priority = models.TaskPriority(*req.Priority)
	}
	if req.Tags != nil {
		task.Tags = req.Tags
	}

	// Обновление дедлайна с переплнированием напоминания
	deadlineChanged := false
	if req.Deadline != nil {
		deadline, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			return nil, apperror.ErrInvalidInput
		}
		task.Deadline = &deadline
		task.ReminderSent = false
		deadlineChanged = true
	}

	if req.ReminderBefore != nil {
		duration, err := time.ParseDuration(*req.ReminderBefore)
		if err != nil {
			return nil, apperror.ErrInvalidInput
		}
		task.ReminderBefore = &duration
		task.ReminderSent = false
		deadlineChanged = true
	}

	task.UpdatedAt = time.Now()

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Перепланируем напоминание, если изменился дедлайн
	if deadlineChanged && task.Deadline != nil && uc.rabbitConn != nil {
		uc.scheduleReminder(ctx, task)
	}

	return task, nil
}

// Delete - удаляет задачу и уменьшает счётчик задач проекта
func (uc *TaskUseCase) Delete(ctx context.Context, taskID, userID primitive.ObjectID) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	if err := uc.taskRepo.Delete(ctx, taskID, userID); err != nil {
		return err
	}

	// уменьшаем счётчик задач проекта
	if !task.ProjectID.IsZero() {
		if err := uc.projectRepo.IncrementTaskCount(ctx, task.ProjectID, -1); err != nil {
			logger.Logger.Warnf("Не удалось обновить счётчик задач проекта: %v", err)
		}
	}

	return nil
}

// scheduleReminder - планирует отложенное напоминание через RabbitMQ
// Вычисляет задержку как deadline - reminder_before - now
func (uc *TaskUseCase) scheduleReminder(ctx context.Context, task *models.Task) {
	if task.Deadline == nil || task.ReminderBefore == nil {
		return
	}

	reminderTime := task.Deadline.Add(-*task.ReminderBefore)
	delay := time.Until(reminderTime)

	// Если время напоминания уже прошло - отправляем сразу
	if delay < 0 {
		delay = 0
	}

	msg := models.ReminderMessage{
		TaskID: task.ID.Hex(),
		UserID: task.UserID.Hex(),
		Type:   string(models.NotificationTypeReminder),
	}

	if err := uc.rabbitConn.PublishReminder(ctx, msg, delay); err != nil {
		logger.Logger.Errorf("Не удалось запланировать напоминание для задачи %s: %v", task.ID.Hex(), err)
	}

}
