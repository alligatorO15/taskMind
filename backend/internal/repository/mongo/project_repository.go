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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const projectsCollection = "projects"

// ProjectRepository — MongoDB-реализация интерфейса repository.ProjectRepository
type ProjectRepository struct {
	coll *mongo.Collection
}

// NewProjectRepository — конструктор репозитория проектов
func NewProjectRepository(db *mongo.Database) *ProjectRepository {
	return &ProjectRepository{
		coll: db.Collection(projectsCollection),
	}
}

// Create — создает новый проект
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	if project.ID.IsZero() {
		project.ID = primitive.NewObjectID()
	}

	_, err := r.coll.InsertOne(ctx, project)
	return err
}

// FindByID — находит проект по идентификатору и идентификатору пользователя
func (r *ProjectRepository) FindByID(ctx context.Context, id, userID primitive.ObjectID) (*models.Project, error) {
	var project models.Project
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &project, nil
}

// FindByUserID — возвращает все проекты пользователя, отсортированные по created_at по убыванию
func (r *ProjectRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Project, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// Update — обновляет проект по идентификатору и идентификатору пользователя
func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	project.UpdatedAt = time.Now()

	result, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": project.ID, "user_id": project.UserID},
		bson.D{
			{Key: "$set", Value: bson.M{
				"name":        project.Name,
				"description": project.Description,
				"color":       project.Color,
				"updated_at":  project.UpdatedAt,
			}},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// Delete — удаляет проект по идентификатору и идентификатору пользователя
func (r *ProjectRepository) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// IncrementTaskCount — атомарно увеличивает счетчик задач проекта на delta с помощью $inc
func (r *ProjectRepository) IncrementTaskCount(ctx context.Context, id primitive.ObjectID, delta int) error {
	result, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$inc", Value: bson.M{"task_count": delta}}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// Проверка, что ProjectRepository реализует интерфейс repository.ProjectRepository
var _ repository.ProjectRepository = (*ProjectRepository)(nil)
