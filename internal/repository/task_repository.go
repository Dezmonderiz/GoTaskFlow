package repository

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"GoTaskFlow/internal/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(ctx context.Context, task model.Task) (model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id int64) (model.Task, error)
	UpdateStatus(ctx context.Context, id int64, status model.TaskStatus) (model.Task, error)
	Delete(ctx context.Context, id int64) error
	Stats(ctx context.Context) (model.TaskStats, error)
}

type MemoryTaskRepository struct {
	mu     sync.RWMutex
	nextID int64
	tasks  map[int64]model.Task
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		nextID: 1,
		tasks:  make(map[int64]model.Task),
	}
}

func (r *MemoryTaskRepository) Create(_ context.Context, task model.Task) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	task.ID = r.nextID
	task.Status = model.StatusTodo
	task.CreatedAt = now
	task.UpdatedAt = now

	r.tasks[task.ID] = task
	r.nextID++

	return task, nil
}

func (r *MemoryTaskRepository) List(_ context.Context) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]model.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks, nil
}

func (r *MemoryTaskRepository) GetByID(_ context.Context, id int64) (model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (r *MemoryTaskRepository) UpdateStatus(_ context.Context, id int64, status model.TaskStatus) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.tasks[id]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}

	task.Status = status
	task.UpdatedAt = time.Now().UTC()
	r.tasks[id] = task

	return task, nil
}

func (r *MemoryTaskRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

func (r *MemoryTaskRepository) Stats(_ context.Context) (model.TaskStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := model.TaskStats{Total: len(r.tasks)}
	for _, task := range r.tasks {
		switch task.Status {
		case model.StatusTodo:
			stats.Todo++
		case model.StatusDone:
			stats.Done++
		case model.StatusInProgress:
			stats.InProgress++
		}
	}

	return stats, nil
}
