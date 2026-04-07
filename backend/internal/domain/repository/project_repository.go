package repository

import (
	"context"

	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectRepository — интерфейс для работы с хранилищем проектов
type ProjectRepository interface {
	// Create — создает новый проект
	Create(ctx context.Context, project *models.Project) error
	// FindByID — находит проект по идентификатору
	FindByID(ctx context.Context, id, userID primitive.ObjectID) (*models.Project, error)
	// FindByUserID — возвращает все проекты пользователя
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Project, error)
	// Update — обновляет проект
	Update(ctx context.Context, project *models.Project) error
	// Delete — удаляет проект
	Delete(ctx context.Context, id, userID primitive.ObjectID) error
	// IncrementTaskCount — увеличивает счетчик задач проекта на указанное значение
	IncrementTaskCount(ctx context.Context, id primitive.ObjectID, delta int) error
}
