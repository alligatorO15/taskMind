package http

import (
	"github.com/alligatorO15/taskMind/backend/internal/delivery/http/handler"
	"github.com/alligatorO15/taskMind/backend/internal/delivery/http/middleware"
	"github.com/alligatorO15/taskMind/backend/internal/delivery/websocket"
	"github.com/alligatorO15/taskMind/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewRouter - создает и настраивает маршруты HTTP API
// Регистрирует публичные (auth) и защищённые (tasks, projects,notifications) эндпоинты
func NewRouter(
	authUC *usecase.AuthUseCase,
	taskUC *usecase.TaskUseCase,
	projectUC *usecase.ProjectUseCase,
	notificationUC *usecase.NotificationUseCase,
	hub *websocket.Hub,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// handlers
	authHandler := handler.NewAuthHandler(authUC)
	taskHandler := handler.NewTaskHandler(taskUC)
	projectHandler := handler.NewProjectHandler(projectUC)
	notificationHandler := handler.NewNotificationHandler(notificationUC)

	// API V1
	api := router.Group("/api/v1")

	// Публичные маршруты аутентификации
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Защищенные маршруты (требуют JWT)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authUC))
	{
		// Профиль
		protected.GET("/profile", authHandler.GetProfile)

		// Задачи
		tasks := protected.Group("/tasks")
		{
			tasks.POST("", taskHandler.Create)
			tasks.GET("", taskHandler.GetAll)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
		}

		// Проекты
		projects := protected.Group("/projects")
		{
			projects.POST("", projectHandler.Create)
			projects.GET("", projectHandler.GetAll)
			projects.GET("/:id", projectHandler.GetByID)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
		}

		// Уведомления (специфичные маршруты до параметризованных)
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", notificationHandler.GetAll)
			notifications.GET("/unread-count", notificationHandler.CountUnread)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		}

		// WebSocket для уведомлений в реальном времени
		protected.GET("/ws", func(c *gin.Context) {
			userIDVal, exists := c.Get("userID")
			if !exists {
				return // не можем в websocket отправить http json ответ
			}
			userID := userIDVal.(primitive.ObjectID)
			hub.HandleWebSocket(c.Writer, c.Request, userID)
		})
	}

	return router
}
