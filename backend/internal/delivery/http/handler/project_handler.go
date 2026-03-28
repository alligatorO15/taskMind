package handler

import (
	"fmt"
	"net/http"

	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/usecase"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectHandler — HTTP-хендлер для эндпоинтов управления проектами.
type ProjectHandler struct {
	projectUC *usecase.ProjectUseCase
}

// NewProjectHandler — создает экземпляр ProjectHandler.
func NewProjectHandler(projectUC *usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{projectUC: projectUC}
}

// Create — создает новый проект.
// Извлекает userID из контекста, возвращает 201.
func (h *ProjectHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("некорректные данные запроса: %s", err.Error())})
		return
	}

	project, err := h.projectUC.Create(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetByID — возвращает проект по идентификатору.
func (h *ProjectHandler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор проекта"})
		return
	}

	project, err := h.projectUC.GetByID(c.Request.Context(), projectID, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, project)
}

// GetAll — возвращает все проекты текущего пользователя.
func (h *ProjectHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	projects, err := h.projectUC.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	if projects == nil {
		projects = []*models.Project{}
	}

	c.JSON(http.StatusOK, projects)
}

// Update — обновляет проект по идентификатору.
func (h *ProjectHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор проекта"})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("некорректные данные запроса: %s", err.Error())})
		return
	}

	project, err := h.projectUC.Update(c.Request.Context(), projectID, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, project)
}

// Delete — удаляет проект по идентификатору.
func (h *ProjectHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор проекта"})
		return
	}

	err = h.projectUC.Delete(c.Request.Context(), projectID, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
