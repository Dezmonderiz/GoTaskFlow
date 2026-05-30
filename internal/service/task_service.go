package service

import (
	"context"
	"errors"
	"strings"

	"GoTaskFlow/internal/cache"
	"GoTaskFlow/internal/model"
	"GoTaskFlow/internal/repository"
)

var (
	ErrTaskNotFound  = repository.ErrTaskNotFound
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidStatus = errors.New("invalid task status")
)

type TaskService struct {
	repository repository.TaskRepository
	statsCache cache.StatsCache
}

func NewTaskService(repository repository.TaskRepository, statsCache cache.StatsCache) *TaskService {
	return &TaskService{
		repository: repository,
		statsCache: statsCache,
	}
}

func (s *TaskService) Create(ctx context.Context, request model.CreateTaskRequest) (model.Task, error) {
	title := strings.TrimSpace(request.Title)
	if title == "" || len(title) > 255 {
		return model.Task{}, ErrInvalidInput
	}

	task := model.Task{
		Title:       title,
		Description: strings.TrimSpace(request.Description),
	}

	created, err := s.repository.Create(ctx, task)
	if err != nil {
		return model.Task{}, err
	}

	s.invalidateStatsCache(ctx)
	return created, nil
}

func (s *TaskService) List(ctx context.Context) ([]model.Task, error) {
	return s.repository.List(ctx)
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (model.Task, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *TaskService) UpdateStatus(ctx context.Context, id int64, status model.TaskStatus) (model.Task, error) {
	if !status.IsValid() {
		return model.Task{}, ErrInvalidStatus
	}

	task, err := s.repository.UpdateStatus(ctx, id, status)
	if err != nil {
		return model.Task{}, err
	}

	s.invalidateStatsCache(ctx)
	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	s.invalidateStatsCache(ctx)
	return nil
}

func (s *TaskService) Stats(ctx context.Context) (model.TaskStats, error) {
	if s.statsCache != nil {
		stats, err := s.statsCache.Get(ctx)
		if err == nil {
			return stats, nil
		}
	}

	stats, err := s.repository.Stats(ctx)
	if err != nil {
		return model.TaskStats{}, err
	}

	if s.statsCache != nil {
		_ = s.statsCache.Set(ctx, stats)
	}

	return stats, nil
}

func (s *TaskService) invalidateStatsCache(ctx context.Context) {
	if s.statsCache != nil {
		_ = s.statsCache.Delete(ctx)
	}
}
