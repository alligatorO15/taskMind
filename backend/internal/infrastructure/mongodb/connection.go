package mongodb

import (
	"context"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/config"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewConnection — устанавливает подключение к MongoDB.
// Создает клиент с заданным таймаутом, проверяет доступность сервера (ping),
// и возвращает ссылку на базу данных.
func NewConnection(cfg config.MongoDBConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Проверяем доступность сервера MongoDB
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	logger.Logger.Infof("Подключение к MongoDB установлено: %s", cfg.Database)
	return client.Database(cfg.Database), nil
}

// CreateIndexes — создает необходимые индексы в коллекциях MongoDB.
// Вызывается при старте приложения для обеспечения производительности запросов.
func CreateIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Уникальный индекс по email пользователя
	_, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Уникальный индекс по username пользователя
	_, err = db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"username": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Индекс для поиска задач по user_id и статусу
	_, err = db.Collection("tasks").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"user_id": 1, "status": 1},
	})
	if err != nil {
		return err
	}

	// Индекс для поиска задач по дедлайну (используется воркером напоминаний)
	_, err = db.Collection("tasks").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"deadline": 1, "reminder_sent": 1},
	})
	if err != nil {
		return err
	}

	// Индекс для поиска проектов по user_id
	_, err = db.Collection("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"user_id": 1},
	})
	if err != nil {
		return err
	}

	// Индекс для поиска уведомлений по user_id
	_, err = db.Collection("notifications").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"user_id": 1, "created_at": -1},
	})
	if err != nil {
		return err
	}

	logger.Logger.Info("Индексы MongoDB успешно созданы")
	return nil
}
