package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskStatus — статус задачи в системе
type TaskStatus string

const (
	TaskStatusNew        TaskStatus = "new"         // Новая задача
	TaskStatusInProgress TaskStatus = "in_progress" // В работе
	TaskStatusDone       TaskStatus = "done"        // Завершена
	TaskStatusOverdue    TaskStatus = "overdue"     // Просрочена
)

// TaskPriority — приоритет задачи
type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"    // Низкий приоритет
	TaskPriorityMedium TaskPriority = "medium" // Средний приоритет
	TaskPriorityHigh   TaskPriority = "high"   // Высокий приоритет
)

// Task — модель задачи.
// Содержит описание, статус, приоритет, привязку к проекту и теги.
type Task struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProjectID      primitive.ObjectID `json:"project_id,omitempty" bson:"project_id,omitempty"`
	Title          string             `json:"title" bson:"title"`
	Description    string             `json:"description" bson:"description"`
	Status         TaskStatus         `json:"status" bson:"status"`
	Priority       TaskPriority       `json:"priority" bson:"priority"`
	Tags           []string           `json:"tags" bson:"tags"`
	Deadline       *time.Time         `json:"deadline,omitempty" bson:"deadline,omitempty"`
	ReminderBefore *time.Duration     `json:"reminder_before,omitempty" bson:"reminder_before,omitempty"`
	ReminderSent   bool               `json:"reminder_sent" bson:"reminder_sent"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateTaskRequest — запрос на создание новой задачи
type CreateTaskRequest struct {
	ProjectID      string   `json:"project_id"`
	Title          string   `json:"title" binding:"required,min=1,max=200"`
	Description    string   `json:"description" binding:"max=2000"`
	Priority       string   `json:"priority" binding:"required,oneof=low medium high"`
	Tags           []string `json:"tags"`
	Deadline       *string  `json:"deadline"`
	ReminderBefore *string  `json:"reminder_before"` // формат: "30m", "1h", "24h"
}

// UpdateTaskRequest — запрос на обновление задачи
type UpdateTaskRequest struct {
	Title          *string  `json:"title" binding:"omitempty,min=1,max=200"`
	Description    *string  `json:"description" binding:"omitempty,max=2000"`
	Status         *string  `json:"status" binding:"omitempty,oneof=new in_progress done"`
	Priority       *string  `json:"priority" binding:"omitempty,oneof=low medium high"`
	Tags           []string `json:"tags"`
	Deadline       *string  `json:"deadline"`
	ReminderBefore *string  `json:"reminder_before"`
}

// TaskFilter — фильтры для поиска задач
type TaskFilter struct {
	Status    *TaskStatus   `json:"status"`
	Priority  *TaskPriority `json:"priority"`
	ProjectID *string       `json:"project_id"`
	Tag       *string       `json:"tag"`
	Search    *string       `json:"search"`
}

// TaskResponse — ответ API с данными задачи
type TaskResponse struct {
	ID             primitive.ObjectID `json:"id"`
	UserID         primitive.ObjectID `json:"user_id"`
	ProjectID      primitive.ObjectID `json:"project_id,omitempty"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	Status         TaskStatus         `json:"status"`
	Priority       TaskPriority       `json:"priority"`
	Tags           []string           `json:"tags"`
	Deadline       *time.Time         `json:"deadline,omitempty"`
	ReminderBefore *time.Duration     `json:"reminder_before,omitempty"`
	ReminderSent   bool               `json:"reminder_sent"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
