package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/config"
	httpDelivery "github.com/alligatorO15/taskMind/backend/internal/delivery/http"
	"github.com/alligatorO15/taskMind/backend/internal/delivery/websocket"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/logger"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/mongodb"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/rabbitmq"
	mongoRepo "github.com/alligatorO15/taskMind/backend/internal/repository/mongo"
	"github.com/alligatorO15/taskMind/backend/internal/usecase"
	"github.com/alligatorO15/taskMind/backend/internal/worker"
	"github.com/gin-gonic/gin"
)

// точка входа: инициализирует все компоненты, запускает HTTP-сервер, воркеры, graceful shutdown
func main() {
	// Загрузка конфига
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	if err := logger.Init(cfg.Server.Mode); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации логгера: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Logger.Info("Запуск TaskMind...")

	gin.SetMode(cfg.Server.Mode)

	// контекст для воркеров
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// подключение к mongodb
	db, err := mongodb.NewConnection(cfg.MongoDB)
	if err != nil {
		logger.Logger.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}
	logger.Logger.Info("MongoDB подключена")

	// создаем инлексы в бд
	if err := mongodb.CreateIndexes(db); err != nil {
		logger.Logger.Warnf("Ошибка создания индексов: %v", err)
	}

	// Подключение к RabbitMQ
	rabbitConn, err := rabbitmq.NewConnection(cfg.RabbitMQ)
	if err != nil {
		logger.Logger.Warnf("Не удалось подключиться к RabbitMQ: %v (напоминания будут недоступны)", err)
	} else {
		defer rabbitConn.Close()
	}

	// инициализация репозиториев
	userRepo := mongoRepo.NewUserRepository(db)
	taskRepo := mongoRepo.NewTaskRepository(db)
	projectRepo := mongoRepo.NewProjectRepository(db)
	notificationRepo := mongoRepo.NewNotificationRepository(db)

	// инициализация юзкейсов
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWT)
	taskUC := usecase.NewTaskUseCase(taskRepo, projectRepo, rabbitConn, cfg.Reminder) // немного нарушил dependency inversion (можно было сделать интерфейс и конструктор)
	projectUC := usecase.NewProjectUseCase(projectRepo)
	notificationUC := usecase.NewNotificationUseCase(notificationRepo)

	// инициализация websocket-хаба для real-time уведомлений
	wsHub := websocket.NewHub(cfg.WebSocket)
	go wsHub.Run()

	// Запуск фоновых воркеров (только если rabbitmq доступен)
	if rabbitConn != nil {
		reminderWorker := worker.NewReminderWorker(rabbitConn, taskRepo, notificationRepo, wsHub)
		if err := reminderWorker.Start(ctx); err != nil {
			logger.Logger.Errorf("Ошибка создания воркера напоминаний: %v", err)
		}

		deadlineChecker := worker.NewDeadlineChecker(rabbitConn, taskRepo, notificationRepo, wsHub, cfg.Reminder.CheckInterval)
		go deadlineChecker.Start(ctx)
	}

	// Настройка HTTP-маршрутов
	router := httpDelivery.NewRouter(authUC, taskUC, projectUC, notificationUC, wsHub)

	// Создание HTTP-сервераа
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// запускаем в отдельнйо горутине чтобы не прерывало main
	go func() {
		logger.Logger.Infof("HTTP-сервер запущен на порту %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatalf("Ошибка HTTP-сервера: %v", err)
		}
	}()

	// грейсфулшатдаун
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.Info("Получен сигнал завершения, выключение сервера...")
	// timout 10 c
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// останавливаем воркеры
	cancel()

	// останавливаем http сервак
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Errorf("Ошибка при завершении HTTP-сервера: %v", err)
	}

	logger.Logger.Info("TaskMind завершил работу")
}
