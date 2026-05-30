package service_test

import (
	"context"
	"errors"
	"testing"

	"GoTaskFlow/internal/model"
	"GoTaskFlow/internal/repository"
	"GoTaskFlow/internal/service"
)

func TestCreateTaskWithValidData(t *testing.T) {
	taskService := newTestTaskService()

	task, err := taskService.Create(context.Background(), model.CreateTaskRequest{
		Title:       "Write README",
		Description: "Describe project setup",
	})
	if err != nil {
		t.Fatalf("expected task to be created, got error: %v", err)
	}

	if task.ID == 0 {
		t.Fatal("expected created task to have id")
	}
	if task.Title != "Write README" {
		t.Fatalf("expected title %q, got %q", "Write README", task.Title)
	}
	if task.Description != "Describe project setup" {
		t.Fatalf("expected description %q, got %q", "Describe project setup", task.Description)
	}
	if task.Status != model.StatusTodo {
		t.Fatalf("expected default status %q, got %q", model.StatusTodo, task.Status)
	}
}

func TestCreateTaskWithoutTitleReturnsError(t *testing.T) {
	taskService := newTestTaskService()

	_, err := taskService.Create(context.Background(), model.CreateTaskRequest{
		Description: "Missing title",
	})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	taskService := newTestTaskService()

	task, err := taskService.Create(context.Background(), model.CreateTaskRequest{Title: "Implement API"})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	updated, err := taskService.UpdateStatus(context.Background(), task.ID, model.StatusInProgress)
	if err != nil {
		t.Fatalf("expected status to be updated, got error: %v", err)
	}

	if updated.Status != model.StatusInProgress {
		t.Fatalf("expected status %q, got %q", model.StatusInProgress, updated.Status)
	}
	if !updated.UpdatedAt.After(updated.CreatedAt) && !updated.UpdatedAt.Equal(updated.CreatedAt) {
		t.Fatalf("updated_at should not be before created_at: created=%s updated=%s", updated.CreatedAt, updated.UpdatedAt)
	}
}

func TestStatsCountsTasksByStatus(t *testing.T) {
	taskService := newTestTaskService()

	todoTask, err := taskService.Create(context.Background(), model.CreateTaskRequest{Title: "Todo task"})
	if err != nil {
		t.Fatalf("failed to create todo task: %v", err)
	}
	inProgressTask, err := taskService.Create(context.Background(), model.CreateTaskRequest{Title: "In progress task"})
	if err != nil {
		t.Fatalf("failed to create in-progress task: %v", err)
	}
	doneTask, err := taskService.Create(context.Background(), model.CreateTaskRequest{Title: "Done task"})
	if err != nil {
		t.Fatalf("failed to create done task: %v", err)
	}

	if _, err := taskService.UpdateStatus(context.Background(), inProgressTask.ID, model.StatusInProgress); err != nil {
		t.Fatalf("failed to update in-progress task: %v", err)
	}
	if _, err := taskService.UpdateStatus(context.Background(), doneTask.ID, model.StatusDone); err != nil {
		t.Fatalf("failed to update done task: %v", err)
	}

	stats, err := taskService.Stats(context.Background())
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}

	if stats.Total != 3 || stats.Todo != 1 || stats.InProgress != 1 || stats.Done != 1 {
		t.Fatalf("unexpected stats after task updates: %+v; todo task id=%d", stats, todoTask.ID)
	}
}

func TestStatsUsesCache(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	statsCache := &fakeStatsCache{}
	taskService := service.NewTaskService(repo, statsCache)

	if _, err := taskService.Create(context.Background(), model.CreateTaskRequest{Title: "Cached task"}); err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	stats, err := taskService.Stats(context.Background())
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}
	if stats.Total != 1 || stats.Todo != 1 {
		t.Fatalf("unexpected stats from repository: %+v", stats)
	}
	if statsCache.setCalls != 1 {
		t.Fatalf("expected stats to be cached once, got %d", statsCache.setCalls)
	}

	statsCache.stats = model.TaskStats{Total: 99, Todo: 99}
	stats, err = taskService.Stats(context.Background())
	if err != nil {
		t.Fatalf("failed to get cached stats: %v", err)
	}
	if stats.Total != 99 || stats.Todo != 99 {
		t.Fatalf("expected cached stats, got %+v", stats)
	}
}

func TestTaskChangesInvalidateStatsCache(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	statsCache := &fakeStatsCache{stats: model.TaskStats{Total: 1, Todo: 1}, exists: true}
	taskService := service.NewTaskService(repo, statsCache)

	task, err := taskService.Create(context.Background(), model.CreateTaskRequest{Title: "Invalidate cache"})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	if statsCache.deleteCalls != 1 {
		t.Fatalf("expected create to invalidate cache once, got %d", statsCache.deleteCalls)
	}

	if _, err := taskService.UpdateStatus(context.Background(), task.ID, model.StatusDone); err != nil {
		t.Fatalf("failed to update task status: %v", err)
	}
	if statsCache.deleteCalls != 2 {
		t.Fatalf("expected update to invalidate cache twice, got %d", statsCache.deleteCalls)
	}

	if err := taskService.Delete(context.Background(), task.ID); err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}
	if statsCache.deleteCalls != 3 {
		t.Fatalf("expected delete to invalidate cache three times, got %d", statsCache.deleteCalls)
	}
}

type fakeStatsCache struct {
	stats       model.TaskStats
	exists      bool
	setCalls    int
	deleteCalls int
}

func (c *fakeStatsCache) Get(context.Context) (model.TaskStats, error) {
	if !c.exists {
		return model.TaskStats{}, errFakeCacheMiss
	}

	return c.stats, nil
}

func (c *fakeStatsCache) Set(_ context.Context, stats model.TaskStats) error {
	c.stats = stats
	c.exists = true
	c.setCalls++
	return nil
}

func (c *fakeStatsCache) Delete(context.Context) error {
	c.exists = false
	c.deleteCalls++
	return nil
}

type fakeCacheMissError struct{}

func (fakeCacheMissError) Error() string {
	return "cache miss"
}

var errFakeCacheMiss fakeCacheMissError

func newTestTaskService() *service.TaskService {
	return service.NewTaskService(repository.NewMemoryTaskRepository(), nil)
}
