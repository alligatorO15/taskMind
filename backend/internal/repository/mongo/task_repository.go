package mongo

import (
	"context"
	"regexp"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/domain/repository"
	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// regexMetaChars - регулярное выражение для экранирования спецсимволов в поисковом запросе
var regexMetaChars = regexp.MustCompile(`[.*+?^${}()|[\]\\]`)

const tasksCollection = "tasks"

// TaskRepository - MongoDB-реализация интерфейса repository.TaskRepository
type TaskRepository struct {
	coll *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		coll: db.Collection(tasksCollection),
	}
}

// Create - создает новую задачу
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}

	_, err := r.coll.InsertOne(ctx, task)
	return err
}

// FindByID - находит задачу по идентификатору и идентификатору пользователя
func (r *TaskRepository) FindByID(ctx context.Context, id, userID primitive.ObjectID) (*models.Task, error) {
	var task models.Task
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &task, nil
}

// FindByUserID - возвращаем все задачи пользователя с применением фильтров
func (r *TaskRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, filter models.TaskFilter) ([]*models.Task, error) {
	q := bson.M{"user_id": userID}

	if filter.Status != nil {
		q["status"] = *filter.Status
	}
	if filter.Priority != nil {
		q["priority"] = *filter.Priority
	}
	if filter.ProjectID != nil && *filter.ProjectID != "" {
		projectOID, err := primitive.ObjectIDFromHex(*filter.ProjectID)
		if err == nil {
			q["project_id"] = projectOID
		}
	}
	if filter.Tag != nil && *filter.Tag != "" {
		q["tags"] = *filter.Tag
	}
	if filter.Search != nil && *filter.Search != "" {
		escaped := regexMetaChars.ReplaceAllString(*filter.Search, `\$0`)
		q["title"] = bson.M{"$regex": primitive.Regex{Pattern: escaped, Options: "i"}}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, q, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Update - обновляем существующую задачу
func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	task.UpdatedAt = time.Now()

	result, err := r.coll.UpdateOne(
		ctx,
		bson.M{"_id": task.ID, "user_id": task.UserID},
		bson.D{
			{Key: "$set", Value: bson.M{
				"project_id":      task.ProjectID,
				"title":           task.Title,
				"description":     task.Description,
				"status":          task.Status,
				"priority":        task.Priority,
				"tags":            task.Tags,
				"deadline":        task.Deadline,
				"reminder_before": task.ReminderBefore,
				"reminder_sent":   task.ReminderSent,
				"updated_at":      task.UpdatedAt,
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

// Delete - удаляет задачу по идентификатору
func (r *TaskRepository) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// FindOverdue -  находит задачи с просроченным дедлайномб которые еще не отмечены как overdue
func (r *TaskRepository) FindOverdue(ctx context.Context, now time.Time) ([]*models.Task, error) {
	q := bson.M{
		"deadline": bson.M{"$lt": now},
		"status":   bson.M{"$nin": []models.TaskStatus{models.TaskStatusDone, models.TaskStatusOverdue}},
	}

	cursor, err := r.coll.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindPendingReminders -  находит задачи, которым нужно отправить напоминание.
// Условие: deadline существует, reminder_sent=false, deadline - reminder_before <= now, status != done
// time.Time - миллисек., а time.Duration - наносек. -  поэтому приводим
func (r *TaskRepository) FindPendingReminders(ctx context.Context, now time.Time) ([]*models.Task, error) {
	nowMillis := now.UnixMilli()
	q := bson.M{
		"deadline":        bson.M{"$exists": true, "$ne": nil},
		"reminder_sent":   false,
		"status":          bson.M{"$ne": models.TaskStatusDone},
		"reminder_before": bson.M{"$exists": true, "$ne": nil},
		"$expr": bson.M{
			"$lte": bson.A{
				bson.M{"$subtract": bson.A{"$deadline", bson.M{"$divide": bson.A{"$reminder_before", 1000000}}}},
				nowMillis,
			},
		},
	}

	cursor, err := r.coll.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// MarkReminderSent - помечает, что напоминание для задачи было отправлено
func (r *TaskRepository) MarkReminderSent(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.coll.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$set", Value: bson.M{"reminder_sent": true}}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// UpdateStatus — обновляет статус задачи
func (r *TaskRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.TaskStatus) error {
	result, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$set", Value: bson.M{"status": status, "updated_at": time.Now()}}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

// Проверка, что TaskRepository реализует интерфейс repository.TaskRepository
var _ repository.TaskRepository = (*TaskRepository)(nil)
