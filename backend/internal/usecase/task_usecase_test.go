package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/config"
	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockTaskRepo — мок-реализация TaskRepository для тестирования
type mockTaskRepo struct {
	tasks map[primitive.ObjectID]*models.Task
}

func newMockTaskRepo() *mockTaskRepo {
	return &mockTaskRepo{
		tasks: make(map[primitive.ObjectID]*models.Task),
	}
}

func (m *mockTaskRepo) Create(ctx context.Context, task *models.Task) error {
	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) FindByID(ctx context.Context, id, userID primitive.ObjectID) (*models.Task, error) {
	task, ok := m.tasks[id]
	if !ok || task.UserID != userID {
		return nil, apperror.ErrNotFound
	}
	return task, nil
}

func (m *mockTaskRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID, filter models.TaskFilter) ([]*models.Task, error) {
	var result []*models.Task
	for _, t := range m.tasks {
		if t.UserID != userID {
			continue
		}
		if filter.Status != nil && t.Status != *filter.Status {
			continue
		}
		if filter.Priority != nil && t.Priority != *filter.Priority {
			continue
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *mockTaskRepo) Update(ctx context.Context, task *models.Task) error {
	if _, ok := m.tasks[task.ID]; !ok {
		return apperror.ErrNotFound
	}
	task.UpdatedAt = time.Now()
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepo) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	task, ok := m.tasks[id]
	if !ok || task.UserID != userID {
		return apperror.ErrNotFound
	}
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepo) FindOverdue(ctx context.Context, now time.Time) ([]*models.Task, error) {
	var result []*models.Task
	for _, t := range m.tasks {
		if t.Deadline != nil && t.Deadline.Before(now) &&
			t.Status != models.TaskStatusDone && t.Status != models.TaskStatusOverdue {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTaskRepo) FindPendingReminders(ctx context.Context, now time.Time) ([]*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepo) MarkReminderSent(ctx context.Context, id primitive.ObjectID) error {
	if task, ok := m.tasks[id]; ok {
		task.ReminderSent = true
		return nil
	}
	return apperror.ErrNotFound
}

func (m *mockTaskRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.TaskStatus) error {
	if task, ok := m.tasks[id]; ok {
		task.Status = status
		return nil
	}
	return apperror.ErrNotFound
}

// mockProjectRepo — мок-реализация ProjectRepository
type mockProjectRepo struct {
	projects map[primitive.ObjectID]*models.Project
}

func newMockProjectRepo() *mockProjectRepo {
	return &mockProjectRepo{
		projects: make(map[primitive.ObjectID]*models.Project),
	}
}

func (m *mockProjectRepo) Create(ctx context.Context, project *models.Project) error {
	if project.ID.IsZero() {
		project.ID = primitive.NewObjectID()
	}
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepo) FindByID(ctx context.Context, id, userID primitive.ObjectID) (*models.Project, error) {
	p, ok := m.projects[id]
	if !ok || p.UserID != userID {
		return nil, apperror.ErrNotFound
	}
	return p, nil
}

func (m *mockProjectRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Project, error) {
	var result []*models.Project
	for _, p := range m.projects {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockProjectRepo) Update(ctx context.Context, project *models.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepo) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	delete(m.projects, id)
	return nil
}

func (m *mockProjectRepo) IncrementTaskCount(ctx context.Context, id primitive.ObjectID, delta int) error {
	if p, ok := m.projects[id]; ok {
		p.TaskCount += delta
		return nil
	}
	return apperror.ErrNotFound
}

// testReminderConfig — тестовая конфигурация напоминаний
func testReminderConfig() config.ReminderConfig {
	return config.ReminderConfig{
		DefaultBefore: 30 * time.Minute,
		CheckInterval: 1 * time.Minute,
	}
}

// TestCreateTask_Success — проверяет успешное создание задачи
func TestCreateTask_Success(t *testing.T) {
	taskRepo := newMockTaskRepo()
	projectRepo := newMockProjectRepo()
	uc := NewTaskUseCase(taskRepo, projectRepo, nil, testReminderConfig())

	userID := primitive.NewObjectID()
	req := models.CreateTaskRequest{
		Title:       "Тестовая задача",
		Description: "Описание тестовой задачи",
		Priority:    "high",
		Tags:        []string{"test", "important"},
	}

	task, err := uc.Create(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("Ошибка создания задачи: %v", err)
	}

	if task.Title != "Тестовая задача" {
		t.Errorf("Неверное название: ожидалось 'Тестовая задача', получили '%s'", task.Title)
	}
	if task.Status != models.TaskStatusNew {
		t.Errorf("Неверный статус: ожидался 'new', получили '%s'", task.Status)
	}
	if task.Priority != models.TaskPriorityHigh {
		t.Errorf("Неверный приоритет: ожидался 'high', получили '%s'", task.Priority)
	}
	if len(task.Tags) != 2 {
		t.Errorf("Неверное количество тегов: ожидалось 2, получили %d", len(task.Tags))
	}
}

// TestCreateTask_WithDeadline — проверяет создание задачи с дедлайном
func TestCreateTask_WithDeadline(t *testing.T) {
	taskRepo := newMockTaskRepo()
	projectRepo := newMockProjectRepo()
	uc := NewTaskUseCase(taskRepo, projectRepo, nil, testReminderConfig())

	userID := primitive.NewObjectID()
	deadline := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	req := models.CreateTaskRequest{
		Title:    "Задача с дедлайном",
		Priority: "medium",
		Deadline: &deadline,
	}

	task, err := uc.Create(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("Ошибка создания задачи: %v", err)
	}

	if task.Deadline == nil {
		t.Fatal("Дедлайн не должен быть nil")
	}

	// Проверяем, что reminder_before установлен по умолчанию
	if task.ReminderBefore == nil {
		t.Fatal("ReminderBefore не должен быть nil (должен быть установлен по умолчанию)")
	}
	if *task.ReminderBefore != 30*time.Minute {
		t.Errorf("Неверный ReminderBefore: ожидался 30m, получили %v", *task.ReminderBefore)
	}
}

// TestCreateTask_WithProject — проверяет создание задачи с привязкой к проекту
func TestCreateTask_WithProject(t *testing.T) {
	taskRepo := newMockTaskRepo()
	projectRepo := newMockProjectRepo()
	uc := NewTaskUseCase(taskRepo, projectRepo, nil, testReminderConfig())

	userID := primitive.NewObjectID()

	// Создаем проект
	project := &models.Project{
		UserID: userID,
		Name:   "Тестовый проект",
		Color:  "#FF0000",
	}
	_ = projectRepo.Create(context.Background(), project)

	req := models.CreateTaskRequest{
		Title:     "Задача в проекте",
		Priority:  "low",
		ProjectID: project.ID.Hex(),
	}

	task, err := uc.Create(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("Ошибка создания задачи: %v", err)
	}

	if task.ProjectID != project.ID {
		t.Errorf("ProjectID не совпадает")
	}

	// Проверяем, что счётчик задач увеличился
	if project.TaskCount != 1 {
		t.Errorf("Счётчик задач проекта: ожидался 1, получили %d", project.TaskCount)
	}
}

// TestUpdateTask_Success — проверяет успешное обновление задачи
func TestUpdateTask_Success(t *testing.T) {
	taskRepo := newMockTaskRepo()
	projectRepo := newMockProjectRepo()
	uc := NewTaskUseCase(taskRepo, projectRepo, nil, testReminderConfig())

	userID := primitive.NewObjectID()
	req := models.CreateTaskRequest{
		Title:    "Исходная задача",
		Priority: "low",
	}

	task, err := uc.Create(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("Ошибка создания задачи: %v", err)
	}

	// Обновляем
	newTitle := "Обновленная задача"
	newStatus := "in_progress"
	updateReq := models.UpdateTaskRequest{
		Title:  &newTitle,
		Status: &newStatus,
	}

	updated, err := uc.Update(context.Background(), task.ID, userID, updateReq)
	if err != nil {
		t.Fatalf("Ошибка обновления задачи: %v", err)
	}

	if updated.Title != "Обновленная задача" {
		t.Errorf("Неверное название после обновления: '%s'", updated.Title)
	}
	if updated.Status != models.TaskStatusInProgress {
		t.Errorf("Неверный статус после обновления: '%s'", updated.Status)
	}
}

// TestDeleteTask_Success — проверяет успешное удаление задачи
func TestDeleteTask_Success(t *testing.T) {
	taskRepo := newMockTaskRepo()
	projectRepo := newMockProjectRepo()
	uc := NewTaskUseCase(taskRepo, projectRepo, nil, testReminderConfig())

	userID := primitive.NewObjectID()
	req := models.CreateTaskRequest{
		Title:    "Задача для удаления",
		Priority: "low",
	}

	task, err := uc.Create(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("Ошибка создания задачи: %v", err)
	}

	err = uc.Delete(context.Background(), task.ID, userID)
	if err != nil {
		t.Fatalf("Ошибка удаления задачи: %v", err)
	}

	// Проверяем, что задача удалена
	_, err = uc.GetByID(context.Background(), task.ID, userID)
	if err != apperror.ErrNotFound {
		t.Errorf("Ожидалась ошибка ErrNotFound после удаления, получили: %v", err)
	}
}

// TestDeleteTask_NotFound — проверяет ошибку при удалении несуществующей задачи
func TestDeleteTask_NotFound(t *testing.T) {
	taskRepo := newMockTaskRepo()
	projectRepo := newMockProjectRepo()
	uc := NewTaskUseCase(taskRepo, projectRepo, nil, testReminderConfig())

	userID := primitive.NewObjectID()
	fakeID := primitive.NewObjectID()

	err := uc.Delete(context.Background(), fakeID, userID)
	if err != apperror.ErrNotFound {
		t.Errorf("Ожидалась ошибка ErrNotFound, получили: %v", err)
	}
}

// TestGetTasksByFilter — проверяет фильтрацию задач по статусу
func TestGetTasksByFilter(t *testing.T) {
	taskRepo := newMockTaskRepo()
	projectRepo := newMockProjectRepo()
	uc := NewTaskUseCase(taskRepo, projectRepo, nil, testReminderConfig())

	userID := primitive.NewObjectID()

	// Создаем задачи с разными статусами
	_, _ = uc.Create(context.Background(), userID, models.CreateTaskRequest{
		Title: "Задача 1", Priority: "low",
	})

	task2, _ := uc.Create(context.Background(), userID, models.CreateTaskRequest{
		Title: "Задача 2", Priority: "high",
	})
	// Меняем статус второй задачи
	newStatus := "in_progress"
	_, _ = uc.Update(context.Background(), task2.ID, userID, models.UpdateTaskRequest{
		Status: &newStatus,
	})

	// Фильтруем по статусу "new"
	statusNew := models.TaskStatusNew
	tasks, err := uc.GetByUserID(context.Background(), userID, models.TaskFilter{Status: &statusNew})
	if err != nil {
		t.Fatalf("Ошибка получения задач: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Ожидалась 1 задача со статусом 'new', получили %d", len(tasks))
	}
}
