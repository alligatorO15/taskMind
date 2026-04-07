package handler

import (
	"fmt"
	"net/http"

	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"github.com/alligatorO15/taskmind-backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskHandler - HTTP-хендлер для эндпоинтов управления задачами
type TaskHandler struct {
	taskUC *usecase.TaskUseCase
}

// NewTaskHandler - создает экземпляр TaskHandler
func NewTaskHandler(taskUC *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{taskUC: taskUC}
}

// Create - создает новую задачу
// Привязывает задачу к текущему пользователю, возвращает 201
func (h *TaskHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не автороизован"})
		return
	}

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("некорректные данные запроса: %s", err.Error())})
		return
	}

	task, err := h.taskUC.Create(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetByID - возвразает задачу по идентфиикатору
func (h *TaskHandler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентфикатор задачи"})
		return
	}

	task, err := h.taskUC.GetByID(c.Request.Context(), taskID, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetAll - возвращает список задач пользвоателя с опциональынми фильтрами
// Поддерживает query-параметры: status, priority, project_id, tag, search
func (h *TaskHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	filter := models.TaskFilter{}

	if s := c.Query("status"); s != "" {
		status := models.TaskStatus(s)
		filter.Status = &status
	}
	if p := c.Query("priority"); p != "" {
		priority := models.TaskPriority(p)
		filter.Priority = &priority
	}
	if pid := c.Query("project_id"); pid != "" {
		filter.ProjectID = &pid
	}
	if tag := c.Query("tag"); tag != "" {
		filter.Tag = &tag
	}
	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	tasks, err := h.taskUC.GetByUserID(c.Request.Context(), userID, filter)
	if err != nil {
		handleError(c, err)
		return
	}

	if tasks == nil {
		tasks = []*models.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

// Update — обновляет задачу по идентификатору.
func (h *TaskHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор задачи"})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные запроса: " + err.Error()})
		return
	}

	task, err := h.taskUC.Update(c.Request.Context(), taskID, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// Delete - удаляет задачу по идентфикатору
func (h *TaskHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректныый идентфикатор задачи"})
		return
	}

	if err := h.taskUC.Delete(c.Request.Context(), taskID, userID); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// getUserID - вспомогательная функция для извлечения userID из gin.Context (устанавливает AuthMiddleware)
func getUserID(c *gin.Context) (primitive.ObjectID, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return primitive.NilObjectID, false
	}

	userID, ok := userIDVal.(primitive.ObjectID)
	return userID, ok
}
