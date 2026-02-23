package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project — модель проекта для группировки задач.
// Каждый проект принадлежит конкретному пользователю.
type Project struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Color       string             `json:"color" bson:"color"` // HEX-цвет для UI
	TaskCount   int                `json:"task_count" bson:"task_count"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateProjectRequest — запрос на создание проекта
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	Color       string `json:"color" binding:"required,hexcolor"`
}

// UpdateProjectRequest — запрос на обновление проекта
type UpdateProjectRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Color       *string `json:"color" binding:"omitempty,hexcolor"`
}
