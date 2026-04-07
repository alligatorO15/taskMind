package mongo

import (
	"context"
	"time"

	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"github.com/alligatorO15/taskmind-backend/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const usersCollection = "users"

// UserRepository - MongoDB-реализация интерфейса repository.UserRepository (обертка нужна чтобы скрыть о сервисов конкретную бд)
type UserRepository struct {
	coll *mongo.Collection
}

// NewUserRepository - конструктор репозитория пользователя
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		coll: db.Collection(usersCollection),
	}
}

// Create - создает нового пользователя  в базе данных
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		// Обработка дубликата email или username (код 11000)
		if mongo.IsDuplicateKeyError(err) {
			return apperror.ErrAlreadyExists
		}
		return err
	}
	return nil
}

// FindByID - находит пользователя по его уникальном идентификатору
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail - находит пользователя по email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername - находит пользвоателя по нику
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Проверка, что UserRepository реализует интерфейс repository.UserRepository
var _ repository.UserRepository = (*UserRepository)(nil)
