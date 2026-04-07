package repository

import (
	"context"

	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository - интерфейс для работы с хранилищем пользователей.
// Определяет контракт, которому должна следовать любая реализация (в нашем случае, mongodb)
// Соответствует принципу инверсии зависисмостей
type UserRepository interface {
	// Create - создает нового пользователя в бд
	Create(ctx context.Context, user *models.User) error
	// FindByID - находит пользователя по его id
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	// FindByEmail - находит пользователя по его email
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	// FindByUsername - находит пользователя по его username
	FindByUsername(ctx context.Context, username string) (*models.User, error)
}
