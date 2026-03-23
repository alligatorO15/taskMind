package usecase

import (
	"context"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectUseCase — бизнес-логика управления проектами
type ProjectUseCase struct {
	projectRepo repository.ProjectRepository
}

// NewProjectUseCase — создает экземпляр ProjectUseCase
func NewProjectUseCase(projectRepo repository.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo: projectRepo,
	}
}

// Create — создает новый проект для пользователя
func (uc *ProjectUseCase) Create(ctx context.Context, userID primitive.ObjectID, req models.CreateProjectRequest) (*models.Project, error) {
	project := &models.Project{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		TaskCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetByID — возвращает проект по идентификатору
func (uc *ProjectUseCase) GetByID(ctx context.Context, projectID, userID primitive.ObjectID) (*models.Project, error) {
	return uc.projectRepo.FindByID(ctx, projectID, userID)
}

// GetByUserID — возвращает все проекты пользователя
func (uc *ProjectUseCase) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Project, error) {
	return uc.projectRepo.FindByUserID(ctx, userID)
}

// Update — обновляет проект
func (uc *ProjectUseCase) Update(ctx context.Context, projectID, userID primitive.ObjectID, req models.UpdateProjectRequest) (*models.Project, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Color != nil {
		project.Color = *req.Color
	}
	project.UpdatedAt = time.Now()

	if err := uc.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete — удаляет проект
func (uc *ProjectUseCase) Delete(ctx context.Context, projectID, userID primitive.ObjectID) error {
	return uc.projectRepo.Delete(ctx, projectID, userID)
}
