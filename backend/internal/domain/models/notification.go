package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType — тип уведомления
type NotificationType string

const (
	NotificationTypeReminder NotificationType = "reminder" // Напоминание о дедлайне
	NotificationTypeOverdue  NotificationType = "overdue"  // Задача просрочена
	NotificationTypeSystem   NotificationType = "system"   // Системное уведомление
)

// Notification — модель уведомления пользователя.
// Уведомления создаются воркерами и доставляются через WebSocket.
type Notification struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	TaskID    primitive.ObjectID `json:"task_id,omitempty" bson:"task_id,omitempty"`
	Type      NotificationType   `json:"type" bson:"type"`
	Title     string             `json:"title" bson:"title"`
	Message   string             `json:"message" bson:"message"`
	Read      bool               `json:"read" bson:"read"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// ReminderMessage — сообщение для очереди RabbitMQ.
// Содержит идентификаторы задачи и пользователя для обработки воркером.
type ReminderMessage struct {
	TaskID string `json:"task_id"`
	UserID string `json:"user_id"`
	Type   string `json:"type"` // "reminder" или "overdue"
}
